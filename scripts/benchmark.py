"""Compare two dct builds with deterministic fixtures and per-process peak memory."""

import argparse
from dataclasses import dataclass
import hashlib
import json
import os
from pathlib import Path
import platform
import signal
from statistics import median
import subprocess
import sys
from tempfile import TemporaryDirectory
from time import perf_counter_ns


MIB = 1024 * 1024
NEEDLE = "DCT_BENCHMARK_MATCH"
BENCH_DIR = Path(__file__).resolve().parents[1] / "bin" / ".benchmark"
PROCESS_TIMEOUT_SECONDS = 180


@dataclass(frozen=True)
class Workload:
    name: str
    root: Path
    input_bytes: int
    records: int
    stdout_bytes: int
    stdout_sha256: str

    @property
    def exit_code(self):
        return 0 if self.records else 1


def write_repeated(file, unit, count):
    units_per_block = min(count, MIB // len(unit))
    if not units_per_block:
        return
    full_blocks, tail_units = divmod(count, units_per_block)
    block = unit * units_per_block
    for _ in range(full_blocks):
        file.write(block)
    file.write(unit * tail_units)


def record(path, line, text):
    output = {"path": path, "line": line, "text": text}
    return json.dumps(output, separators=(",", ":")).encode("ascii") + b"\n"


def short_line_fixture(root, name, file_count, lines_per_file, every_line_matches):
    directory = root / name
    directory.mkdir()
    miss = b"x" * 127 + b"\n"
    matched = NEEDLE.ljust(127, "x")
    expected = hashlib.sha256()
    record_count = 0
    output_count = 0
    name_width = max(2, len(str(file_count - 1)))
    for index in range(file_count):
        filename = f"file-{index:0{name_width}d}.txt"
        with (directory / filename).open("wb") as file:
            if every_line_matches:
                write_repeated(file, matched.encode("ascii") + b"\n", lines_per_file)
                matching_lines = range(1, lines_per_file + 1)
            else:
                write_repeated(file, miss, lines_per_file - 1)
                file.write(matched.encode("ascii") + b"\n")
                matching_lines = (lines_per_file,)
        for line in matching_lines:
            output = record(filename, line, matched)
            expected.update(output)
            output_count += len(output)
            record_count += 1
    return Workload(name, directory, file_count * lines_per_file * 128,
                    record_count, output_count, expected.hexdigest())


def long_line_fixture(root, name, payload_size, matches):
    directory = root / name
    directory.mkdir()
    expected = hashlib.sha256()
    with (directory / "line.txt").open("wb") as file:
        write_repeated(file, b"x", payload_size)
        if matches:
            file.write(NEEDLE.encode("ascii") + b"\r\n")
    if not matches:
        return Workload(name, directory, payload_size, 0, 0, expected.hexdigest())

    prefix = b'{"path":"line.txt","line":1,"text":"'
    suffix = NEEDLE.encode("ascii") + b'"}\n'
    expected.update(prefix)
    block = b"x" * MIB
    full_blocks, remainder = divmod(payload_size, MIB)
    for _ in range(full_blocks):
        expected.update(block)
    expected.update(block[:remainder])
    expected.update(suffix)
    return Workload(name, directory, payload_size + len(NEEDLE) + 2, 1,
                    len(prefix) + payload_size + len(suffix), expected.hexdigest())


def make_workloads(root):
    return (
        short_line_fixture(root, "sparse-short-lines", 4, 24 * MIB // 128, False),
        short_line_fixture(root, "dense-short-lines", 1, 2 * MIB // 128, True),
        short_line_fixture(root, "many-small-files", 1024, 4 * 1024 // 128, False),
        long_line_fixture(root, "long-no-match", 48 * MIB, False),
        long_line_fixture(root, "long-one-match", 16 * MIB, True),
    )


def windows_peak_rss(process):
    import ctypes
    from ctypes import wintypes

    class ProcessMemoryCounters(ctypes.Structure):
        _fields_ = [
            ("size", wintypes.DWORD),
            ("page_faults", wintypes.DWORD),
            ("peak_working_set", ctypes.c_size_t),
            ("working_set", ctypes.c_size_t),
            ("peak_paged_pool", ctypes.c_size_t),
            ("paged_pool", ctypes.c_size_t),
            ("peak_nonpaged_pool", ctypes.c_size_t),
            ("nonpaged_pool", ctypes.c_size_t),
            ("pagefile", ctypes.c_size_t),
            ("peak_pagefile", ctypes.c_size_t),
        ]

    counters = ProcessMemoryCounters(size=ctypes.sizeof(ProcessMemoryCounters))
    query = ctypes.WinDLL("psapi", use_last_error=True).GetProcessMemoryInfo
    query.argtypes = (wintypes.HANDLE, ctypes.POINTER(ProcessMemoryCounters), wintypes.DWORD)
    query.restype = wintypes.BOOL
    # CPython retains the child handle after wait, which allows reading the high-water mark.
    if not query(int(process._handle), ctypes.byref(counters), counters.size):
        raise ctypes.WinError(ctypes.get_last_error())
    return counters.peak_working_set


def wait_posix(process):
    def timeout(_signum, _frame):
        raise subprocess.TimeoutExpired(process.args, PROCESS_TIMEOUT_SECONDS)

    previous_handler = signal.signal(signal.SIGALRM, timeout)
    try:
        signal.setitimer(signal.ITIMER_REAL, PROCESS_TIMEOUT_SECONDS)
        return os.wait4(process.pid, 0)
    finally:
        signal.setitimer(signal.ITIMER_REAL, 0)
        signal.signal(signal.SIGALRM, previous_handler)


def timed_process(command, environment, expected_exit):
    start = perf_counter_ns()
    with subprocess.Popen(command, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL,
                          env=environment) as process:
        try:
            if sys.platform == "win32":
                process.wait(timeout=PROCESS_TIMEOUT_SECONDS)
                elapsed = perf_counter_ns() - start
                peak_rss = windows_peak_rss(process)
            else:
                _, status, usage = wait_posix(process)
                elapsed = perf_counter_ns() - start
                process.returncode = os.waitstatus_to_exitcode(status)
                peak_rss = usage.ru_maxrss * (1 if sys.platform == "darwin" else 1024)
        except subprocess.TimeoutExpired:
            process.kill()
            process.wait()
            raise
    if process.returncode != expected_exit or peak_rss <= 0:
        raise RuntimeError(f"{command[0]} returned {process.returncode}, expected "
                           f"{expected_exit}, peak RSS reported: {peak_rss}")
    return {"elapsed_nanoseconds": elapsed, "peak_rss_bytes": peak_rss}


def sha256_file(path):
    with path.open("rb") as file:
        return hashlib.file_digest(file, "sha256").hexdigest()


def assert_expected_output(command, environment, workload, workdir):
    stdout = workdir / "stdout"
    stderr = workdir / "stderr"
    with stdout.open("wb") as output, stderr.open("wb") as errors:
        completed = subprocess.run(command, stdout=output, stderr=errors, env=environment,
                                   timeout=PROCESS_TIMEOUT_SECONDS, check=False)
    if (completed.returncode != workload.exit_code
            or stderr.stat().st_size
            or stdout.stat().st_size != workload.stdout_bytes
            or sha256_file(stdout) != workload.stdout_sha256):
        raise RuntimeError(f"{command[0]} changed the {workload.name} fixture output "
                           f"(exit {completed.returncode}, stdout {stdout.stat().st_size} bytes, "
                           f"stderr {stderr.stat().st_size} bytes)")


def summarize(trials, input_bytes):
    elapsed = median(trial["elapsed_nanoseconds"] for trial in trials)
    return {
        "runs": trials,
        "median_elapsed_ms": round(elapsed / 1_000_000, 3),
        "median_input_mib_per_second": round(input_bytes / MIB / (elapsed / 1_000_000_000), 3),
        "median_peak_rss_bytes": median(trial["peak_rss_bytes"] for trial in trials),
        "max_peak_rss_bytes": max(trial["peak_rss_bytes"] for trial in trials),
    }


def run_workload(workload, binaries, environment, workdir, runs):
    commands = {
        label: [str(binary), "search", NEEDLE, str(workload.root), "--ext", ".txt", "--json"]
        for label, binary in binaries.items()
    }
    for command in commands.values():
        assert_expected_output(command, environment, workload, workdir)
        timed_process(command, environment, workload.exit_code)
    measurements = {label: [] for label in commands}
    order = tuple(commands)
    for iteration in range(runs):
        for label in order if iteration % 2 == 0 else reversed(order):
            measurements[label].append(timed_process(commands[label], environment, workload.exit_code))
    results = {label: summarize(trials, workload.input_bytes) for label, trials in measurements.items()}
    for label, result in results.items():
        peak = result["median_peak_rss_bytes"] / MIB
        throughput = result["median_input_mib_per_second"]
        print(f"{workload.name:20s} {label:9s} {result['median_elapsed_ms']:8.3f} ms "
              f"{throughput:8.1f} MiB/s {peak:7.2f} MiB peak RSS", flush=True)
    return {
        "name": workload.name,
        "input_bytes": workload.input_bytes,
        "expected_exit": workload.exit_code,
        "expected_records": workload.records,
        "expected_stdout_bytes": workload.stdout_bytes,
        "expected_stdout_sha256": workload.stdout_sha256,
        "builds": results,
    }


def commit_sha(value):
    if len(value) != 40 or any(character not in "0123456789abcdef" for character in value):
        raise argparse.ArgumentTypeError("expected a 40-character lowercase Git commit SHA")
    return value


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--baseline", type=Path, required=True)
    parser.add_argument("--baseline-commit", type=commit_sha, required=True)
    parser.add_argument("--candidate", type=Path, required=True)
    parser.add_argument("--candidate-commit", type=commit_sha, required=True)
    parser.add_argument("--runs", type=int, default=5)
    parser.add_argument("--output", type=Path, help="Write the full results to this JSON file.")
    args = parser.parse_args()
    if sys.version_info < (3, 11) or platform.python_implementation() != "CPython":
        parser.error("CPython 3.11 or newer is required")
    if sys.platform not in ("win32", "darwin", "linux"):
        parser.error("peak RSS measurement is implemented for Windows, macOS, and Linux")
    if args.runs < 1:
        parser.error("--runs must be positive")
    binaries = {label: getattr(args, label).resolve(strict=True) for label in ("baseline", "candidate")}
    evidence = {
        "environment": {
            "os": platform.system(),
            "os_release": platform.release(),
            "os_version": platform.version(),
            "architecture": platform.machine(),
            "processor": platform.processor(),
            "logical_processors": os.cpu_count(),
            "python": platform.python_version(),
            "go_runtime_overrides": {name: os.environ[name] for name in
                                      ("GOMAXPROCS", "GOGC", "GOMEMLIMIT") if name in os.environ},
        },
        "method": {
            "runs_per_build": args.runs,
            "warmups_per_build": 1,
            "timeout_seconds_per_process": PROCESS_TIMEOUT_SECONDS,
            "stdout_validation": "One separate digest, length, stderr and exit check per build and workload.",
            "time": "Elapsed wall time from before process creation through exit. Both output streams "
                    "are sent to the null device while timing. Build order alternates each iteration.",
            "cache": "Generated locally, verified and warmed before timing. OS caches are shared.",
            "peak_rss": ("Win32 GetProcessMemoryInfo PeakWorkingSetSize in bytes" if sys.platform == "win32"
                         else "wait4 ru_maxrss converted to bytes"),
        },
        "builds": {
            label: {
                "commit": getattr(args, f"{label}_commit"),
                "binary_sha256": sha256_file(binary),
                "version": subprocess.check_output([str(binary), "--version"], text=True, timeout=5).strip(),
            }
            for label, binary in binaries.items()
        },
    }
    BENCH_DIR.mkdir(parents=True, exist_ok=True)
    with TemporaryDirectory(prefix="run-", dir=BENCH_DIR) as directory:
        workdir = Path(directory)
        temp_dir = workdir / "temp"
        temp_dir.mkdir()
        environment = {**os.environ, "TEMP": str(temp_dir), "TMP": str(temp_dir), "TMPDIR": str(temp_dir)}
        workloads = make_workloads(workdir)
        evidence["workloads"] = [
            run_workload(workload, binaries, environment, workdir, args.runs) for workload in workloads
        ]
    if args.output:
        args.output.write_text(json.dumps(evidence, indent=2) + "\n", encoding="utf-8", newline="\n")


if __name__ == "__main__":
    main()

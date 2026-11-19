# A02 streaming search benchmarks

On November 19, 2026, the Windows amd64 build of `0.1.1-rc.1` reduced median peak working set by approximately 87% on both the sparse 96 MiB fixture and the 48 MiB single-line no-match fixture. The candidate had a lower median peak working set in all five measured workloads. End-to-end time varied with the input and number of matches. [Raw measurements and SHA-256 values](BENCHMARKS_WINDOWS.json)

## Builds and method

- Baseline `0.1.0` uses the Windows amd64 executable from the official private release at source commit `02472793619d8a644634068759ecb2e826a39797`. Executable SHA-256 `a84e51fff68fec0e450e5e2ca763d4ebf48d8595c281705b99a1427a94e598a5`. [Private release](https://github.com/soroushneyestani/developer-cli-toolkit/releases/tag/v0.1.0)
- Candidate `0.1.1-rc.1` was built locally from exact code commit `aaf9a9c05b9a93fb77a517215c99dddc93a26df3` using Go `1.27.4`, `CGO_ENABLED=0`, `-trimpath` and `-buildvcs=false`. Executable SHA-256 `c11ecb7fef161c1ae0933a8dff5da2862bd8dd55cd7b365852cda61b27dce56b`.
- Both builds run `dct search DCT_BENCHMARK_MATCH <fixture> --ext .txt --json`. The harness creates deterministic fixtures and separately verifies the expected exit, empty stderr, and exact NDJSON length and SHA-256 before timing. Each build gets one warmup and five measured invocations per fixture. The build order alternates and both builds use the same temporary directory.
- Wall time includes child creation, traversal, matching and output encoding. The timing pass sends stdout and stderr to the null device and uses warmed filesystem caches. Effective MiB/s divides fixture bytes by median wall time. Windows memory is the peak working set of the individual child process throughout its lifetime.
- The measured host was Windows 11 amd64 build `10.0.26300`, with 16 logical processors, CPython `3.13.1`, and default Go runtime settings. Linux and macOS performance measurements remain open. The test summary records the separate integration-test scope. [Harness](../scripts/benchmark.py) · [Test summary](TEST_SUMMARY.md)

## Median results

| Synthetic fixture | Input | Records | v0.1.0 time (ms) | RC time (ms) | v0.1.0 MiB/s | RC MiB/s | v0.1.0 peak (MiB) | RC peak (MiB) |
| --- | ---: | ---: | ---: | ---: | ---: | ---: | ---: | ---: |
| Sparse short lines, 4 files × 24 MiB | 96 MiB | 4 | 67.403 | 56.175 | 1424.3 | 1709.0 | 57.24 | 7.57 |
| Dense short lines, every line matches | 2 MiB | 16,384 | 28.731 | 29.326 | 69.6 | 68.2 | 13.09 | 11.70 |
| Many small files, 1,024 files × 4 KiB | 4 MiB | 1,024 | 90.100 | 79.352 | 44.4 | 50.4 | 12.52 | 9.22 |
| One unterminated line without a match | 48 MiB | 0 | 41.872 | 50.269 | 1146.4 | 954.9 | 56.31 | 7.32 |
| One long CRLF line with a match near the end | 16 MiB + 21 bytes | 1 | 40.292 | 69.589 | 397.1 | 229.9 | 56.93 | 41.28 |

The sparse and small-file fixtures took less time in this run. Dense results can spill to temporary storage and the dense case took slightly longer. The oversized no-match probe kept the 48 MiB line out of memory and increased elapsed time by about 20%. The oversized matching line required a verified positional read, temporary staging and a complete rendered record. Its peak working set decreased by about 27% and elapsed time increased by about 73%. These figures describe the specified synthetic inputs and measured host. [Streaming ADR](adr/0003-streaming-search.md) · [Raw trials](BENCHMARKS_WINDOWS.json)

## Reproduce the comparison

Extract `dct.exe` from the private `dct-windows-amd64.zip` for v0.1.0 after verifying the archive against its release manifest. In a clean source checkout of `aaf9a9c05b9a93fb77a517215c99dddc93a26df3`, build the candidate with Go `1.27.4` on Windows amd64. [Release manifest](releases/v0.1.0-SHA256SUMS.txt)

```powershell
$env:CGO_ENABLED = '0'
go build -mod=readonly -trimpath -buildvcs=false -o dct-rc.exe .
```

From the review checkout containing the benchmark driver, point it to the two executables. CPython 3.11 or newer is required. The driver writes disposable fixtures under the ignored `bin/.benchmark/` directory. The driver can collect per-child `wait4` peak RSS on Linux and macOS when given binaries for those platforms.

```powershell
python -B scripts/benchmark.py --baseline path/to/v0.1.0/dct.exe `
  --baseline-commit 02472793619d8a644634068759ecb2e826a39797 `
  --candidate path/to/candidate/dct-rc.exe `
  --candidate-commit aaf9a9c05b9a93fb77a517215c99dddc93a26df3 `
  --runs 5 --output bin/.benchmark/reproduced.json
```

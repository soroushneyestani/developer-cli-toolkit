# Changelog

## 0.1.1-rc.1 — Streaming search review

- Process files in a reusable read buffer while preserving the literal CLI contract.
- Stage matches per file and spill large result sets to temporary storage. Emit results only after complete UTF-8 and NUL checks.
- Scan oversized lines in fragments. When staging an oversized match, reread its range from the same open file and verify the copied bytes against the scan's SHA-256.
- Cover boundary-spanning literals, long lines, late invalid bytes, read errors, and spill cleanup with regression tests.
- Measure search time and peak process memory against the private v0.1.0 release using reproducible synthetic workloads. Results and tradeoffs are recorded in `docs/BENCHMARKS.md` in the source repository.
- Prepare the project MIT license, contributor and security guides, and source distribution. New binary bundles require the project license, third-party notices and license texts.

## 0.1.0 — Private release

- Search UTF-8 text recursively for a case-sensitive literal.
- Inspect regular files and report their relative paths and byte sizes.
- Filter either command by repeatable suffixes and emit newline-delimited JSON.
- Return documented exit codes for matches, empty results, and failures.
- Build archives for Windows amd64, Linux amd64, and macOS arm64 with SHA-256 checksums and dependency notices.
- Keep Go sources in LF format during Git checkout so native Windows formatting verification sees the same source bytes as Linux and macOS.

Native CI passed on Windows, Linux, and macOS. The packaging job uploaded all three target archives, and a Windows executable extracted from the verified hosted package passed a synthetic CLI smoke check. See `docs/RELEASE_EVIDENCE.md` in the source repository for the exact commit, verification coverage, and review items.

# Changelog

## 0.1.0 — Release candidate

- Search UTF-8 text recursively for a case-sensitive literal.
- Inspect regular files and report their relative paths and byte sizes.
- Filter either command by repeatable suffixes and emit newline-delimited JSON.
- Return documented exit codes for matches, empty results, and failures.
- Build archives for Windows amd64, Linux amd64, and macOS arm64 with SHA-256 checksums and dependency notices.
- Keep Go sources in LF format during Git checkout so native Windows formatting verification sees the same source bytes as Linux and macOS.

Native CI passed on Windows, Linux, and macOS. The packaging job uploaded all three target archives, and a Windows executable extracted from the verified hosted package passed a synthetic CLI smoke check. See `docs/RELEASE_EVIDENCE.md` in the source repository for the exact commit, verification coverage, and review items.

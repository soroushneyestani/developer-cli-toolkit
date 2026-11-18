# Evidence

## 2026-11-15 — Initial contract and verification gate

Bounded task: define the original v0.1 interface and executable verification gate.
Read AGENTS.md, README.md, TASKS.md, EVIDENCE.md, PROJECT_STATE.json, git history/status, and all existing file paths. Baseline was clean at 1dbf562, with no code or unfinished implementation.

Changes: docs/SPEC.md defines commands, filters, JSON schemas, error codes, safety rules, and synthetic cross-platform integration cases. ADR 0001 retains Go/Cobra. scripts/verify.ps1 requires source/tests and runs formatting, vet, tests, and native build once bootstrapped.

Observed checks:
- PowerShell parser: PASS, no syntax errors in scripts/verify.ps1.
- PROJECT_STATE.json parsing: PASS; status in_progress, verified false.
- git diff --check: PASS (line-ending normalization warnings only).
- powershell.exe -NoProfile -ExecutionPolicy Bypass -File scripts/verify.ps1: exit 1, expected diagnostic: Go is required on PATH. This verifies the missing-tool failure path only.
- pwsh -NoProfile -File scripts/verify.ps1: not run successfully; installed launcher cannot be accessed by the system.
- Go formatting/vet/tests/build, native integration, packaging: NOT RUN. No Go toolchain on PATH; no implementation/module exists yet.

Exact next action: install Go with go and gofmt on PATH, reopen the shell, then implement the Go/Cobra literal-search slice and contract/misuse tests. Use Windows PowerShell to invoke verification if pwsh remains inaccessible. Tool installation was not attempted outside the authorized repository.

Git ownership mismatch was handled using per-command safe.directory configuration scoped to this repository; no global configuration changed. No push or external target testing occurred.

## 2026-11-15 - Bootstrap prerequisite recheck

Attempted next task: Go/Cobra literal-search bootstrap. Read all required files, SPEC, ADR, verification script, history and status. Baseline was clean at 170dc70 with no unfinished implementation.

- Get-Command found neither go nor gofmt on PATH.
- Windows PowerShell invocation of scripts/verify.ps1: exit 1, diagnostic: Go is required on PATH. Install Go, reopen the shell, and rerun scripts/verify.ps1.
- Formatting, vet, tests, native build, integration and packaging: NOT RUN; toolchain unavailable and no module exists.

Human action: install Go with go and gofmt on PATH, reopen the shell and confirm go version and gofmt -h resolve. No installation outside the repository was attempted. Status remains in_progress; v0.1 is not runnable or verified. Existing substantive verification checks were preserved.

## 2026-11-18 — Runnable CLI and Windows contract tests

- Go and gofmt resolved on PATH. `go version` reported `go1.27.4 windows/amd64`. No install was required. The `pwsh` launcher remains inaccessible, so local script checks use Windows PowerShell.
- Added search and inspect, suffix filters, NDJSON, deterministic sorting, scoped read-only traversal, and the documented exit statuses.
- `go test -count=1 -v ./...` passed on Windows (3.608 s). File and directory symlink fixtures ran and passed. The POSIX chmod and colon/tab filename fixtures skipped with explicit reasons.
- `go test -mod=readonly -run '^TestSearchReportsWindowsSharingViolation$' -v` passed (1.902 s). A Windows sharing lock prevented reading a fixture. The process returned exit 2 with a diagnostic after emitting the prior match.
- `go vet -mod=readonly ./...` and `go build -mod=readonly -buildvcs=false -trimpath -o bin/dct.exe .` passed. The executable printed `dct version 0.1.0`.
- Cross-compilation, workflow lint and native Linux/macOS checks remain separate release evidence.

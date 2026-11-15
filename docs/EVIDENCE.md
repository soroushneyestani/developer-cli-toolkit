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

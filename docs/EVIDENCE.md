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

## 2026-11-18 — Release packaging and verification

- Windows PowerShell `scripts/verify.ps1` exited 0. Formatting, module integrity, vet, process-level contract tests and native build passed. The Go test package reported 4.961 s. The Windows sharing-violation test passed, and two POSIX-specific fixtures skipped. An isolated unformatted Go fixture caused the verifier to fail with exit 1.
- The README search and inspection examples passed on a synthetic workspace.
- `scripts/package.ps1` built CGO-disabled Linux amd64, Windows amd64 and macOS arm64 archives. Independently checked ZIP CRCs, SHA-256 digests, target binary headers, release notes and dependency notices. The extracted Windows binary produced the expected search result. Invalid build flags produced a nonzero exit and left prior release directories intact.
- Local archives: `dist/dct-20261118T114738Z-993774766c8e4d81bedd4c47b3feb139/`, with digests in `SHA256SUMS.txt`.
- `actionlint v1.7.12` found zero GitHub Actions workflow errors. ShellCheck and Pyflakes were unavailable. The workflow invokes checked-in PowerShell scripts.
- `govulncheck v1.10.0` exited 0 with no vulnerabilities found. `gitleaks v8.30.1` exited 0 with no leaks across four commits and the staged release changes. Scans ran locally against this repository.
- `go test -c -mod=readonly -buildvcs=false` compiled the Linux amd64 and macOS arm64 test binaries separately.
- Native Linux/macOS tests and hosted CI were not run. WSL access was denied in this session and the repository has no remote. An authorized private feature-branch matrix run or native Linux and macOS hosts are needed before release review. Cross-compilation demonstrates buildability for those targets.

## 2026-11-18 — Native CI checkout consistency

- The first [private hosted matrix](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44305452488) passed on Linux and macOS. Windows reported all five Go sources as unformatted and packaging was skipped.
- A local index checkout with `core.autocrlf=true` reproduced the same five `gofmt` findings. After declaring LF checkout for Go files, the same fixture produced zero formatting findings. The native Windows verifier passed.
- The corrected hosted Windows and packaging jobs await the next feature-branch run.

## 2026-11-19 — Release-review verification

- GitHub returned a private Draft PR with head `482ebacf9d2dcc940c1e6931652b49610dee797a`. [Run 44307143679](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44307143679) returned success for native Windows, Linux, macOS, and packaging. The Windows CRLF checkout blocker is resolved.
- The Windows-hosted delimiter and chmod-based unreadable-file fixtures skipped. The sharing-violation case accounts for Windows read-error coverage. The Linux dependency and history scanners completed successfully.
- Downloaded artifact `14428662773` and matched its reported SHA-256. All three target ZIP hashes matched `SHA256SUMS.txt`. ZIP CRCs, binary headers, build settings, and included notices passed. The executable payloads match the earlier local archives.
- Ran the hosted Windows package's executable from ignored local scratch against synthetic fixtures. Version, filtered search, inspection, no-match exit 1, and malformed suffix exit 2 matched their contracts. The README and education examples use verified outputs.
- Completed release review materials and recorded package and platform limits in [release evidence](RELEASE_EVIDENCE.md).
- The first synchronization preflight found that an educational `.go` fixture under a visible `bin/` subdirectory entered `go list ./...`. The tutorial now generates valid Go text under a dot-prefixed scratch directory. The revised example and the Windows verifier passed together.

## 2026-11-19 — Release-evidence audit

- Authenticated GitHub Actions evidence confirmed the successful native Windows, Linux, macOS, and packaging jobs for candidate source `1d41dcc0184a7fad5e09e1e7c9ee6cefc37dd930`. The Linux logs reported no vulnerabilities from `govulncheck` and no leaks from `gitleaks`. The two platform-specific Windows skips are recorded in the release evidence. [Private workflow run 44523796632](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44523796632)
- The authenticated download for artifact `14528472864` matched the retained local outer archive byte for byte. SHA-256 and CRC verification passed for the outer ZIP and three target ZIPs. The retained manifest matches the outer ZIP member, and each target digest matches that manifest. The three candidate executable payloads match their historical equivalents byte for byte.
- Ran the Windows executable extracted from artifact `14528472864` against synthetic fixtures. It reported `dct version 0.1.0` and produced the expected README and Persian tutorial NDJSON with inspection sizes of 24 and 29 bytes. The no-match case returned 1 and the malformed suffix case returned 2. Linux and macOS archive payloads were checked for integrity and executable format. Native execution on those platforms was exercised by CI's native builds.
- The canonical artifact, retained local paths, and checksum manifest are now aligned for `release_review`. Final acceptance is pending. [Release evidence](RELEASE_EVIDENCE.md) · [Pinned manifest](RELEASE_SHA256SUMS.txt)

## 2026-11-19 — Streaming engineering candidate

- Native Windows verification passed locally. Eighteen synthetic differential cases preserved the v0.1.0 text and JSON output, status and absence of spilled files after exit.
- CI for `541aff83f8675b5d6aac9c98a28adcc1228fb052` passed on all three native platforms and uploaded the candidate package. Checked its outer and per-platform hashes and ran the direct packaged Windows smoke. [Exact CI run 44540116482](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44540116482) · [Verified build artifact 14537687323](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44540116482/artifacts/14537687323)
- Recorded the representative working-set and elapsed-time changes and their tradeoffs separately from the product contract. [Benchmark evidence](BENCHMARKS.md) · [Release evidence](RELEASE_EVIDENCE.md)

## 2026-11-19 — Open-source preparation

- Rechecked the native Windows verifier with Go 1.27.4. Seventy-two additional seeded comparisons against the released binary preserved stdout, stderr and exit status across short and oversized lines, invalid bytes, and dense output.
- Confirmed the 50 recorded benchmark samples and both benchmark executable digests. The three local preparation ZIPs retained the previously reviewed runtime bytes and passed CRC, SHA-256, release-material and machine-header checks after adding required license packaging. [Benchmark evidence](BENCHMARKS.md) · [Release procedure](releases/OPEN_SOURCE_RELEASE.md)
- Compared the historical Git refs with the private remote, verified the dependency notices, and screened the existing history and hosted records. The owner's decisions about rights, historical Git author metadata and the older release remain open. [Publication checklist](../FINAL_PUBLICATION_CHECKLIST.md)

## 2026-11-19 — Public stable v0.1.1

- The tagged main commit `2623811347200ff5433332479ac0e2b75e3fb958` passed all four checks in run `44582995235`. Verified platform, license, CRC and SHA-256 checks for the three binary archives and confirmed both public source archives against their exact tag tree. [Public release](https://github.com/soroushneyestani/developer-cli-toolkit/releases/tag/v0.1.1) · [Stable build evidence](releases/v0.1.1.md)

## 2026-11-19 — A02 closure

The owner accepted public v0.1.1. Required macOS, Linux, Windows and packaging checks passed on the exact pre-closure main commit `3c4c10adca36efc47c56aabed72b4ef7f7835b8f`. The five public assets and both source archives matched the release record. [Completion record](COMPLETION.md) · [Exact post-merge CI](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44587543545)

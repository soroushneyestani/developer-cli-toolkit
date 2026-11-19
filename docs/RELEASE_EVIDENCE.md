# A02 Developer CLI Toolkit release evidence

## Candidate and source

- Version `0.1.0`. The candidate archive was built from commit `1d41dcc0184a7fad5e09e1e7c9ee6cefc37dd930` on `codex/a02`. The pull request CI run started on November 19, 2026 at 13:25:48 CET and recorded success at 13:27:29 CET. [Private workflow run 44523796632](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44523796632)
- Subsequent documentation-only CI runs verify their own source revisions. The candidate package remains pinned to run `44523796632` and artifact `14528472864` so readers can independently verify the exact archive reviewed here.
- The repository is private and the pull request targeting `main` is still Draft. [Draft pull request #1](https://github.com/soroushneyestani/developer-cli-toolkit/pull/1)

## What was verified

| Gate | Result and scope |
| --- | --- |
| Windows native `windows-2025` | Passed formatting, `go mod verify`, `go vet`, process-level integration tests, and build. The Windows filename-delimiter and chmod-based unreadable-file fixtures skipped. The Windows-specific sharing-violation test covers a read failure. |
| Linux native `ubuntu-24.04` | Passed formatting, module integrity, vet, integration tests, and build. The dependency and Git history scans ran in this job. |
| macOS native `macos-15` | Passed formatting, module integrity, vet, integration tests, and build. |
| Dependency and history scans | On Linux, `govulncheck v1.10.0` reported no vulnerabilities and `gitleaks v8.30.1` reported no leaks. Findings are scoped to the workflow's scans and commit. |
| Packaging | Passed and uploaded three target ZIPs with `SHA256SUMS.txt`. Go `1.27.4` used `CGO_ENABLED=0`, `-trimpath`, and `-buildvcs=false`. |
| Hosted artifact integrity | The downloaded outer ZIP matched GitHub's SHA-256 and passed its ZIP CRC check. All three target ZIPs passed CRC checks and matched the bundled manifest. The target executable headers, embedded build flags, candidate changelog, notices and licenses were checked. |
| Packaged Windows execution | The executable extracted from candidate artifact `14528472864` reported `dct version 0.1.0`. A no-match search returned 1 with empty stdout and stderr. A malformed suffix returned 2 with a diagnostic on stderr. |
| Documented examples | The packaged candidate Windows executable produced the documented README, Persian tutorial, and portfolio article NDJSON records. It reported inspection sizes of 24 and 29 bytes and passed the Persian repeated-suffix search across mixed LF and CRLF input. |

The native integration tests build and run executables on their own runners. The separate package build produces the distribution archives. Linux and macOS archive payloads were checked for integrity and architecture during review. Direct execution of those extracted payloads was outside this verification. [Test history](EVIDENCE.md) · [CLI contract](SPEC.md)

## Verified packages

GitHub artifact `14528472864` is **8,982,978 bytes**. Its outer `dct-packages.zip` SHA-256 is `2261c04556de25eb3d793d1b1da5ebc716171ccf450ec944d841d20652ef7fa1`. GitHub listed expiration on November 26, 2026 at 13:27:23 CET (`2026-11-26T12:27:23Z`). [Private dct-packages artifact](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44523796632/artifacts/14528472864)

| ZIP | Bytes | Verified binary type |
| --- | ---: | --- |
| `dct-linux-amd64.zip` | 3,006,276 | ELF64 x86-64 |
| `dct-windows-amd64.zip` | 3,096,446 | PE32+ AMD64 |
| `dct-darwin-arm64.zip` | 2,888,079 | Mach-O 64-bit ARM64 |

The verified outer ZIP, the three target ZIPs, and their original `SHA256SUMS.txt` are retained privately in `dist/hosted-44523796632/` in the original workspace. This copy remains available locally when the GitHub Actions download expires. The checked-in manifest records the three target ZIP digests for this artifact. [Pinned SHA-256 manifest](RELEASE_SHA256SUMS.txt)

## Historical verification

The implementation baseline is `482ebacf9d2dcc940c1e6931652b49610dee797a`. Its native matrix and package job passed on November 18, 2026 in run `44307143679`, resolving the Windows CRLF checkout blocker. The dependency and history scanners reported no vulnerabilities and no leaks in that run. [Private workflow run 44307143679](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44307143679)

Historical artifact `14428662773` is 8,982,494 bytes with outer SHA-256 `1630089350a96cc38169b634a7a5a5cf4ab8407c2b85f7cdf2c7df028b7c1d46`. Its listed expiration was November 25, 2026 at 13:49:02 CET. The corresponding local archives and manifest remain under `dist/hosted-44307143679/`. Each package set has its own ZIP digests. The three candidate executable payloads match the corresponding historical executable payloads byte for byte. [Historical private artifact](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44307143679/artifacts/14428662773)

## Reproduce the checks

From the project directory with Go `1.27.4` and PowerShell 7 available, run the real verifier and packaging script as needed.

```powershell
pwsh -NoProfile -File scripts/verify.ps1
pwsh -NoProfile -File scripts/package.ps1
```

To verify the candidate download or retained copy, check `dct-packages.zip` with `Get-FileHash -Algorithm SHA256` against the outer digest above. Extract the bundle and verify the three target ZIPs using the manifest for artifact `14528472864`. The packaging script creates a new archive set with its own manifest. [README verification and extraction steps](../README.md#install-the-010-candidate)

## Release review items

- Product scope is recursive, literal, case-sensitive UTF-8 search and regular-file inspection. Hidden entries follow normal selection. Regex, stdin input, ignore-file processing, and numeric performance targets are outside the candidate scope.
- The selected path list and a full search file affect memory use. Keep the live file tree stable for repeatable output. A read or write error can follow earlier results. Consumers should check the exit code before using all records as complete output.
- Unix JSON paths containing invalid UTF-8 use replacement characters. Linux and macOS archives may extract with mode `0644`; set the execute bit when needed.
- The packages are unsigned and include third-party notices and licenses. SHA-256 verifies package bytes against the pinned digests. Signing, public licensing and wider distribution remain separate review decisions.
- The project remains in `release_review` pending explicit user acceptance. Tagging `v0.1.0`, creating a private GitHub Release with the verified assets and checksums, and any PR merge require approval. [Draft pull request #1](https://github.com/soroushneyestani/developer-cli-toolkit/pull/1)
- The review includes the user documentation and release notes. [Persian tutorial](LEARNING_FA.md) · [English website article](WEBSITE_EN.md) · [English Instagram copy](INSTAGRAM_EN.md) · [Changelog](../CHANGELOG.md)

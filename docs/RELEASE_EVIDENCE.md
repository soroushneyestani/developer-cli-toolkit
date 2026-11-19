# A02 Developer CLI Toolkit release evidence

## Stable v0.1.1

The stable version uses streaming search with the reviewed literal-search and inspection contract. The public release supplies Windows amd64, Linux amd64 and macOS arm64 packages and tagged source archives. Its SHA-256 manifest applies to that exact package set. The stable tag resolves to `2623811347200ff5433332479ac0e2b75e3fb958`. Exact main CI run `44582995235` passed and produced artifact `14559337068`. The checked-in proof records anonymous verification of five assets and both source archives. Older RC and historical records below retain their own provenance. [Stable provenance](releases/v0.1.1.md) · [Stable SHA-256 manifest](releases/v0.1.1-SHA256SUMS.txt) [Stable release](https://github.com/soroushneyestani/developer-cli-toolkit/releases/tag/v0.1.1) · [Changelog](../CHANGELOG.md)

## Final closure

The owner accepted public `v0.1.1` on 2026-11-19. The final documentation PR #3 merge at `3c4c10adca36efc47c56aabed72b4ef7f7835b8f` passed macOS, Linux, Windows and packaging in run `44587543545`. The public binaries, checksums, source archives, MIT license, and documentation were checked again for closure. [Completion record](COMPLETION.md) · [Post-merge CI](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44587543545)

## Open-source preparation

PR #2 added the project MIT license, retained the existing dependency licenses, and required them in new binary packages. All three local preparation archives passed CRC, checksum, license and machine-header checks and contained executables identical to the pinned RC. The PR's preparation summary records the exact hosted asset and test runs. [Publication checklist](../FINAL_PUBLICATION_CHECKLIST.md) · [Release procedure](releases/OPEN_SOURCE_RELEASE.md) · [Merged PR #2](https://github.com/soroushneyestani/developer-cli-toolkit/pull/2)

## 0.1.1-rc.1 engineering candidate

The pinned candidate was built by the private branch push workflow from commit `541aff83f8675b5d6aac9c98a28adcc1228fb052` on `codex/a02-streaming`. Its code is at checkpoint `aaf9a9c05b9a93fb77a517215c99dddc93a26df3`. Native Windows, Linux, and macOS verification and packaging passed on both the branch push and Draft PR runs. The subsequent evidence checkpoint keeps the runtime, tests, and packaging inputs from the pinned source. [Draft PR #2](https://github.com/soroushneyestani/developer-cli-toolkit/pull/2) · [Exact CI run 44540116482](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44540116482) · [Verified build artifact 14537687323](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44540116482/artifacts/14537687323)

The authenticated download's outer ZIP is 9,474,410 bytes with SHA-256 `53afa2a5d787f023210e39a127d4bc2da12c86913efec3909bc7b80b798d6677`. Its three platform archives passed CRC and SHA-256 checks and contained the expected executable targets and licenses. The manifest's SHA-256 is `f770371ec7a65240881ffa898745426bbbc42765ebe64e3151ecd6e63066d45a`. The outer ZIP, platform archives, and manifest are retained privately under `dist/hosted-44540116482/` in the workspace. [Candidate SHA-256 manifest](releases/v0.1.1-rc.1-SHA256SUMS.txt) · [0.1.1-rc.1 build evidence](releases/v0.1.1-rc.1-build.json)

The packaged Windows binary passed a direct smoke test with 5,002 ordered NDJSON records, dense temporary staging, late NUL and truncated UTF-8 suppression, boundary-spanning Unicode, CRLF and final CR, exits 0 through 2, and temporary-file cleanup. The packaged Windows executable matches the benchmarked candidate binary byte for byte. Linux and macOS native CI executed the CLI built on each runner. Distribution packages for those two targets were inspected for integrity and build settings. [Test summary](TEST_SUMMARY.md)

The benchmark compares the released v0.1.0 binary with an exact streaming checkpoint build. It verifies output equivalence and measures time and peak working set on generated Windows workloads. This candidate keeps the ordinary validation buffers bounded and still materializes each complete matching line for output. Dense output can require space in the temporary directory. [Benchmark results](BENCHMARKS.md) · [Streaming ADR](adr/0003-streaming-search.md)

## Approved v0.1.0 release

GitHub PR #1 was merged with both original parents preserved. The historical `v0.1.0` release points to the accepted commit `02472793619d8a644634068759ecb2e826a39797`. It was originally published privately. Its uploaded packages came from the exact main workflow run `44536777289`, artifact `14536722607`. Every uploaded asset matched its verified local copy byte for byte. [Historical release](https://github.com/soroushneyestani/developer-cli-toolkit/releases/tag/v0.1.0) · [Source and asset verification](releases/v0.1.0.md) · [Historical MIT grant](releases/v0.1.0-MIT.md)

## Archived v0.1.0 review evidence

The following sections record the artifact seen during review. Its checksums apply to those archived bytes. [Published asset manifest](releases/v0.1.0-SHA256SUMS.txt)

### Candidate and source

- Version `0.1.0`. The candidate archive was built from commit `1d41dcc0184a7fad5e09e1e7c9ee6cefc37dd930` on `codex/a02`. The pull request CI run started on November 19, 2026 at 13:25:48 CET and recorded success at 13:27:29 CET. [Private workflow run 44523796632](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44523796632)
- Subsequent documentation-only CI runs verify their own source revisions. The candidate package remains pinned to run `44523796632` and artifact `14528472864` so readers can independently verify the exact archive reviewed here.
- The repository was private and PR #1 was Draft during this verification. The approved PR has since been merged. [Pull request #1](https://github.com/soroushneyestani/developer-cli-toolkit/pull/1)

### What was verified

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

### Verified packages

GitHub artifact `14528472864` is **8,982,978 bytes**. Its outer `dct-packages.zip` SHA-256 is `2261c04556de25eb3d793d1b1da5ebc716171ccf450ec944d841d20652ef7fa1`. GitHub listed expiration on November 26, 2026 at 13:27:23 CET (`2026-11-26T12:27:23Z`). [Private dct-packages artifact](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44523796632/artifacts/14528472864)

| ZIP | Bytes | Verified binary type |
| --- | ---: | --- |
| `dct-linux-amd64.zip` | 3,006,276 | ELF64 x86-64 |
| `dct-windows-amd64.zip` | 3,096,446 | PE32+ AMD64 |
| `dct-darwin-arm64.zip` | 2,888,079 | Mach-O 64-bit ARM64 |

The verified outer ZIP, the three target ZIPs, and their original `SHA256SUMS.txt` are retained privately in `dist/hosted-44523796632/` in the original workspace. This copy remains available locally when the GitHub Actions download expires. The checked-in manifest records the three target ZIP digests for this artifact. [Pinned SHA-256 manifest](RELEASE_SHA256SUMS.txt)

### Historical verification

The implementation baseline is `482ebacf9d2dcc940c1e6931652b49610dee797a`. Its native matrix and package job passed on November 18, 2026 in run `44307143679`, resolving the Windows CRLF checkout blocker. The dependency and history scanners reported no vulnerabilities and no leaks in that run. [Private workflow run 44307143679](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44307143679)

Historical artifact `14428662773` is 8,982,494 bytes with outer SHA-256 `1630089350a96cc38169b634a7a5a5cf4ab8407c2b85f7cdf2c7df028b7c1d46`. Its listed expiration was November 25, 2026 at 13:49:02 CET. The corresponding local archives and manifest remain under `dist/hosted-44307143679/`. Each package set has its own ZIP digests. The three candidate executable payloads match the corresponding historical executable payloads byte for byte. [Historical private artifact](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44307143679/artifacts/14428662773)

### Reproduce the checks

From the project directory with Go `1.27.4` and PowerShell 7 available, run the real verifier and packaging script as needed.

```powershell
pwsh -NoProfile -File scripts/verify.ps1
pwsh -NoProfile -File scripts/package.ps1
```

To verify the candidate download or retained copy, check `dct-packages.zip` with `Get-FileHash -Algorithm SHA256` against the outer digest above. Extract the bundle and verify the three target ZIPs using the manifest for artifact `14528472864`. The packaging script creates a new archive set with its own manifest. [Pinned review manifest](RELEASE_SHA256SUMS.txt)

### Release review items

- Product scope is recursive, literal, case-sensitive UTF-8 search and regular-file inspection. Hidden entries follow normal selection. Regex, stdin input, ignore-file processing, and numeric performance targets are outside the candidate scope.
- The selected path list and a full search file affect memory use. Keep the live file tree stable for repeatable output. A read or write error can follow earlier results. Consumers should check the exit code before using all records as complete output.
- Unix JSON paths containing invalid UTF-8 use replacement characters. Linux and macOS archives may extract with mode `0644`; set the execute bit when needed.
- The packages are unsigned and include third-party notices and licenses. SHA-256 verifies package bytes against the pinned digests. Signing, public licensing and wider distribution remain separate review decisions.
- The review gate for these materials was satisfied by the user's explicit v0.1.0 approval. The accepted merge and published binaries are recorded separately. [Published v0.1.0 evidence](releases/v0.1.0.md)
- The review includes the user documentation and release notes. [Persian tutorial](LEARNING_FA.md) · [English website article](WEBSITE_EN.md) · [English Instagram copy](INSTAGRAM_EN.md) · [Changelog](../CHANGELOG.md)

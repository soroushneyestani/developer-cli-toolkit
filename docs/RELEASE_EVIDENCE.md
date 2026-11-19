# A02 Developer CLI Toolkit release evidence

## Candidate and source

- Version `0.1.0`. Implementation baseline `482ebacf9d2dcc940c1e6931652b49610dee797a` on `codex/a02`. The release-review follow-up changes documentation and project state.
- Native verification and package creation completed on November 18, 2026 between 13:46:34 and 13:49:09 CET. Candidate verification and archive checks were reviewed on November 19. [Private workflow run 44307143679](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44307143679)
- Private review target `main`. [Draft pull request #1](https://github.com/soroushneyestani/developer-cli-toolkit/pull/1)

## What was verified

| Gate | Result and scope |
| --- | --- |
| Windows native `windows-2025` | Passed formatting, `go mod verify`, `go vet`, process-level integration tests, and build. The Windows filename-delimiter and chmod-based unreadable-file fixtures skipped. The Windows-specific sharing-violation test covers a read failure. |
| Linux native `ubuntu-24.04` | Passed formatting, module integrity, vet, integration tests, and build. The dependency and Git history scans ran in this job. |
| macOS native `macos-15` | Passed formatting, module integrity, vet, integration tests, and build. |
| Dependency and history scans | `govulncheck v1.10.0` reported no vulnerabilities found. `gitleaks v8.30.1` reported no leaks found. Findings are scoped to the workflow's scan and commit. |
| Packaging | Passed and uploaded three target ZIPs with `SHA256SUMS.txt`. Go `1.27.4` used `CGO_ENABLED=0`, `-trimpath`, and `-buildvcs=false`. |
| Hosted artifact integrity | The downloaded outer ZIP matched GitHub's SHA-256 and passed its ZIP CRC check. The three inner archive digests matched the bundled manifest. Inner CRCs, target executable headers, embedded build flags, changelog, and third-party notices and licenses were checked. |
| Packaged Windows execution | The executable from the verified Windows ZIP reported `dct version 0.1.0`. A synthetic `TODO` search and an inspection produced the expected NDJSON. No match returned 1 with empty stdout. A malformed suffix returned 2 with a diagnostic on stderr. |
| Documented examples | The Persian examples passed on an existing Windows executable, including mixed LF and CRLF, repeated suffixes, inspection byte counts, and misuse. The website example returned its stated NDJSON record. |

The native integration tests build and run executables on their own runners. The separate package build produces the distribution archives. Linux and macOS archive payloads were checked for integrity and architecture during review. Direct execution of those extracted payloads was outside this verification. See [test history](EVIDENCE.md) and the [CLI contract](SPEC.md) for details.

## Verified packages

GitHub artifact `14428662773` is **8,982,494 bytes**. Its outer archive SHA-256 is `1630089350a96cc38169b634a7a5a5cf4ab8407c2b85f7cdf2c7df028b7c1d46`. GitHub listed expiration on November 25, 2026 at 13:49:02 CET. The verified copy and the three inner ZIPs are retained locally in `dist/hosted-44307143679/`. [Private dct-packages artifact](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44307143679/artifacts/14428662773)

| ZIP | Bytes | Verified binary type |
| --- | ---: | --- |
| `dct-linux-amd64.zip` | 3,006,114 | ELF64 x86-64 |
| `dct-windows-amd64.zip` | 3,096,284 | PE32+ AMD64 |
| `dct-darwin-arm64.zip` | 2,887,917 | Mach-O 64-bit ARM64 |

[Pinned SHA-256 manifest](RELEASE_SHA256SUMS.txt) belongs to artifact `14428662773`. Each package set has its own manifest. The earlier local archives differ in ZIP bytes, and their executable payloads match the hosted executables byte for byte.

## Reproduce the checks

From the project directory with Go `1.27.4` and PowerShell 7 available, run the real verifier and packaging script as needed.

```powershell
pwsh -NoProfile -File scripts/verify.ps1
pwsh -NoProfile -File scripts/package.ps1
```

Use `Get-FileHash -Algorithm SHA256` on the outer archive and each target ZIP to compare with the digest and manifest for that exact artifact. The [README](../README.md) provides a manifest check and the Windows extraction steps.

## Release review items

- Product scope is recursive, literal, case-sensitive UTF-8 search and regular-file inspection. Hidden entries follow normal selection. Regex, stdin input, ignore-file processing, and numeric performance targets are outside the candidate scope.
- The selected path list and a full search file affect memory use. Keep the live file tree stable for repeatable output. A read or write error can follow earlier results. Consumers should check the exit code before using all records as complete output.
- Unix JSON paths containing invalid UTF-8 use replacement characters. Linux and macOS archives may extract with mode `0644`; set the execute bit when needed.
- SHA-256 verifies package bytes against the selected manifest. Public code licensing, signing, and wider distribution are separate review decisions. The packages already carry third-party notices.
- The candidate is ready for private release review and user acceptance. Review the [Persian tutorial](LEARNING_FA.md), [English website article](WEBSITE_EN.md), and [English Instagram copy](INSTAGRAM_EN.md) with the [changelog](../CHANGELOG.md).

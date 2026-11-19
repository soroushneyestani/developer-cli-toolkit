# A02 v0.1.0 test summary

The binary candidate was built and verified at `1d41dcc0184a7fad5e09e1e7c9ee6cefc37dd930` on `codex/a02`. All three native verification jobs and the subsequent packaging job passed. The candidate is pinned to this run and artifact for release review. [Private candidate CI run 44523796632](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44523796632) · [Private artifact 14528472864](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44523796632/artifacts/14528472864)

## Native verification

| GitHub Actions runner | Result | Distribution package and binary type |
| --- | --- | --- |
| Windows `windows-2025` | Passed formatting, module integrity, vet, process-level integration tests, and native build. | Windows amd64, PE32+ AMD64 |
| Linux `ubuntu-24.04` | Passed the same native checks and the Linux-only dependency and Git history scans. | Linux amd64, ELF64 x86-64 |
| macOS `macos-15` | Passed formatting, module integrity, vet, process-level integration tests, and native build. | macOS arm64, Mach-O 64-bit ARM64 |

The native tests compile and run the CLI against synthetic fixtures. They cover deterministic ordering, NDJSON, repeated suffix filters, Unicode and line endings, long lines, binary and invalid UTF-8 text handling, malformed arguments, symlinks, and read-error behavior. The candidate Windows log shows two fixture skips, `inspect_JSON_escapes_filename_delimiters` and `unreadable_file_is_an_error`, due to Windows filename restrictions and the permissions available to the test. The Windows sharing-violation test passed and exercises a read failure on Windows. The reviewed Linux and macOS verification logs showed no test skip lines. [Search tests](../search_test.go) · [Inspection and options tests](../options_test.go) · [Windows sharing-violation test](../search_windows_test.go)

The candidate Linux scan steps finished successfully using `govulncheck v1.10.0` and `gitleaks v8.30.1`. These are scoped to the workflow's dependency and checked-out Git history checks. [CI workflow](../.github/workflows/ci.yml)

## Verified candidate artifact

The retained `dct-packages.zip` for artifact `14528472864` is 8,982,978 bytes. Its SHA-256 is `2261c04556de25eb3d793d1b1da5ebc716171ccf450ec944d841d20652ef7fa1`. The outer ZIP and all three inner ZIPs passed CRC checks. Every inner ZIP matched the bundled `SHA256SUMS.txt`, whose SHA-256 is `032ce1fe00289796f0f1b831b62c3c987e816e97d93b00589f918638daea0afa`. The archives include executables, release notes, and third-party notices and licenses. [Pinned checksum manifest](RELEASE_SHA256SUMS.txt) · [Release evidence](RELEASE_EVIDENCE.md)

| Target ZIP | Bytes | SHA-256 |
| --- | ---: | --- |
| `dct-windows-amd64.zip` | 3,096,446 | `87d0b52217a506bccce9647ecf5df7ffdf072a0714db4ec68427187172517de8` |
| `dct-linux-amd64.zip` | 3,006,276 | `f6a324433f7ab0e3c66f565e08fdb5c72b7933d13c2b2d05493cc31bbd131dc6` |
| `dct-darwin-arm64.zip` | 2,888,079 | `a7a8a26fbfa557480396b9628f4c11c9446ad6584491d8565c1edf136333b84d` |

Local review used the retained private copy in `dist/hosted-44523796632/`. The extracted review package carries the same artifact bytes in `artifacts/dct-packages.zip`.

The Windows executable was extracted directly from the verified candidate archive. It reported `dct version 0.1.0`, reproduced the README's filtered search NDJSON and 24-byte inspection, matched the Persian guide's ordered search NDJSON and 29-byte inspection, and produced the portfolio article's search record. A no-match search returned exit 1 with empty stdout. A malformed suffix returned exit 2 with a diagnostic on stderr. [README examples](../README.md) · [Persian guide](LEARNING_FA.md) · [Portfolio article](WEBSITE_EN.md)

The Linux and macOS distribution executables were checked for archive integrity and target architecture. Direct execution of those extracted distribution executables remains outside this review. Native CLI execution on both platforms was covered by the integration tests.

## Historical verification

Run `44307143679` at implementation checkpoint `482ebacf9d2dcc940c1e6931652b49610dee797a` passed the same runner matrix and packaging. The historical Windows archive's executable passed a separate smoke test. The executable payloads in the candidate and historical packages were compared byte for byte and matched for all three targets. [Historical CI run 44307143679](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44307143679) · [Verification history](EVIDENCE.md)

Historical artifact `14428662773` had outer SHA-256 `1630089350a96cc38169b634a7a5a5cf4ab8407c2b85f7cdf2c7df028b7c1d46`. Its ZIP digests apply to that artifact only.

## Outstanding limitations

- Search is literal, case-sensitive, and operates on UTF-8 text. Hidden entries participate in selection. Regex, stdin input, and ignore-file processing are outside v0.1.0.
- Memory grows with the selected path list and the largest search file. Filesystem changes during execution can affect results. A read or write error can leave earlier records on stdout before exit 2.
- Unix paths containing invalid UTF-8 use replacement characters in JSON. Linux and macOS executables may require `chmod +x dct` after extraction.
- Evidence covers the listed runner and package targets. The review contains no measured performance targets or package signing.

See the [CLI specification](SPEC.md) for the complete command and exit-status contract.

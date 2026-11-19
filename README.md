# Developer CLI Toolkit

`dct` searches workspace text for a literal and lists regular files with their byte sizes. Both commands produce sorted, relative paths and can emit newline-delimited JSON (NDJSON) for scripts.

The original Go/Cobra implementation takes inspiration from recursive search and file filters. `dct` includes hidden files, uses case-sensitive literal matching, and applies explicit suffix filters. [CLI contract](docs/SPEC.md) · [ripgrep reference](https://github.com/BurntSushi/ripgrep)

The source and release review are currently private. The approved baseline is `v0.1.0`. The follow-up `0.1.1-rc.1` adds streaming and is being prepared for open-source publication. [Review index](REVIEW_INDEX.md) · [Publication checklist](FINAL_PUBLICATION_CHECKLIST.md)

## Install v0.1.0

Download the ZIPs and `SHA256SUMS.txt` from the private GitHub Release. The ZIPs were built by the verified workflow for the tagged `main` commit. Compare the downloaded manifest with the versioned copy and verify the ZIPs before extracting the package for your platform. The packages are unsigned. [Private v0.1.0 release](https://github.com/soroushneyestani/developer-cli-toolkit/releases/tag/v0.1.0) · [Release provenance](docs/releases/v0.1.0.md) · [SHA-256 manifest](docs/releases/v0.1.0-SHA256SUMS.txt)

| Target | Archive | Executable |
| --- | --- | --- |
| Windows amd64 | `dct-windows-amd64.zip` | `dct.exe` |
| Linux amd64 | `dct-linux-amd64.zip` | `dct` |
| macOS arm64 | `dct-darwin-arm64.zip` | `dct` |

```powershell
Get-Content SHA256SUMS.txt | ForEach-Object {
    $expected, $name = $_ -split '\s+', 2
    $actual = (Get-FileHash -LiteralPath $name -Algorithm SHA256).Hash
    if ($actual -ine $expected) { throw "Checksum mismatch for $name" }
}
Expand-Archive -LiteralPath dct-windows-amd64.zip -DestinationPath dct-windows
.\dct-windows\dct.exe --version
```

On Linux, verify with `sha256sum -c SHA256SUMS.txt`. On macOS, use `shasum -a 256 -c SHA256SUMS.txt`. Extract the target ZIP, set its run permission with `chmod +x dct` if needed, and run `./dct --version`. The packaged executable runs without a local Go installation.

## Build and verify

The complete source lives in the repository. Clone it, or extract the source ZIP or tarball for the revision you want to build. Private review requires repository access. To build the streaming candidate, check out the reviewed commit in the engineering review. [Source repository](https://github.com/soroushneyestani/developer-cli-toolkit) · [Engineering review](REVIEW_INDEX.md)

```sh
git clone https://github.com/soroushneyestani/developer-cli-toolkit.git
cd developer-cli-toolkit
```

Use Go 1.26 or newer. Go 1.27.4 matches the verified CI and packaging toolchain. On Linux and macOS, run `go build -mod=readonly -trimpath -o dct .` and `./dct --version`. On Windows, run `go build -mod=readonly -trimpath -o dct.exe .` and `.\dct.exe --version`. A clone defaults to `main`. Check the reported version when switching revisions.

```powershell
powershell.exe -NoProfile -ExecutionPolicy Bypass -File scripts/verify.ps1
```

The verifier checks formatting, module integrity, vet, tests, and the native build. It writes `bin/dct.exe` on Windows and `bin/dct` on Linux and macOS. Use `pwsh -NoProfile -File scripts/verify.ps1` with PowerShell 7.

Git checks out Go sources with LF endings on every platform so `gofmt` sees the same source bytes on Windows and Unix.

## Usage

```text
dct search <literal> [root] [--ext <suffix>]... [--json]
dct inspect [root] [--ext <suffix>]... [--json]
dct --help
dct --version
```

`root` defaults to the current directory. Put any flags before `--` when passing a literal or a path that starts with `-`. Repeat `--ext` to include files matching any supplied suffix. Suffixes are case-sensitive, start with a dot, and may contain multiple dots.

Search emits one record per matching line and skips files containing NUL or invalid UTF-8. Inspect includes all regular files that pass the suffix filter. Both commands skip symlink entries and reject symlink roots.

| Exit code | Meaning |
| --- | --- |
| 0 | Search found matches, or inspection succeeded. |
| 1 | Search completed with no matches. |
| 2 | Invalid invocation, filesystem failure, or output failure. Diagnostics go to stderr. |

Check the exit code before treating output as complete. Errors may follow partial results.

## Sample workflow

In PowerShell from the project directory, create a synthetic workspace and run the built binary.

```powershell
$demo = Join-Path (Get-Location).Path ('bin/dct-demo-' + [guid]::NewGuid().ToString('N'))
New-Item -ItemType Directory -Path $demo | Out-Null
[IO.File]::WriteAllBytes((Join-Path $demo 'notes.txt'), [Text.Encoding]::ASCII.GetBytes("TODO: review the change`n"))
.\bin\dct.exe search TODO $demo --ext .txt --json
# {"path":"notes.txt","line":1,"text":"TODO: review the change"}
.\bin\dct.exe inspect $demo --json
# {"path":"notes.txt","bytes":24}
```

Use NDJSON when paths or contents can contain output delimiters. Read it one line at a time and decode each record independently.

## Packaging

Run `pwsh -NoProfile -File scripts/package.ps1` with Go 1.27.4. On Windows PowerShell, use `powershell.exe -NoProfile -ExecutionPolicy Bypass -File scripts/package.ps1`. A successful run creates a directory under `dist/` with ZIP archives for Linux amd64, Windows amd64, and macOS arm64, plus `SHA256SUMS.txt`. Verify each archive's SHA-256 before extraction. On a Unix host, run `chmod +x dct` if the extracted binary lacks the execute bit.

All three binary packages produced from this branch include the project license and dependency license texts. CI runs native integration tests and builds on Windows, Linux and macOS before uploading packages. [CI workflow](.github/workflows/ci.yml) · [Release evidence](docs/RELEASE_EVIDENCE.md) · [Release procedure](docs/releases/OPEN_SOURCE_RELEASE.md)

## Architecture and limits

`main.go` parses commands with Cobra, renders output and maps errors to exit codes. `workspace.go` scopes filesystem access to the opened directory and sorts relative paths. The search engine reads files in chunks and stages results until each file reaches a clean, valid EOF. Oversized nonmatching lines are scanned without accumulating their text. When staging an oversized match, the engine rereads its range from the same open file and verifies the bytes by SHA-256. Dense match sets can use temporary disk. Memory still grows with the selected paths and the size of a matching record during output. [Streaming design](docs/adr/0003-streaming-search.md) · [Measured benchmark cases](docs/BENCHMARKS.md)

Results describe a live filesystem, so keep the tree stable for reproducible output. On Unix, filenames with invalid UTF-8 bytes are represented with replacement characters in JSON. Search takes its literal and root from command arguments. Regex matching, stdin search, and ignore-file rules are outside the CLI contract.

## License and participation

The MIT license covers the original source and documentation in this review branch. The Go toolchain and libraries retain their respective licenses. Historical `v0.1.0` assets keep their recorded provenance. [MIT license](LICENSE) · [Third-party notices](THIRD_PARTY_NOTICES.md)

See the [contribution guide](CONTRIBUTING.md), [security policy](SECURITY.md), [Persian tutorial](docs/LEARNING_FA.md), [technical article](docs/WEBSITE_EN.md), and [test evidence](docs/EVIDENCE.md).

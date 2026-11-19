# Developer CLI Toolkit

`dct` searches workspace text for a literal and lists regular files with their byte sizes. Both commands produce sorted, relative paths and can emit newline-delimited JSON (NDJSON) for scripts.

The original Go/Cobra implementation takes inspiration from recursive search and file filters. `dct` includes hidden files, uses case-sensitive literal matching, and applies explicit suffix filters. [CLI contract](docs/SPEC.md) · [ripgrep reference](https://github.com/BurntSushi/ripgrep)

Version `v0.1.1` adds streaming search and keeps the existing command and exit-code contract. The original project source and documentation are available under the MIT license. [Changelog](CHANGELOG.md) · [MIT license](LICENSE)

## Install v0.1.1

Download the ZIP for your platform and `SHA256SUMS.txt` from the same release. Verify the ZIP before extracting it. The packages are unsigned. [v0.1.1 release assets](https://github.com/soroushneyestani/developer-cli-toolkit/releases/tag/v0.1.1) · [Release evidence](docs/RELEASE_EVIDENCE.md)

| Target | Archive | Executable |
| --- | --- | --- |
| Windows amd64 | `dct-windows-amd64.zip` | `dct.exe` |
| Linux amd64 | `dct-linux-amd64.zip` | `dct` |
| macOS arm64 | `dct-darwin-arm64.zip` | `dct` |

```powershell
$archive = 'dct-windows-amd64.zip'
$rows = @(Get-Content -LiteralPath SHA256SUMS.txt | Where-Object {
    ($_ -split '\s+', 2)[1] -ceq $archive
})
if ($rows.Count -ne 1) { throw "Missing or ambiguous checksum for $archive" }
$expected = ($rows[0] -split '\s+', 2)[0]
$actual = (Get-FileHash -LiteralPath $archive -Algorithm SHA256).Hash
if ($actual -ine $expected) { throw "Checksum mismatch for $archive" }
Expand-Archive -LiteralPath $archive -DestinationPath dct-windows
.\dct-windows\dct.exe --version
```

On Linux, run `sha256sum dct-linux-amd64.zip`. On macOS, run `shasum -a 256 dct-darwin-arm64.zip`. Compare the result with the matching line in `SHA256SUMS.txt`. Extract the target ZIP, set its run permission with `chmod +x dct` if needed, and run `./dct --version`. The packaged executable reports `dct version 0.1.1` and runs without a local Go installation.

## Build and verify

Clone the repository and check out the stable tag, or extract its source ZIP or tarball. [Source repository](https://github.com/soroushneyestani/developer-cli-toolkit) · [Source archives](https://github.com/soroushneyestani/developer-cli-toolkit/releases/tag/v0.1.1)

```sh
git clone https://github.com/soroushneyestani/developer-cli-toolkit.git
cd developer-cli-toolkit
git switch --detach v0.1.1
```

Use Go 1.26 or newer. Go 1.27.4 matches the CI and packaging toolchain. On Linux and macOS, run `go build -mod=readonly -trimpath -o dct .` and `./dct --version`. On Windows, run `go build -mod=readonly -trimpath -o dct.exe .` and `.\dct.exe --version`.

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
| 2 | Invalid invocation, filesystem, temporary-storage, or output failure. Diagnostics go to stderr. |

Check the exit code before treating output as complete. Errors may follow partial results.

## Sample workflow

After running the verifier in Windows PowerShell, create a synthetic workspace and run the resulting binary.

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

All three `v0.1.1` binary packages include the project license and dependency license texts. CI runs native integration tests on Windows, Linux and macOS before building the release packages. [CI workflow](.github/workflows/ci.yml) · [Release evidence](docs/RELEASE_EVIDENCE.md) · [Release procedure](docs/releases/OPEN_SOURCE_RELEASE.md)

## Architecture and limits

`main.go` parses commands with Cobra, renders output and maps errors to exit codes. `workspace.go` scopes filesystem access to the opened directory and sorts relative paths. The search engine reads files in chunks and stages results until each file reaches a clean, valid EOF. Oversized nonmatching lines are scanned without accumulating their text. When staging an oversized match, the engine rereads its range from the same open file and verifies the bytes by SHA-256. Dense match sets can use temporary disk. Memory still grows with the selected paths and the size of a matching record during output. [Streaming design](docs/adr/0003-streaming-search.md) · [Measured benchmark cases](docs/BENCHMARKS.md)

Results describe a live filesystem, so keep the tree stable for reproducible output. On Unix, filenames with invalid UTF-8 bytes are represented with replacement characters in JSON. Search takes its literal and root from command arguments. Regex matching, stdin search, and ignore-file rules are outside the CLI contract.

## License and participation

The MIT license covers the original source and documentation. The Go toolchain and libraries retain their respective licenses. The historical `v0.1.0` archives contain the earlier project tree and dependency notices. Consult that release's notes for its licensing and provenance. [MIT license](LICENSE) · [Third-party notices](THIRD_PARTY_NOTICES.md) · [Historical v0.1.0 release](https://github.com/soroushneyestani/developer-cli-toolkit/releases/tag/v0.1.0)

See the [contribution guide](CONTRIBUTING.md), [security policy](SECURITY.md), [Persian tutorial](docs/LEARNING_FA.md), [technical article](docs/WEBSITE_EN.md), and [test evidence](docs/EVIDENCE.md).

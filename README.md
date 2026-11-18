# Developer CLI Toolkit

`dct` searches workspace text for a literal and lists regular files with their byte sizes. Both commands produce sorted, relative paths and can emit newline-delimited JSON (NDJSON) for scripts.

The original Go/Cobra implementation takes inspiration from the recursive search and file filters in [ripgrep](https://github.com/BurntSushi/ripgrep). `dct` includes hidden files, uses case-sensitive literal matching, and applies only explicit suffix filters. Its scope and performance goals are defined by the [v0.1 contract](docs/SPEC.md).

## Build and verify

Use Go 1.26 or newer. Use Go 1.27.4 to match the verified Windows toolchain.

```powershell
go build -mod=readonly -buildvcs=false -trimpath -o bin/dct.exe .
powershell.exe -NoProfile -ExecutionPolicy Bypass -File scripts/verify.ps1
```

On Linux or macOS, build to `bin/dct` and run verification with `pwsh -NoProfile -File scripts/verify.ps1`.

## Usage

```text
dct search <literal> [root] [--ext <suffix>]... [--json]
dct inspect [root] [--ext <suffix>]... [--json]
dct --help
dct --version
```

`root` defaults to the current directory. Pass `--` before a literal or a path that starts with `-`. Repeat `--ext` to include files matching any supplied suffix. Suffixes are case-sensitive, start with a dot, and may contain multiple dots.

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
New-Item -ItemType Directory -Force bin/demo | Out-Null
Set-Content -Encoding ascii bin/demo/notes.txt 'TODO: review the change'
.\bin\dct.exe search TODO bin/demo --ext .txt --json
# {"path":"notes.txt","line":1,"text":"TODO: review the change"}
.\bin\dct.exe inspect bin/demo --json
```

Use NDJSON when paths or contents can contain output delimiters. Read it one line at a time and decode each record independently.

## Architecture and limits

`main.go` parses commands with Cobra, renders output and maps errors to exit codes. `workspace.go` scopes filesystem access to the opened directory, sorts relative paths, and checks each search file before matching its lines. Memory usage grows with the number of selected paths and the size of the largest searched file. Results describe a live filesystem, so keep the tree stable for reproducible output. On Unix, filenames with invalid UTF-8 bytes are represented with replacement characters in JSON.

See the [design decisions](docs/adr/), [test evidence](docs/EVIDENCE.md), and [third-party notices](THIRD_PARTY_NOTICES.md).

# Developer CLI Toolkit

A workspace command should tell a script when to continue, when a search finished empty, and when to stop after an error. I built Developer CLI Toolkit around that requirement. Its command, `dct`, searches UTF-8 text and inventories regular files through a small Go and Cobra interface.

## Using the toolkit

`dct search` finds a case-sensitive literal recursively and emits each matching line once. `dct inspect` lists regular files and their sizes in bytes. Both commands default to the current directory, accept repeatable `--ext` suffix filters, and can emit newline-delimited JSON with `--json`. A file passes the filter if its name ends in any supplied suffix in a case-sensitive comparison.

That makes the tool useful for finding `TODO` markers before a review, checking the contents of a handoff directory, or feeding a file inventory into a reporting script. For example, with `TODO: review the change` on the first line of `demo/notes.txt`, this command emits one record.

```sh
dct search TODO demo --ext .txt --json
```

```json
{"path":"notes.txt","line":1,"text":"TODO: review the change"}
```

Search records contain `path`, `line`, and `text`. Inspection records contain `path` and `bytes`. Exit code `0` means a search found results or an inspection completed. Exit code `1` means a search completed without a match. Exit code `2` covers invalid input and runtime errors. Diagnostics go to standard error. A script should check the exit code before treating output as complete because a read error can occur after earlier results have been emitted.

## Design decisions

Cobra parses commands and flags in `main.go`, where text and JSON output are also formatted. The workspace layer in `workspace.go` opens the requested directory using Go's `os.OpenRoot`, enumerates matching regular files, and sorts their relative paths. Paths use forward slashes on each platform. Enumeration includes qualifying hidden files and files listed in ignore rules, skips symbolic-link entries, and applies the requested suffix filters. The command rejects a symbolic link supplied as the root.

Search reads one file at a time and validates its whole contents before emitting matches from that file. This matters when a file appears to contain a valid match before an invalid UTF-8 sequence or a NUL byte. The file is skipped in either case. Search handles LF, CRLF, and an unterminated final line. Inspection reads metadata for byte sizes, with no need to open file contents.

The selected paths remain in memory for sorting, and whole-file validation uses memory proportional to the largest file read during the search. A stable workspace gives reproducible results because the filesystem stays live during a scan. Newline-delimited JSON preserves record boundaries when names or matched text contain delimiters. On Unix, filenames containing invalid UTF-8 bytes are represented with replacement characters in JSON.

## Verification and packaging

The release candidate passed native CI on Windows, Linux, and macOS. Integration tests build and invoke the executable against synthetic workspaces, checking standard output, standard error, and exit status separately. The cases cover deterministic ordering, long lines, Unicode, malformed usage, and file handling. Some permission and filename fixtures depend on host capabilities. A Windows-specific test uses a sharing lock to exercise partial output followed by a read error.

The first hosted Windows run exposed a Git checkout issue. CRLF conversion made the Go source files appear unformatted to `gofmt`. I set a repository attribute to keep Go source checkouts in LF, then verified the corrected native matrix.

After the native jobs passed, a separate packaging job built ZIP archives for Windows amd64, Linux amd64, and macOS arm64. Each archive includes the executable and dependency notices. The job also generates a SHA-256 manifest for checking the ZIPs before extraction. The native results cover executables built for the tests, and the packaging job records which archives were produced.

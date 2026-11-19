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

## Streaming search with whole-file validation

Cobra parses commands and flags in `main.go`, where text and JSON output are also formatted. The workspace layer in `workspace.go` opens the requested directory using Go's `os.OpenRoot`, enumerates matching regular files, and sorts their relative paths. Paths use forward slashes on each platform. Enumeration includes qualifying hidden files and files listed in ignore rules, skips symbolic-link entries, and applies the requested suffix filters. The command rejects a symbolic link supplied as the root.

Version `0.1.0` read an entire file into memory before scanning it. I kept the whole-file validation guarantee and changed how the follow-up candidate uses memory. It checks for invalid UTF-8 and NUL through a reusable 64 KiB reader. Matching records stay pending until the source reaches a clean end of file. If invalid data appears near the end, the command discards the file's pending matches, even if the first line looked usable.

The candidate holds up to 256 KiB of encoded matches in memory per file and uses temporary disk when the pending results exceed that budget. It checks long lines in pieces, retaining only enough UTF-8 and literal bytes to handle buffer boundaries. Unmatched long lines require no second read. For a matching long line, the CLI rereads just that range from the same open file and checks its SHA-256 against the original scan before output. LF, CRLF, and unterminated final lines retain their original behavior.

The sorted path list and the longest emitted line still affect memory, and dense matches can incur temporary I/O. Inspection reads file sizes from metadata. A stable workspace gives repeatable results because traversal and reads observe a live filesystem. On Unix, filenames containing invalid UTF-8 bytes are represented with replacement characters in JSON.

## Measuring the tradeoff

I compared the streaming candidate with the released `v0.1.0` binary on five generated Windows workloads, checking output equivalence before measuring. Each median uses five fresh CLI processes. On the 96 MiB sparse fixture, median peak working set fell from 57.24 to 7.57 MiB. The 48 MiB unmatched line went from 56.31 to 7.32 MiB, with a longer median runtime. The 16 MiB matching line also reduced memory and took longer, with median time rising from 40.292 to 69.589 ms. These measurements include CLI startup and output encoding and apply to the recorded inputs and machine. [Benchmark method and full results](BENCHMARKS.md)

## Verification and packaging

The private `v0.1.0` release comes from the accepted main commit. Integration tests build and invoke the CLI on Windows, Linux, and macOS, checking standard output, standard error, and exit status separately. The follow-up candidate adds coverage for validation across buffer boundaries, long lines, and a late invalid byte after matches have already been staged. Some permission and filename fixtures depend on host capabilities. A Windows-specific test uses a sharing lock to exercise a read error after output from an earlier file.

The first hosted Windows run exposed a Git checkout issue. CRLF conversion made the Go source files appear unformatted to `gofmt`. I set a repository attribute to keep Go source checkouts in LF, then verified the corrected native matrix.

For the release, the packaging job built ZIP archives for Windows amd64, Linux amd64, and macOS arm64 after native CI passed. Each archive includes the executable and dependency notices. The release includes a SHA-256 manifest, notes, and known limitations. The streaming changes are a separate `0.1.1-rc.1` candidate for final review. [Private v0.1.0 release](https://github.com/soroushneyestani/developer-cli-toolkit/releases/tag/v0.1.0) · [Release evidence](RELEASE_EVIDENCE.md)

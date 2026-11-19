# v0.1 contract and test plan

Executable: `dct`, an original implementation inspired by ripgrep behavior, without source copying or compatibility promises.

## Interface

- `dct search <literal> [root]`: recursive, case-sensitive UTF-8 literal search. Root defaults to the current directory. Emit each matching line once.
- `dct inspect [root]`: recursively list regular files with relative paths and byte sizes.
- Both accept repeatable `--ext .go` (OR across exact case-sensitive suffixes) and `--json` (newline-delimited JSON).
- `--help` and `--version` succeed without traversing files.

Use relative paths with `/` separators and lexicographic path ordering. Search line numbers start at 1. Support LF, CRLF, and unterminated final lines; exclude terminators from output. Reject empty literals, extra arguments, unknown flags, file roots, and extension filters lacking a leading dot or containing path separators.

Include hidden files. Ignore-file rules do not affect selection. Skip symlink entries during traversal and reject symlink roots. Inspect only regular files. Search reads one file at a time and publishes its records only after the source has been read successfully and validated in full. A NUL byte or invalid UTF-8 anywhere in the file discards every match from it. A source read failure returns an error even after invalid content has been detected. Keep the workspace stable while a command runs because traversal and reading observe a live tree.

## Output

Search text: `path:line:text`. Inspect text: `path<TAB>bytes`. Text delimiters can be ambiguous in filenames; JSON is the machine interface.

Search JSON fields: `path` (string), `line` (integer), `text` (string). Inspect JSON fields: `path` (string), `bytes` (integer). Emit diagnostics only to stderr. Empty results emit no records.

- Exit 0 means a search has matches or inspection succeeds, including in an empty directory.
- Exit 1 means a completed search has no matches.
- Exit 2 means usage, root, traversal, read, temporary storage, or output failed. Consumers must check exit status because errors can follow records from earlier files or records already written during output.

An empty root argument is invalid. On Unix, NDJSON represents any invalid UTF-8 in filenames with replacement characters.

## Integration proof

Build and invoke the binary using Go os/exec and synthetic temporary fixtures, without shell-specific commands. Assert stdout, stderr, and exit status separately.

1. Search: nested files, repeated matches on one line, case sensitivity, Unicode, CRLF, unterminated lines, no matches, deterministic order.
2. Filters and misuse: repeated extensions, excluded files, malformed filters, empty literals, unknown commands/flags, extra arguments, nonexistent and file roots.
3. Inspect: nested/empty directories, byte sizes, hidden files, spaces and delimiters in filenames, JSON escaping.
4. Safety: local symlinks, NUL files, invalid UTF-8, unreadable files. Explicitly skip permission/symlink tests with a reason when the host cannot enforce/create fixtures. Never test external targets.
5. Packaging: CGO-disabled Linux amd64, Windows amd64, macOS arm64 builds. Cross-compilation proves buildability only; run integration tests natively on all three OS families before release.
6. Streaming regressions cover buffer boundaries, long lines, a matched source range that changes before it can be staged, and suppression when NUL or invalid UTF-8 appears after earlier matches. [Resource limits and validation order in ADR 0003](adr/0003-streaming-search.md)

## Implementation order

The v0.1.0 baseline established the Cobra interface, contract tests, suffix filters, NDJSON, inspection, packaging, and native CI. The follow-up keeps that interface and replaces whole-file allocation with streaming and staged matches. Measure performance on reproducible synthetic datasets against the v0.1.0 baseline. The contract makes no numeric speed promise.

# v0.1 contract and test plan

Executable: `dct`, an original implementation inspired by ripgrep behavior, without source copying or compatibility promises.

## Interface

- `dct search <literal> [root]`: recursive, case-sensitive UTF-8 literal search. Root defaults to the current directory. Emit each matching line once.
- `dct inspect [root]`: recursively list regular files with relative paths and byte sizes.
- Both accept repeatable `--ext .go` (OR across exact case-sensitive suffixes) and `--json` (newline-delimited JSON).
- `--help` and `--version` succeed without traversing files.

Use relative paths with `/` separators and lexicographic path ordering. Search line numbers start at 1. Support LF, CRLF, and unterminated final lines; exclude terminators from output. Reject empty literals, extra arguments, unknown flags, file roots, and extension filters lacking a leading dot or containing path separators.

Include hidden files. Do not interpret ignore files or follow symlinks; reject symlink roots. Inspect only regular files. Search skips files containing NUL or invalid UTF-8. Process one file at a time. Traversal/read failures are errors, not silent nonmatches.

## Output

Search text: `path:line:text`. Inspect text: `path<TAB>bytes`. Text delimiters can be ambiguous in filenames; JSON is the machine interface.

Search JSON fields: `path` (string), `line` (integer), `text` (string). Inspect JSON fields: `path` (string), `bytes` (integer). Emit diagnostics only to stderr. Empty results emit no records.

Exit 0: search has matches, or inspection succeeds (including empty directories). Exit 1: completed search has no matches. Exit 2: usage, root, traversal, or read error. Partial results may precede errors; consumers must check exit status.

## Integration proof

Build and invoke the binary using Go os/exec and synthetic temporary fixtures, without shell-specific commands. Assert stdout, stderr, and exit status separately.

1. Search: nested files, repeated matches on one line, case sensitivity, Unicode, CRLF, unterminated lines, no matches, deterministic order.
2. Filters and misuse: repeated extensions, excluded files, malformed filters, empty literals, unknown commands/flags, extra arguments, nonexistent and file roots.
3. Inspect: nested/empty directories, byte sizes, hidden files, spaces and delimiters in filenames, JSON escaping.
4. Safety: local symlinks, NUL files, invalid UTF-8, unreadable files. Explicitly skip permission/symlink tests with a reason when the host cannot enforce/create fixtures. Never test external targets.
5. Packaging: CGO-disabled Linux amd64, Windows amd64, macOS arm64 builds. Cross-compilation proves buildability only; run integration tests natively on all three OS families before release.

## Implementation order

Cobra skeleton and literal search with contract/misuse tests; extension filters; JSON; inspection; packaging and native CI. Performance benchmarking follows correctness; there is no numeric speed promise for v0.1.

# ADR 0003: Stream source files and stage matches until validation

Status: accepted for the 0.1.1-rc.1 implementation.

## Context

The original implementation read an entire file into memory before looking for matches. Complete validation matters because a match can appear before a NUL byte or an invalid UTF-8 sequence. The CLI contract discards all matches from such a file. On common inputs with short lines, keeping a complete source file in memory also makes process memory depend on file size.

## Decision

Read one source file at a time from the opened workspace. A reusable 64 KiB buffered reader supplies the search. If a logical line fits in the buffer, validate its UTF-8 and NUL rules and match it directly, using the text without its documented terminator. A final line without LF is processed as-is.

Scan a longer line in fixed-size fragments. Keep up to three bytes for an incomplete UTF-8 rune, a pending carriage return for CRLF detection, and overlap proportional to the literal length for matches crossing fragments. Hash the text of the long line while scanning it. An unmatched line requires no full-line allocation, staging, or second read. For a matching line, read the bounded source range again through the same open file handle and confirm its SHA-256 before accepting the staged result. Changed or truncated content can raise a staging error. The workspace still needs to remain stable during traversal and reads.

Stage matching records until the source reaches a clean end of file. Up to 256 KiB of encoded matches per file fit in memory. A larger set of matches spills to a temporary file created on demand. Each record has a 16-byte little-endian header with the line number and text byte length, followed by the line text without its terminator. Once the file is validated, replay the staged records in line order through the usual text or JSON renderer. Only matching records go to staging.

If invalid data follows earlier matches, discard the stage and consume the remaining source bytes so a later read error still produces exit 2. Temporary staging is closed and removed as the operation unwinds. A staging I/O failure can produce exit 2 for an otherwise valid file.

## Consequences

Ordinary scans read the source once. Matching lines that cross the reader capacity require a second read limited to the matched line. Moderately sized result sets need only in-memory staging. Dense results can use temporary disk proportional to their encoded matches and incur an additional write and read. This path requires writable temporary storage. A temporary file can contain input text, so it is created with restricted permissions where supported and removed during ordinary cleanup. Process termination can interrupt cleanup.

During the scan, memory includes the sorted path list, the literal and its boundary overlap, and the fixed-size buffers. Emitting a result still requires memory proportional to the longest matching line. The implementation has no fixed maximum line length. A single very long matching line can require substantial output memory, and a large tree can retain many paths.

This decision preserves the command interface, skip rules, relative path ordering, and exit codes. It supersedes only the whole-file allocation decision. [Previous workspace I/O decision](0002-workspace-io.md) · [CLI contract](../SPEC.md)

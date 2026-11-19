# ADR 0002: Scoped filesystem reads and complete-file validation

Status: accepted

`dct` checks that the supplied root is a directory and rejects a symbolic link at that path. It opens the directory with `os.OpenRoot` to constrain descendant operations to the opened tree. The filesystem walk skips symlink entries and sorts all selected, slash-separated paths because recursive directory order can differ from relative path order.

Search reads and validates one complete file before emitting its matches. A file can contain a matching line before its first NUL or invalid UTF-8 byte. Emitting matches before validation would include records from files the contract requires the command to skip. Complete-file reads trade memory proportional to the largest file for a single read per path, complete binary-file detection, and no fixed line length. A two-pass design would lower memory use and require repeated I/O and consistency checks between passes.

Inspection reads metadata and leaves file contents unopened. Enumeration and reading observe live filesystem state. A concurrent replacement can change which in-root file a path resolves to after enumeration. Keep the input tree stable for a consistent result.

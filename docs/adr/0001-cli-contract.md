# ADR 0001: Go/Cobra and deterministic literal search

Status: accepted

Retain the preferred Go/Cobra stack. Separate traversal/matching from command parsing to support unit tests and process-level contract tests. Pin Go and Cobra versions when bootstrapping go.mod, and commit go.sum with dependencies.

Start with literal UTF-8 search and sequential processing, sorting paths for reproducible output. Sorting consumes memory proportional to discovered paths. Regex, ignore-file compatibility, and parallelism are deferred. No ripgrep performance parity is claimed.

Use standard Go testing and synthetic temporary fixtures. The module requires Go 1.26.0 or newer and pins Cobra 1.10.2. Windows checks use Go 1.27.4.

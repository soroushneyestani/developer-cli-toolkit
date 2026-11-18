# Tasks

- [x] Define original v0.1 contract, stack ADR, and integration test matrix.
- [x] Add verification script with explicit prerequisites and substantive Go checks.
- [x] Bootstrap Go/Cobra and literal search with contract and misuse tests.
- [x] Add extension filters and JSON with tests.
- [x] Add workspace inspection with tests.
- [x] Package binaries and add a native Windows/Linux/macOS CI matrix.
- [x] Document runnable sample workflows, release candidate, and limitations.
- [x] Run the native Linux and macOS jobs and record their results before release review.
- [ ] Verify the fixed Windows checkout and the complete hosted packaging pipeline.

Windows native verification, local release packaging, workflow lint, and security scans passed. Hosted Linux and macOS verified successfully. The first hosted Windows checkout used CRLF for Go sources and failed formatting. The checkout fix awaits the complete CI run.

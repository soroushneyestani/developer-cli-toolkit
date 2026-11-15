# Developer CLI Toolkit

Reference: ripgrep — https://github.com/BurntSushi/ripgrep

Original MVP: Build a fast file-search and workspace-inspection CLI with filters, JSON output and error codes.

Proof: Cross-platform integration tests, binary packaging and sample workflows.

See docs/TASKS.md for work items.


The repository is at the specification stage; no CLI binary exists yet.
See [the CLI contract and test plan](docs/SPEC.md) and [stack decision](docs/adr/0001-cli-contract.md).

Run `pwsh -NoProfile -File scripts/verify.ps1` to check prerequisites, formatting,
Go vet, tests, and native build. Go and gofmt must be on PATH. Verification fails
until the module, implementation, and tests exist. Packaging and native CI remain
future work.
On Windows, `powershell.exe -NoProfile -ExecutionPolicy Bypass -File scripts/verify.ps1`
also works when the pwsh launcher is unavailable.

Latest prerequisite recheck: neither go nor gofmt resolves on PATH. Install Go and reopen the shell before resuming bootstrap.

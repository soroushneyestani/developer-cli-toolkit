# PROJECT 01 — A02 — COMPLETED

Owner acceptance and closure dated 2026-11-19 (Europe/Berlin). The accepted public release is Developer CLI Toolkit `v0.1.1`. [Public release](https://github.com/soroushneyestani/developer-cli-toolkit/releases/tag/v0.1.1) · [MIT license](../LICENSE)

| Checkpoint | Exact reference |
| --- | --- |
| Released source and tag | `2623811347200ff5433332479ac0e2b75e3fb958`, `v0.1.1` |
| Main verified for closure | `3c4c10adca36efc47c56aabed72b4ef7f7835b8f`, the merge of documentation PR #3 |
| Tagged source CI | Successful run `44582995235` |
| Documentation PR #3 CI | Successful run `44584336637` |
| Post-PR #3 main CI | Successful run `44587543545` · attempt `2` |

The release archives were built from the tagged source. The main commit verified for closure includes the release-provenance documentation. The closure documentation follows through a separate protected PR, and the workspace progress and synchronization reports record its resulting main commit. [Tagged main CI](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44582995235) · [Merged PR #3](https://github.com/soroushneyestani/developer-cli-toolkit/pull/3) · [PR checks](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44584336637) · [Post-merge main CI](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44587543545)

## Required post-merge checks

- Native verification (macos-15) passed. [Job](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44587543545/job/135793992208)
- Native verification (ubuntu-24.04) passed. [Job](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44587543545/job/135793996181)
- Native verification (windows-2025) passed. [Job](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44587543545/job/135794087899)
- Package binaries passed. [Job](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44587543545/job/135794200175)

GitHub retried failed jobs after a hosted macOS capacity cancellation. Attempt 2 passed.

## Public release verification

Anonymous downloads of the five published assets matched the release SHA-256 records. The binaries are Windows amd64, Linux amd64 and macOS arm64, with the published manifest and release notes. Both public source archive formats matched all 51 tagged files. GitHub detected the MIT license and the documentation remained accessible. [Stable binary manifest](releases/v0.1.1-SHA256SUMS.txt) · [Release evidence](releases/v0.1.1.md) · [Public documentation](https://github.com/soroushneyestani/developer-cli-toolkit/tree/main/docs)

Main uses the protected PR workflow with the three native checks and packaging, an up-to-date branch, and resolved conversations. The project has one maintainer and zero required approving reviews. Force pushes and branch deletion are blocked. [Public merge record](https://github.com/soroushneyestani/developer-cli-toolkit/pull/3)

## Known limitations

- Search handles case-sensitive literals in UTF-8 files. NUL and invalid UTF-8 invalidate the source file for search. The CLI uses a live filesystem and expects a stable tree during traversal.
- The selected-path list and largest emitted record can grow memory use. Dense output may spill to temporary disk. A long matching line used less peak memory and more elapsed time in the recorded synthetic Windows test. [Benchmark scope and data](BENCHMARKS.md)
- Binaries are unsigned. Native integration tests run on Windows, Linux and macOS. The recorded direct smoke tests of the packaged executable ran on Windows. [Test and platform evidence](RELEASE_EVIDENCE.md)

Project 02 awaits separate explicit authorization.

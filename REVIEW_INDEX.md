# A02 Developer CLI Toolkit engineering review

A02 is `completed` with public `v0.1.1`. The historical `v0.1.0` release is the benchmark baseline. [Completion record](docs/COMPLETION.md) · [Historical v0.1.0 release](https://github.com/soroushneyestani/developer-cli-toolkit/releases/tag/v0.1.0) · [Published asset provenance](docs/releases/v0.1.0.md)

PR #2 merged the open-source release preparation. The checklist records release verification, and the release procedure explains how to publish the source and packages. [Final publication checklist](FINAL_PUBLICATION_CHECKLIST.md) · [Release procedure](docs/releases/OPEN_SOURCE_RELEASE.md) · [PR #2](https://github.com/soroushneyestani/developer-cli-toolkit/pull/2)

The public stable release, exact CI, source and verified assets are recorded separately. [v0.1.1 release evidence](docs/releases/v0.1.1.md) [Final publication decisions](FINAL_PUBLICATION_DECISIONS.md)

## Documentation

- [README](README.md) covers setup, installation, and command examples.
- [Changelog](CHANGELOG.md) covers the v0.1.0 baseline, the streaming candidate, and stable v0.1.1.
- [Persian learning guide](docs/LEARNING_FA.md), [English portfolio article](docs/WEBSITE_EN.md), and [English Instagram copy](docs/INSTAGRAM_EN.md).
- [Project state](PROJECT_STATE.json) and [project tasks](docs/TASKS.md) record the accepted closure.

## Implementation and verification

- [CLI specification](docs/SPEC.md), [CLI contract decision](docs/adr/0001-cli-contract.md), and [streaming search decision](docs/adr/0003-streaming-search.md) explain the implemented behavior.
- [Command entry points](main.go), [workspace traversal and inspection](workspace.go), and [streaming search](stream_search.go) contain the implementation.
- [Search integration tests](search_test.go), [inspection and options integration tests](options_test.go), and [Windows sharing-violation test](search_windows_test.go) exercise the executable.
- [Test summary](docs/TEST_SUMMARY.md), [release evidence](docs/RELEASE_EVIDENCE.md), [benchmark evidence](docs/BENCHMARKS.md), and [verification history](docs/EVIDENCE.md) record the checked results and their scope.
- [Local verifier](scripts/verify.ps1), [packaging script](scripts/package.ps1), and [GitHub Actions workflow](.github/workflows/ci.yml) define the runnable checks.
- [Published v0.1.0 SHA-256 manifest](docs/releases/v0.1.0-SHA256SUMS.txt) lists the three released distribution archive digests. [Third-party notices](THIRD_PARTY_NOTICES.md) identify their redistributed components.

## Verified 0.1.1-rc.1 package

The candidate is pinned to `541aff83f8675b5d6aac9c98a28adcc1228fb052`. Verify the three platform archives against the candidate manifest before extracting a platform package. The retained package path and outer checksum are in the release evidence. [Release evidence](docs/RELEASE_EVIDENCE.md) · [Candidate SHA-256 manifest](docs/releases/v0.1.1-rc.1-SHA256SUMS.txt) · [PR #2](https://github.com/soroushneyestani/developer-cli-toolkit/pull/2) · [Exact CI run 44540116482](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44540116482) · [Verified build artifact 14537687323](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44540116482/artifacts/14537687323)

## Archived v0.1.0 review package

The v0.1.0 review ZIP preserved artifact `14528472864` from run `44523796632`, previously reviewed in PR #1. That source revision and manifest differ from the subsequently published release assets. [Historical review evidence](docs/RELEASE_EVIDENCE.md#archived-v010-review-evidence) · [Merged PR #1](https://github.com/soroushneyestani/developer-cli-toolkit/pull/1)

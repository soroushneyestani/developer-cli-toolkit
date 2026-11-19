# A02 Developer CLI Toolkit v0.1.0 review index

A02 is in `release_review`. The implementation, tests, and review documents follow their repository paths in the review ZIP. Extract the whole archive into one directory and open this file to follow the relative links offline.

## Documentation

- [README](README.md) covers setup, installation, and command examples.
- [Changelog](CHANGELOG.md) describes the v0.1.0 changes.
- [Persian learning guide](docs/LEARNING_FA.md), [English portfolio article](docs/WEBSITE_EN.md), and [English Instagram copy](docs/INSTAGRAM_EN.md) are ready for content review.
- [Project state](PROJECT_STATE.json) records the release review status. [Project tasks](docs/TASKS.md) record the outstanding approval.

## Implementation and verification

- [CLI specification](docs/SPEC.md), [CLI contract decision](docs/adr/0001-cli-contract.md), and [workspace I/O decision](docs/adr/0002-workspace-io.md) explain the implemented behavior.
- [Command entry points](main.go) and [workspace traversal, search, and inspection](workspace.go) contain the implementation.
- [Search integration tests](search_test.go), [inspection and options integration tests](options_test.go), and [Windows sharing-violation test](search_windows_test.go) exercise the executable.
- [Test summary](docs/TEST_SUMMARY.md), [release evidence](docs/RELEASE_EVIDENCE.md), and [verification history](docs/EVIDENCE.md) record the checked results and their scope.
- [Local verifier](scripts/verify.ps1), [packaging script](scripts/package.ps1), and [GitHub Actions workflow](.github/workflows/ci.yml) define the runnable checks.
- [Release SHA-256 manifest](docs/RELEASE_SHA256SUMS.txt) lists the three distribution archive digests. [Third-party notices](THIRD_PARTY_NOTICES.md) identify their redistributed components.

## Release candidate package

The v0.1.0 binary candidate comes from run `44523796632` at `1d41dcc0184a7fad5e09e1e7c9ee6cefc37dd930` on `codex/a02`. The link to the artifact requires repository access and is subject to GitHub Actions retention. [Private candidate CI run](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44523796632) · [Private candidate artifact 14528472864](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44523796632/artifacts/14528472864) · [Private draft PR #1](https://github.com/soroushneyestani/developer-cli-toolkit/pull/1)

In the extracted review ZIP, `artifacts/dct-packages.zip` is the preserved copy of the candidate artifact. Its outer digest is recorded in the [test summary](docs/TEST_SUMMARY.md). Extract it into a separate directory to see the three target ZIPs and their bundled checksum manifest. The review ZIP also includes `SOURCE_OVERVIEW.md`, `RELEASE_CHECKLIST.md`, `REVIEW_PACKAGE_CHANGES.md`, and saved CI metadata and verification output for offline review.

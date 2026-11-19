# A02 final publication checklist

Developer CLI Toolkit is `completed` after the owner accepted public `v0.1.1`. The owner authorized public release of `v0.1.1` in this repository. The historical baseline is `v0.1.0`. [PR #2](https://github.com/soroushneyestani/developer-cli-toolkit/pull/2) · [CI checks](https://github.com/soroushneyestani/developer-cli-toolkit/pull/2/checks)

The final decisions report records the verified address categories, publication paths, and historical asset treatment. [Final publication decisions](FINAL_PUBLICATION_DECISIONS.md)

## Source and license audit

- [x] Matched all three remote branches and `v0.1.0` against the local refs. The pre-preparation history at `f6c3b9727c2fcfbf067de2c6472862818c6bcde7` contains 15 commits, 44 paths and 110 unique tracked blobs. Inspected the original Go source, scripts, fixtures, documentation, attribution and historical file inventory. The tracked blobs are regular UTF-8 text. The audit found no vendored application code or recorded outside contributors. One unreachable local blob was outside the ref history.
- [x] Redacted Gitleaks and focused blob and commit-message screening found no credentials or matching private paths in the reviewed refs. Screened two PRs, their conversations, one release record and 14 of 14 then-retained Actions log archives for credential and personal-data patterns. The profile-path indications matched hosted Windows runner names. Findings are scoped to the inspected records and patterns.
- [x] Prepared the standard MIT notice for original source and documentation, naming Soroush Neyestani based on recorded authorship and repository ownership. The existing third-party license texts remain in place. Verified their text against the pinned module and Go distributions. The three build dependency closures include Cobra and pflag, Go's runtime and standard library with its bundled DNS package, and mousetrap on Windows. [MIT license](LICENSE) · [Standard MIT text](https://spdx.org/licenses/MIT.html) · [Third-party notices](THIRD_PARTY_NOTICES.md)
- [x] The owner confirmed this is a personal project and authorized MIT for the original material for which they hold the necessary rights.
- [x] The owner accepted public disclosure of their own historical Git author and committer metadata. Fourteen commits include a non-GitHub-no-reply address. Its value is excluded from this checklist. The existing commits and tag preserve their original hashes. [GitHub commit email guidance](https://docs.github.com/en/account-and-profile/how-tos/email-preferences/setting-your-commit-email-address)

## CI and release candidate

- [x] Verified the `f6c3b972` pre-preparation head on both push and Draft PR CI. The Windows, Linux and macOS jobs and packaging passed. The local Windows verifier passed again with the expected filename and POSIX-permission fixture skips. The sharing-violation regression passed. [Push CI](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44540885418) · [PR CI](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44540899232) · [Test summary](docs/TEST_SUMMARY.md)
- [x] Reviewed full-file validation, oversized line scanning, positional rereads, match replay, spill cleanup and filesystem assumptions. Seventy-two seeded Windows comparisons against the released binary matched stdout, stderr and exit status, including buffer boundaries, UTF-8 and NUL cases and dense spill. The review found no reproducible runtime defect.
- [x] Verify the exact stable preparation revision, merge commit, and release evidence. [Stable evidence](docs/releases/v0.1.1.md) [PR checks](https://github.com/soroushneyestani/developer-cli-toolkit/pull/2/checks)

## Benchmark results

- [x] Recomputed the 50 recorded process measurements across five generated Windows workloads. The sparse 96 MiB case used 7.57 MiB median peak working set versus 57.24 MiB in `v0.1.0`, and fell from 67.403 to 56.175 ms. The unmatched 48 MiB line used 7.32 MiB versus 56.31 MiB, taking 50.269 ms versus 41.872 ms. The matching 16 MiB line used 41.28 MiB versus 56.93 MiB and took 69.589 ms versus 40.292 ms. These are measured inputs on the recorded Windows host. [Reproduction and full results](docs/BENCHMARKS.md) · [Raw samples](docs/BENCHMARKS_WINDOWS.json)

## Release assets and complete source

- [x] Preserved the private `v0.1.0` tag, all five existing release assets, and the earlier pinned RC evidence. Their checksums and source commits remain scoped to those exact builds. [Verified v0.1.0 provenance](docs/releases/v0.1.0.md) · [Pinned RC evidence](docs/RELEASE_EVIDENCE.md)
- [x] Built and inspected all three local preparation ZIPs for Windows amd64, Linux amd64 and macOS arm64. Checked CRC, SHA-256, architecture, the project MIT license, all four upstream license texts, and changelog. The executables match the already-reviewed RC binaries byte for byte. The ZIPs have separate checksums because their contents include new release material.
- [x] Record an express MIT grant for eligible original `v0.1.0` material. Its original archives predate the project MIT notice. [Historical grant](docs/releases/v0.1.0-MIT.md)
- [x] Verify the stable package and full source archive at the accepted commit and tag. Anonymous README, source, CI and asset checks passed. The release procedure covers these checks and the stable version update. [Release procedure](docs/releases/OPEN_SOURCE_RELEASE.md)

## GitHub configuration

- [x] Verified public visibility on the original repository, default branch `main`, and accurate project metadata. [Public repository](https://github.com/soroushneyestani/developer-cli-toolkit)
- [x] Main requires pull requests, all four observed GitHub Actions jobs with current base, and resolved conversations. Force pushes and branch deletion are blocked. Zero separate approving reviews are required, and merge commits are supported. GitHub private vulnerability reporting is enabled. [Required checks and launch settings](docs/releases/OPEN_SOURCE_RELEASE.md#github-settings-for-launch)
- [x] The owner approved visibility on this existing repository subject to the publication audit. Screen third-party data and credentials before making it public. [GitHub visibility effects](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/managing-repository-settings/setting-repository-visibility#making-a-repository-public)

## Documentation

- [x] Prepared source and binary instructions, contribution guidance, private reporting guidance, final learning material, and platform and performance limitations. [README](README.md) · [Contributing](CONTRIBUTING.md) · [Security policy](SECURITY.md) · [Persian learning guide](docs/LEARNING_FA.md) · [Portfolio article](docs/WEBSITE_EN.md) · [Instagram copy](docs/INSTAGRAM_EN.md)
- [x] Verify working public release and source links. The English portfolio and Instagram copy are supplied as repository documentation.

## Remaining risks

The owner approved disclosure and MIT for eligible original material. The public stable assets passed exact-commit CI and checksum verification. Temporary disk can grow with dense output. The sorted paths and the largest emitted record still affect memory and the input tree needs to stay stable during execution. The benchmark samples are Windows-only, and the binaries are unsigned. [Benchmark scope](docs/BENCHMARKS.md) · [Release limits](docs/releases/OPEN_SOURCE_RELEASE.md#limits-to-disclose)

## Final review

- [x] Public stable release and exact build verified. [Release](https://github.com/soroushneyestani/developer-cli-toolkit/releases/tag/v0.1.1) · [Evidence](docs/releases/v0.1.1.md)
- [x] Obtain the owner’s final A02 acceptance and verify the exact post-PR #3 main commit, including macOS and packaging. [Completion record](docs/COMPLETION.md)

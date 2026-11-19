# Open-source release procedure

The owner approved public MIT distribution of eligible original material in the existing repository. The historical `v0.1.0` tag and its assets trace to their original main build. PR #2 carries streaming and the stable `v0.1.1` preparation. [Historical baseline](v0.1.0.md) · [Merged PR #2](https://github.com/soroushneyestani/developer-cli-toolkit/pull/2) · [Publication checklist](../../FINAL_PUBLICATION_CHECKLIST.md)

## Review before a visibility change

Confirm the owner holds the rights to license the original A02 work. Review the historical author metadata, pull requests, older releases and retained Actions records for public disclosure. GitHub exposes the code and Actions history and logs when a repository changes from private to public. Review how to license or retain the historical `v0.1.0` source and binaries, which were built before the project MIT notice was present. Keep their existing tag, assets and checksums traceable to their original build. [Project license](../../LICENSE) · [Third-party notices](../../THIRD_PARTY_NOTICES.md) · [GitHub visibility changes](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/managing-repository-settings/setting-repository-visibility#making-a-repository-public)

## Build the approved release

1. Review the stable version update before merging. Confirm that `dct --version` would report `0.1.1` and that the changelog and known limitations agree. Accept the exact PR revision and preserve its commits through a merge commit.
2. Wait for the native matrix and packaging to pass on the resulting `main` commit. Use the push run for the exact source SHA. Match the final tag to that SHA. [Verify and package workflow](../../.github/workflows/ci.yml)
3. Recheck the three target ZIPs, `SHA256SUMS.txt`, version output, and smoke results. Check each ZIP for the original project `LICENSE`, `THIRD_PARTY_NOTICES.md`, bundled license texts, and changelog. Retain the exact artifact and record its provenance.
4. Check a source archive of that same commit. It should contain the Go source, `go.mod`, `go.sum`, the license, and the build and test instructions. Build and test from the extracted tree.
5. After separate public-visibility approval, verify from an unauthenticated session that visitors can read and clone `main` and download the release tag's source ZIP or tarball. The release should link to the source tag and offer the verified binaries, SHA-256 manifest, release notes and known limitations. [GitHub source archives](https://docs.github.com/en/repositories/working-with-files/using-files/downloading-source-code-archives)
6. Verify the public sources and downloadable assets against the approved ref and SHA-256 manifest. Confirm branch rules and private vulnerability reporting, then publish the portfolio and social copy with working public links.

The `0.1.1-rc.1` review builds have prerelease version output. The historical RC and preparation packages remain useful as evidence for their own revisions. The eventual stable assets need their own exact source and checksums. [Release evidence](../RELEASE_EVIDENCE.md)

## GitHub settings for launch

The default branch is `main`. Keep merge commits enabled. On the public repository, protect `main` by requiring a PR, resolved conversations, current CI and the checks below. Block force pushes and branch deletion. Use merge-commit history. Set the required approval count to zero for the sole maintainer. Keep the PR, CI and conversation-resolution requirements. [GitHub protected branches](https://docs.github.com/en/repositories/configuring-branches-and-merges-in-your-repository/managing-protected-branches/about-protected-branches)

- `Native verification (ubuntu-24.04)`
- `Native verification (windows-2025)`
- `Native verification (macos-15)`
- `Package binaries`

The private repository previously hit a plan restriction. The effective public settings are recorded in the stable release evidence. [Verified configuration](v0.1.1.md#public-repository-and-audit) Enable private vulnerability reporting when GitHub makes it available and check the reporting form. Keep the current descriptive repository title, project homepage, and accurate Go, CLI, search, inspection, and NDJSON topics. [GitHub vulnerability reporting](https://docs.github.com/en/code-security/how-tos/report-and-fix-vulnerabilities/configure-vulnerability-reporting/configure-for-a-repository)

## Limits to disclose

The CLI performs deterministic, literal search on a stable input tree. Paths and the largest emitted record affect memory, and dense matching output can require temporary disk. The Windows benchmark compares specific synthetic fixtures. The oversized matching line uses less memory and takes longer than `v0.1.0` on that fixture. Binaries are unsigned. [CLI contract](../SPEC.md) · [Benchmark evidence](../BENCHMARKS.md)

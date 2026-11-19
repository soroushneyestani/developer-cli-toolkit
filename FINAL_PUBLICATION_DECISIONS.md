# A02 final publication decisions

A02 is in `release_review`. The repository is private, `main` and `v0.1.0` point to `02472793619d8a644634068759ecb2e826a39797`, and PR #2 is Draft. This audit identifies the rights and disclosure decisions that still need the owner's approval. [Draft PR #2](https://github.com/soroushneyestani/developer-cli-toolkit/pull/2) · [Publication checklist](FINAL_PUBLICATION_CHECKLIST.md)

## 1. MIT and ownership

The audited Git history and recorded provenance support original A02 development. The original Go source, scripts, tests, synthetic fixtures and documentation trace to the project owner's commits and a GitHub-generated merge. The review found no identified employer, client or unrelated-project code, vendor source import, or separately identified outside contributor. Git and GitHub record attribution and repository control. Licensing authority still depends on ownership or permission, including any employment, client, contributor or assignment terms that the repository cannot establish. There is no identified competing ownership claim. The unresolved item is the owner's confirmation that they have authority to grant MIT for the original material in the scope they approve. If any file or portion is subject to conflicting terms, identify it before publishing. [WIPO copyright guidance](https://www.wipo.int/en/web/copyright/protection)

MIT is the selected policy for original A02 material, subject to the authority confirmation above. The root `LICENSE` uses the standard MIT text. Preserve its copyright and permission notices in copies and substantial portions of the original work. The dependency notices and terms apply to the respective third-party portions. [Project license](LICENSE) · [SPDX MIT text](https://spdx.org/licenses/MIT.html)

| Linked component | Pinned license and distribution evidence |
| --- | --- |
| Cobra 1.10.2 | Apache-2.0. The full upstream license matches the collected text. The pinned module has no separate upstream `NOTICE` file. |
| mousetrap 1.1.0 | Apache-2.0 and linked on Windows. Its full license matches the collected text and the pinned module has no separate upstream `NOTICE` file. |
| pflag 1.0.9 | BSD-3-Clause. The collected text preserves copyright, conditions, disclaimer and non-endorsement terms. |
| Go runtime and standard library as built with Go 1.27.4 | BSD-3-Clause, including the linked Go-vendored DNS message package with the same license text. |

The pinned module and toolchain texts match the collected texts after line-ending normalization. The audited binary ZIPs carry the four upstream texts and third-party notices. Retain full Apache license texts and applicable source notices, include supplied upstream `NOTICE` material if introduced, and mark relevant modifications to Apache-licensed files if made. Retain BSD copyright, conditions and disclaimer in the appropriate source or binary distribution material. [Bundled notices](THIRD_PARTY_NOTICES.md) · [Apache-2.0 redistribution terms](https://www.apache.org/licenses/LICENSE-2.0#redistribution) · [BSD-3-Clause text](https://spdx.org/licenses/BSD-3-Clause.html)

## 2. Historical Git email privacy

Compared the raw author and committer fields of all 16 commits at the audited preparation checkpoint with authenticated GitHub commit records. The 14 historical commits use one shared address in both fields. The current public GitHub profile does not list that address, and there is no verified designation for public professional contact. Treat it as potentially sensitive. The report omits address values.

| Classification | Distinct addresses | Author fields | Committer fields |
| --- | ---: | ---: | ---: |
| GitHub noreply, owner and GitHub platform identities | 2 | 2 | 2 |
| Public professional address, confirmed | 0 | 0 | 0 |
| Personal or potentially sensitive address | 1 | 14 | 14 |

The affected commit prefixes are `1dbf562`, `170dc70`, `6d78a2e`, `ee3cb74`, `cb952c8`, `482ebac`, `367acfe`, `1d41dcc`, `7927304`, `95fd1cc`, `1bb67e8`, `aaf9a9c`, `541aff8`, and `f6c3b97`. The signed release merge `0247279` and the preparation commit `b43e08c` use noreply in both fields. The release tag is lightweight and points directly to its commit.

| Audited GitHub branch, tag or pull request | History or PR commit count | Commits with the potentially sensitive address |
| --- | ---: | ---: |
| `main` | 10 | 9 |
| `v0.1.0` | 10 | 9 |
| `codex/a02` | 9 | 9 |
| `codex/a02-streaming` at `b43e08c` | 16 | 14 |
| Merged PR #1, commits listed on the PR | 8 | 8 |
| Draft PR #2, commits listed on the PR at `b43e08c` | 6 | 5 |

The server's PR head refs matched the two reviewed head commits. The temporary GitHub merge ref for Draft PR #2 used noreply metadata and also references the existing ancestry. The earlier publication preparation screened 16 retained Actions log archives and the two PR and release records for address and credential patterns. The public-disclosure decision covers the Git objects as well as the existing GitHub surface. A clone exposes the raw fields even when a profile or UI hides the email. [GitHub email settings](https://docs.github.com/en/account-and-profile/how-tos/email-preferences/setting-your-commit-email-address) · [Visibility change effects](https://docs.github.com/en/repositories/managing-your-repositorys-settings-and-features/managing-repository-settings/setting-repository-visibility#making-a-repository-public)

Only the local A02 repository now has its default commit email set to the account's verified GitHub noreply. The active account ID and login matched the address already used in the release merge and preparation commit. Both effective Git identity fields resolve to noreply. This configuration applies to future commits.

| Publication option | Effects on hashes, tags, PRs, releases and CI |
| --- | --- |
| Keep the existing repository private and publish a separate, standalone source snapshot after approval. Recommended when withholding the address. | The private repository continues to hold its original hashes, signed merge, `v0.1.0` tag, PRs, releases and CI. Export the approved source into a fresh repository with verified noreply and a separately reviewed public tree. The public repository gets its own commit hashes, PRs, tags, builds, release records and provenance. Copy source and license material without carrying private Git objects, PR threads, logs or artifacts. |
| Make the existing repository public after informed consent to disclose the address. | Existing hashes, tag, PRs, release and exact-commit CI stay attached to the same history. The raw address in 14 commits becomes accessible through the public history, including nine commits reachable from `v0.1.0`. Existing repository records and Actions history also need to be suitable for public access. |
| Rewrite the existing repository's history to change historical addresses. | Changes affected commit IDs and descendant IDs, including the signed release merge. Preserving a tag at the original SHA retains access to the old metadata. PR diffs and references can break, exact-commit CI and old release provenance refer to the original IDs, and cached objects, forks and clones can retain the old data. This conflicts with the requirement to preserve the original tag and history in place. |

Profile settings and display aliases leave stored commit fields available. An independent source snapshot provides isolation while retaining all original private records. [GitHub's history rewrite caveats](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/removing-sensitive-data-from-a-repository)

## 3. Historical v0.1.0

Live downloads of GitHub's automatic source ZIP and tarball matched all 33 files at `02472793619d8a644634068759ecb2e826a39797`. Both formats lack a project MIT `LICENSE` and contain the four upstream texts. Re-downloaded all five uploaded release assets and matched their recorded IDs, sizes and SHA-256 hashes. The three target ZIPs passed CRC and contain the third-party notice, upstream texts and changelog. They lack a project MIT notice. The published checksum file remains tied to those original ZIP bytes. [v0.1.0 release](https://github.com/soroushneyestani/developer-cli-toolkit/releases/tag/v0.1.0) · [Exact build evidence](docs/releases/v0.1.0.md) · [Original manifest](docs/releases/v0.1.0-SHA256SUMS.txt)

If the person or entity with the relevant rights authorizes it, they can expressly grant MIT for the original A02 material at that exact tagged revision, including their portions of the identified binaries. Use a dated grant identifying the commit, covered originals and original binary payloads by package or digest. Attach or link the exact MIT text and preserve the original upstream terms. This is a new grant for identified existing material. GitHub's automatic tag archives continue to contain the tag's original tree. [GitHub source archives](https://docs.github.com/en/repositories/working-with-files/using-files/downloading-source-code-archives) · [MIT terms](https://spdx.org/licenses/MIT.html)

Keep the original tag and the five original release assets as the historical record. If an approved distribution needs self-contained MIT notices, create distinctly named licensed companion archives. The source companion should contain the exact tagged files plus the new project license and provenance notice. For binary companions, preserve the verified original executable bytes, add the project license and provenance notice, and retain all upstream texts. Compute separate outer archive hashes and a separate manifest. Record the source SHA, original package and executable digests, and repackaging date. Describe the new files as licensed wrappers of the original build. Use a separately approved distribution location if the original release disallows asset additions. [GitHub release asset management](https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository)

With the separate public snapshot option, the existing `v0.1.0` release can remain private. Under the same-repository option the old tag and release become visible with the repository, so explicitly decide how their license status will be presented before switching visibility.

## 4. v0.1.1 readiness

The exact preparation checkpoint `b43e08c36e4ee893164df23fb48ccf11a4a6e886` passed the three native jobs and packaging on both push and Draft PR workflows. The extracted 49-file source archive built and passed the native verifier. The source and three target packages contain the root MIT text and four upstream license texts. Their ZIP CRCs and checksums verified. The Windows executable reports `0.1.1-rc.1`, and all three packaged executables match the reviewed candidate. [Push CI](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44570221378) · [Draft PR CI](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44570237842) · [Preparation artifact](https://github.com/soroushneyestani/developer-cli-toolkit/actions/runs/44570221378/artifacts/14553008712)

After choosing the publication destination and approving the scope, set the [CLI version](main.go) to `0.1.1` and prepare the [stable changelog entry](CHANGELOG.md) and release notes on the chosen path. Review that exact source revision, verify fresh Windows, Linux and macOS CI, and build and check newly checksummed source and binary archives with the project and upstream licenses. If integrating into the existing repository, merge only the approved PR revision. Match the future `v0.1.1` tag and release provenance to the exact verified stable commit in the publication repository. The verified prerelease evidence continues to describe the RC. [Stable release procedure](docs/releases/OPEN_SOURCE_RELEASE.md)

## Owner decisions required

1. Confirm your authority to grant MIT for the original A02 work, including relevant employment, client, contributor and assignment rights. Identify any specific exceptions.
2. Choose the publication destination. Authorize public exposure of the existing Git metadata, or authorize preparation of a separate standalone public repository from a reviewed source snapshot. For the second option, provide the destination owner and repository name. The independent snapshot is recommended if you want the historical address kept private.
3. Decide the scope of the historical release. Authorize a dated, explicit MIT grant for original `v0.1.0` material and separately identified licensed companions, or retain that release as a private archive on the separate-repository path. If publishing the existing repository while withholding a historical grant, approve how to label the visible legacy material and its license status.
4. Approve the exact stable publication work after choosing the route. Specify authorization for the stable version change and, once its revision and fresh evidence are reviewed, any PR merge, new version tag, visibility change or release.

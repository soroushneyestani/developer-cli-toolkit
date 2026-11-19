# Contributing to Developer CLI Toolkit

Use a GitHub issue for bug reports or proposals about the CLI. Show the command, `dct --version`, OS, exit status, and relevant stdout and stderr. Use a small synthetic workspace that reproduces the behavior. Report security concerns through the private process. [Security policy](SECURITY.md)

## Build and test

Create a topic branch from `main` in your fork or authorized checkout. Install Go 1.26 or later and PowerShell 7 on Linux or macOS. CI uses Go 1.27.4. From the project directory, run the verifier. Windows PowerShell can invoke it with `powershell.exe -NoProfile -ExecutionPolicy Bypass -File scripts/verify.ps1`.

```sh
pwsh -NoProfile -File scripts/verify.ps1
```

The verifier checks Go formatting, module checksums, vet, uncached tests, and the native build. Exercise changed command behavior with process-level fixtures that check stdout, stderr, and exit status. Keep sorting, symlink handling, UTF-8 validation, and error semantics covered. Document any intentional contract change and explain its compatibility impact. Cite the reproducible workload when proposing a performance change. [CLI contract](docs/SPEC.md) · [Benchmark method](docs/BENCHMARKS.md)

## Submit a pull request

Open a PR against `main` with the reason for the change, the tests you ran, and any limitations. Use a draft for work that still needs verification. Keep changes small enough to review and preserve the notices of any third-party material you have permission to use. The maintainer reviews changes and merges through a merge commit after CI succeeds and review conversations are resolved. [Verification workflow](.github/workflows/ci.yml)

Submit work you have the right to contribute. Original contributions are made under the project's MIT license. Preserve existing third-party license terms and attribution. [Project license](LICENSE) · [Third-party notices](THIRD_PARTY_NOTICES.md)

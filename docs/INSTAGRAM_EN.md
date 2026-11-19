# Developer CLI Toolkit

A script should be able to distinguish an empty search from a failed run. That requirement shaped `dct`, a Go and Cobra CLI I built for workspace automation.

`dct search` finds literal text. `dct inspect` lists files with their byte sizes. Both support suffix filters, sorted relative paths, and newline-delimited JSON. I designed the exit codes so scripts can handle a match, an empty search, and an error explicitly.

Developer CLI Toolkit v0.1.1 is open source under MIT, with native CI on Windows, Linux and macOS. Search reads files in chunks and waits for full validation before sending matches to stdout. If a NUL byte appears at the end of a file, a match in its first line still gets discarded.

I compared the two versions on five reproducible synthetic Windows workloads and verified their output before timing. On a sparse 96 MiB fixture, median peak process memory fell from 57.24 to 7.57 MiB. The long matching line used less memory and took longer to process.

The engineering lesson was to treat the command line as an interface. Output, error messages, and exit codes all need a documented contract and tests.

Source, downloads and reproducible benchmark data are on GitHub. [Developer CLI Toolkit v0.1.1](https://github.com/soroushneyestani/developer-cli-toolkit/releases/tag/v0.1.1)

#GoLang #DeveloperTools #CLI #SoftwareEngineering #EngineeringMastery

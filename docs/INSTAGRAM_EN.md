# Developer CLI Toolkit

A script should be able to distinguish an empty search from a failed run. That requirement shaped `dct`, a Go and Cobra CLI I built for workspace automation.

`dct search` finds literal text. `dct inspect` lists files with their byte sizes. Both support suffix filters, sorted relative paths, and newline-delimited JSON. I designed the exit codes so scripts can handle a match, an empty search, and an error explicitly.

The release candidate passed native CI on Windows, Linux, and macOS. A separate packaging job built three platform archives and generated SHA-256 checksums.

The engineering lesson was to treat the command line as an interface. Output, error messages, and exit codes all need a documented contract and tests.

#GoLang #DeveloperTools #CLI #SoftwareEngineering #EngineeringMastery

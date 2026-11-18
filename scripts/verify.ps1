[CmdletBinding()]
param()
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest
$projectRoot = Split-Path -Parent $PSScriptRoot
Push-Location $projectRoot
try {
    if (-not (Get-Command go -CommandType Application -ErrorAction SilentlyContinue) -or
            -not (Get-Command gofmt -CommandType Application -ErrorAction SilentlyContinue)) {
        throw 'Go and gofmt are required on PATH. Install the Go toolchain and rerun scripts/verify.ps1.'
    }
    if (-not (Test-Path -LiteralPath 'go.mod')) { throw 'Implementation is not bootstrapped: go.mod is missing.' }
    $packageDirectories = & go list -mod=readonly -buildvcs=false -f '{{.Dir}}' ./...
    if ($LASTEXITCODE -ne 0) { throw 'go list failed.' }
    $sources = @($packageDirectories | ForEach-Object {
        Get-ChildItem -LiteralPath $_ -Filter '*.go' -File
    })
    if ($sources.Count -eq 0) { throw 'No Go source files exist.' }
    if (@($sources | Where-Object { $_.Name -like '*_test.go' }).Count -eq 0) { throw 'Real Go tests are required.' }
    $unformatted = & gofmt -l @($sources.FullName)
    if ($LASTEXITCODE -ne 0) { throw 'gofmt failed.' }
    if ($unformatted) { throw "Unformatted Go files: $($unformatted -join ', ')" }
    & go mod verify
    if ($LASTEXITCODE -ne 0) { throw 'Module integrity check failed.' }
    & go vet -mod=readonly -buildvcs=false ./...
    if ($LASTEXITCODE -ne 0) { throw 'go vet failed.' }
    & go test -mod=readonly -buildvcs=false -count=1 -v ./...
    if ($LASTEXITCODE -ne 0) { throw 'Go tests failed.' }
    $executableSuffix = & go env GOEXE
    if ($LASTEXITCODE -ne 0) { throw 'Could not determine the native executable suffix.' }
    & go build -mod=readonly -trimpath -buildvcs=false -o "bin/dct$executableSuffix" .
    if ($LASTEXITCODE -ne 0) { throw 'Go build failed.' }
    Write-Output 'Formatting, module integrity, vet, tests, and native build passed.'
} finally {
    Pop-Location
}

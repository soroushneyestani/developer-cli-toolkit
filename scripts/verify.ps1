[CmdletBinding()]
param()
$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest
$projectRoot = Split-Path -Parent $PSScriptRoot
Push-Location $projectRoot
try {
    if (-not (Get-Command go -ErrorAction SilentlyContinue)) {
        throw 'Go is required on PATH. Install Go, reopen the shell, and rerun scripts/verify.ps1.'
    }
    if (-not (Test-Path -LiteralPath 'go.mod')) { throw 'Implementation is not bootstrapped: go.mod is missing.' }
    $sources = @(Get-ChildItem -LiteralPath $projectRoot -Recurse -Filter '*.go' -File |
        Where-Object { $_.FullName -notmatch '[\\/](vendor|\.git)[\\/]' })
    if ($sources.Count -eq 0) { throw 'No Go source files exist.' }
    if (@($sources | Where-Object { $_.Name -like '*_test.go' }).Count -eq 0) { throw 'Real Go tests are required.' }
    $unformatted = & gofmt -l @($sources.FullName)
    if ($LASTEXITCODE -ne 0) { throw 'gofmt failed.' }
    if ($unformatted) { throw "Unformatted Go files: $($unformatted -join ', ')" }
    & go vet ./...
    if ($LASTEXITCODE -ne 0) { throw 'go vet failed.' }
    & go test -count=1 ./...
    if ($LASTEXITCODE -ne 0) { throw 'Go tests failed.' }
    & go build ./...
    if ($LASTEXITCODE -ne 0) { throw 'Go build failed.' }
    Write-Output 'Formatting, vet, tests, and native build passed.'
} finally {
    Pop-Location
}

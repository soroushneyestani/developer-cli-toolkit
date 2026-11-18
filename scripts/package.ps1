[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'
Set-StrictMode -Version Latest

if (-not (Get-Command go -CommandType Application -ErrorAction SilentlyContinue)) {
    throw 'Go is required on PATH to build the packages.'
}

$projectRoot = Split-Path -Parent $PSScriptRoot
$dist = Join-Path $projectRoot 'dist'
if (Test-Path -LiteralPath $dist) {
    $distItem = Get-Item -LiteralPath $dist -Force
    if (-not $distItem.PSIsContainer -or ($distItem.Attributes -band [IO.FileAttributes]::ReparsePoint)) {
        throw 'dist must be a local directory.'
    }
}

$targets = @(
    @{ OS = 'linux'; Arch = 'amd64'; Binary = 'dct'; Archive = 'dct-linux-amd64.zip' }
    @{ OS = 'windows'; Arch = 'amd64'; Binary = 'dct.exe'; Archive = 'dct-windows-amd64.zip' }
    @{ OS = 'darwin'; Arch = 'arm64'; Binary = 'dct'; Archive = 'dct-darwin-arm64.zip' }
)
$runId = '{0}-{1}' -f [DateTime]::UtcNow.ToString('yyyyMMddTHHmmssZ'), [Guid]::NewGuid().ToString('N')
$staging = Join-Path $dist ".staging-$runId"
$destination = Join-Path $dist "dct-$runId"
$checksumsName = 'SHA256SUMS.txt'

$requiredNotices = @('THIRD_PARTY_NOTICES.md', 'licenses')
$notices = @($requiredNotices | ForEach-Object {
    $source = Join-Path $projectRoot $_
    if (-not (Test-Path -LiteralPath $source)) { throw "Required package content is missing: $_" }
    $source
})
$notices += @('CHANGELOG.md', 'LICENSE' | ForEach-Object {
    $source = Join-Path $projectRoot $_
    if (Test-Path -LiteralPath $source -PathType Leaf) { $source }
})

$savedEnvironment = @{}
foreach ($name in 'GOOS', 'GOARCH', 'CGO_ENABLED') {
    $savedEnvironment[$name] = [Environment]::GetEnvironmentVariable($name, 'Process')
}

Push-Location $projectRoot
try {
    [void][IO.Directory]::CreateDirectory($staging)
    $env:CGO_ENABLED = '0'
    $checksums = foreach ($target in $targets) {
        $archivePath = Join-Path $staging $target.Archive
        $binaryPath = Join-Path $staging $target.Binary
        $env:GOOS = $target.OS
        $env:GOARCH = $target.Arch

        & go build -mod=readonly -trimpath -buildvcs=false -o $binaryPath .
        if ($LASTEXITCODE -ne 0) { throw "Build failed for $($target.OS)/$($target.Arch)." }
        Compress-Archive -LiteralPath (@($binaryPath) + $notices) -DestinationPath $archivePath
        [IO.File]::Delete($binaryPath)

        $hash = (Get-FileHash -LiteralPath $archivePath -Algorithm SHA256).Hash.ToLowerInvariant()
        "$hash  $($target.Archive)"
    }
    [IO.File]::WriteAllText((Join-Path $staging $checksumsName), (($checksums -join "`n") + "`n"))
    [IO.Directory]::Move($staging, $destination)
    Write-Output "Packages and SHA-256 checksums: $destination"
} finally {
    foreach ($name in $savedEnvironment.Keys) {
        [Environment]::SetEnvironmentVariable($name, $savedEnvironment[$name], 'Process')
    }
    Pop-Location
    if ([IO.Directory]::Exists($staging)) {
        $stagedNames = @($targets | ForEach-Object { $_.Archive }) + @('dct', 'dct.exe', $checksumsName)
        foreach ($name in $stagedNames) {
            [IO.File]::Delete((Join-Path $staging $name))
        }
        [IO.Directory]::Delete($staging)
    }
}

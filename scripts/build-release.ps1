param(
    [string]$Version = "dev"
)

$ErrorActionPreference = "Stop"

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$projectRoot = Split-Path -Parent $scriptDir
$distDir = Join-Path $projectRoot "dist"
$releaseWorkDir = Join-Path $distDir "_release"
$assetsDir = Join-Path $projectRoot "release-assets"

if (Test-Path $releaseWorkDir) {
    Remove-Item -Recurse -Force $releaseWorkDir
}

New-Item -ItemType Directory -Path $distDir -Force | Out-Null
New-Item -ItemType Directory -Path $releaseWorkDir -Force | Out-Null

$targets = @(
    @{ Arch = "amd64"; Label = "windows-amd64" },
    @{ Arch = "arm64"; Label = "windows-arm64" }
)

Push-Location $projectRoot
try {
    foreach ($target in $targets) {
        $folderName = "download-cleaner-$($Version)-$($target.Label)"
        $workDir = Join-Path $releaseWorkDir $folderName

        New-Item -ItemType Directory -Path $workDir -Force | Out-Null

        $env:CGO_ENABLED = "0"
        $env:GOOS = "windows"
        $env:GOARCH = $target.Arch

        $exePath = Join-Path $workDir "download-cleaner.exe"

        Write-Host "Building $($target.Label)..."
        go build -trimpath -ldflags "-s -w" -o $exePath .
        if ($LASTEXITCODE -ne 0) {
            throw "go build failed for $($target.Label)"
        }

        Copy-Item -Path (Join-Path $assetsDir "Run Download Cleaner.bat") -Destination $workDir -Force
        Copy-Item -Path (Join-Path $assetsDir "How to Use Download Cleaner.txt") -Destination $workDir -Force

        $zipPath = Join-Path $distDir "$folderName.zip"
        if (Test-Path $zipPath) {
            Remove-Item -Force $zipPath
        }

        Compress-Archive -Path (Join-Path $workDir "*") -DestinationPath $zipPath -Force
        Write-Host "Created $zipPath"
    }
}
finally {
    Pop-Location
}

Write-Host "Release build complete. ZIP files are in: $distDir"

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$root = Split-Path -Parent $MyInvocation.MyCommand.Path

$modules = Get-ChildItem -Path $root -Recurse -Filter go.mod -File |
  Where-Object { $_.FullName -notmatch "[/\\\\]\\.git([/\\\\]|$)" } |
  ForEach-Object { $_.Directory.FullName } |
  Sort-Object -Unique

$subModules = $modules | Where-Object { $_ -ne $root }

foreach ($module in $subModules) {
  $rel = $module.Substring($root.Length).TrimStart("\", "/")
  Write-Host "tidy: $rel"
  Push-Location $module
  try {
    go mod tidy
  } finally {
    Pop-Location
  }
}

if (Test-Path -Path (Join-Path $root "go.mod") -PathType Leaf) {
  Write-Host "tidy: ./"
  Push-Location $root
  try {
    go mod tidy
  } finally {
    Pop-Location
  }
}

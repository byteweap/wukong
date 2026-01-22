Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

if ($args.Count -lt 1 -or [string]::IsNullOrWhiteSpace($args[0])) {
  Write-Error "usage: .\release.ps1 <version> (e.g. 0.0.1 or v0.0.1)"
  exit 1
}

$version = $args[0].Trim()
if ($version.StartsWith("v")) {
  $version = $version.Substring(1)
}

$root = Split-Path -Parent $MyInvocation.MyCommand.Path

$modules = Get-ChildItem -Path $root -Recurse -Filter go.mod -File |
  Where-Object { $_.FullName -notmatch "[/\\\\]\\.git([/\\\\]|$)" } |
  ForEach-Object { $_.Directory.FullName } |
  Sort-Object -Unique

$tags = New-Object System.Collections.Generic.List[string]

foreach ($module in $modules) {
  if ($module -eq $root) {
    $tags.Add("v$version") | Out-Null
  } else {
    $rel = $module.Substring($root.Length).TrimStart("\", "/")
    $tags.Add("$rel/v$version") | Out-Null
  }
}

foreach ($tag in $tags) {
  git tag $tag
}

Write-Host ("tagged: " + ($tags -join " "))

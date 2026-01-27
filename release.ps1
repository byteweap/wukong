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

$expected = "v$version"
$goMods = Get-ChildItem -Path $root -Recurse -Filter go.mod -File |
  Where-Object { $_.FullName -notmatch "[/\\\\]\\.git([/\\\\]|$)" }

foreach ($gomod in $goMods) {
  $content = Get-Content $gomod.FullName
  $updated = $false
  $inRequireBlock = $false
  for ($i = 0; $i -lt $content.Count; $i++) {
    $line = $content[$i]
    if ($line -match "^\s*require\s*\(\s*$") {
      $inRequireBlock = $true
      continue
    }
    if ($inRequireBlock -and $line -match "^\s*\)\s*$") {
      $inRequireBlock = $false
      continue
    }
    if ($line -match "^\s*require\s+github.com/byteweap/wukong\s+\S+") {
      $content[$i] = $line -replace "github.com/byteweap/wukong\s+\S+", "github.com/byteweap/wukong $expected"
      $updated = $true
      continue
    }
    if ($inRequireBlock -and $line -match "^\s*github.com/byteweap/wukong\s+\S+") {
      $content[$i] = $line -replace "github.com/byteweap/wukong\s+\S+", "github.com/byteweap/wukong $expected"
      $updated = $true
      continue
    }
  }
  if ($updated) {
    Set-Content -Path $gomod.FullName -Value $content
  }
}

$badDeps = New-Object System.Collections.Generic.List[string]
foreach ($gomod in $goMods) {
  $lines = Get-Content $gomod.FullName
  $inRequireBlock = $false
  foreach ($line in $lines) {
    if ($line -match "^\s*require\s*\(\s*$") {
      $inRequireBlock = $true
      continue
    }
    if ($inRequireBlock -and $line -match "^\s*\)\s*$") {
      $inRequireBlock = $false
      continue
    }
    if ($line -match "^\s*require\s+github.com/byteweap/wukong\s+(\S+)") {
      $ver = $Matches[1]
      if ($ver -ne $expected) {
        $badDeps.Add("$($gomod.FullName): $ver") | Out-Null
      }
      continue
    }
    if ($inRequireBlock -and $line -match "^\s*github.com/byteweap/wukong\s+(\S+)") {
      $ver = $Matches[1]
      if ($ver -ne $expected) {
        $badDeps.Add("$($gomod.FullName): $ver") | Out-Null
      }
      continue
    }
  }
}

if ($badDeps.Count -gt 0) {
  Write-Error ("dependency version mismatch, expected " + $expected + ":`n" + ($badDeps -join "`n"))
  exit 1
}

git add -A
$status = git status --porcelain
if ($status) {
  git commit -m "release: v$version"
  $branch = (git rev-parse --abbrev-ref HEAD).Trim()
  git push origin $branch
}

$tags = New-Object System.Collections.Generic.List[string]

foreach ($module in $modules) {
  if ($module -eq $root) {
    $tags.Add("v$version") | Out-Null
  } else {
    $rel = $module.Substring($root.Length).TrimStart("\", "/")
    $rel = $rel -replace "\\", "/"
    $tags.Add("$rel/v$version") | Out-Null
  }
}

foreach ($tag in $tags) {
  git tag $tag
}

git push origin @tags

Write-Host ("tagged: " + ($tags -join " "))

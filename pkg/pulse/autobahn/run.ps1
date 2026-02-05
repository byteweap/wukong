$ErrorActionPreference = 'Stop'

$root = Split-Path -Parent $PSScriptRoot
$reports = Join-Path $PSScriptRoot 'reports'

Write-Host "[1/3] Starting pulse autobahn server..."
$serverProc = Start-Process -FilePath "go" -ArgumentList @("run", ".\autobahn") -WorkingDirectory $root -PassThru -WindowStyle Hidden

for ($i = 0; $i -lt 30; $i++) {
	$ready = Test-NetConnection -ComputerName 127.0.0.1 -Port 9001 -WarningAction SilentlyContinue
	if ($ready.TcpTestSucceeded) { break }
	Start-Sleep -Milliseconds 200
}

Write-Host "[2/3] Running Autobahn test suite..."
if (!(Test-Path $reports)) {
	New-Item -ItemType Directory -Force -Path $reports | Out-Null
}

$dockerArgs = @(
	"run", "--rm",
	"--add-host=host.docker.internal:host-gateway",
	"-v", "${PSScriptRoot}:/config",
	"-v", "${reports}:/reports",
	"crossbario/autobahn-testsuite",
	"wstest", "-m", "fuzzingclient", "-s", "/config/fuzzingserver.json"
)

$docker = Start-Process -FilePath "docker" -ArgumentList $dockerArgs -WorkingDirectory $PSScriptRoot -NoNewWindow -Wait -PassThru

Write-Host "[3/3] Stopping server..."
if ($serverProc -and !$serverProc.HasExited) {
	$serverProc | Stop-Process -Force
}

if ($docker.ExitCode -ne 0) {
	throw "Autobahn tests failed with exit code $($docker.ExitCode)."
}

Write-Host "Done. Report: $reports\servers\index.html"

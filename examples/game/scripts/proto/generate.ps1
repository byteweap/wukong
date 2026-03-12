Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$root = Resolve-Path (Join-Path $PSScriptRoot "../..")
$protoDir = Join-Path $root "api/proto"
$outDir = Join-Path $root "internal/pb"

& docker run --rm `
  -v "${protoDir}:/defs" `
  -v "${outDir}:/out" `
  rvolosatovs/protoc `
  --proto_path=/defs `
  --go_out=/out `
  --go_opt=paths=source_relative `
  /defs/*.proto

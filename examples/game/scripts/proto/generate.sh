#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)"
proto_dir="${root}/proto"
out_dir="${root}/internal/pb"

docker run --rm \
  -v "${proto_dir}:/defs" \
  -v "${out_dir}:/out" \
  rvolosatovs/protoc \
  --proto_path=/defs \
  --go_out=/out \
  --go_opt=paths=source_relative \
  /defs/*.proto

#!/usr/bin/env bash
set -euo pipefail

# generate Go sources for the .proto files using the same Docker image
# that the PowerShell helper uses.  The container mounts the directory
# containing this script as /defs and writes the generated files there.

root="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

docker run --rm -v "${root}:/defs" rvolosatovs/protoc \
    --proto_path=/defs --go_out=/defs /defs/*.proto

# move the generated files out of the temporary "pb" subdirectory and
# clean up the directory the generator creates.
if [[ -d "${root}/pb" ]]; then
    mv -f "${root}/pb"/*.go "${root}/"
    rm -rf "${root}/pb"
fi

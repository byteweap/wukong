#!/bin/bash
set -euo pipefail

version="${1:-}"
if [[ -z "${version}" ]]; then
  echo "usage: ./release.sh <version> (e.g. 0.0.1 or v0.0.1)" >&2
  exit 1
fi

version="${version#v}"

mapfile -t modules < <(
  find "${PWD}" -name go.mod -type f \
    -not -path "*/.git/*" \
    -print |
  sed 's|/go.mod$||' |
  sort
)

tags=()
for module in "${modules[@]}"; do
  if [[ "${module}" == "${PWD}" ]]; then
    tags+=("v${version}")
    continue
  fi
  rel="${module#${PWD}/}"
  tags+=("${rel}/v${version}")
done

for t in "${tags[@]}"; do
  git tag "$t"
done

echo "tagged: ${tags[*]}"

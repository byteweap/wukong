#!/bin/bash
set -euo pipefail

readonly directory=$(cd "$(dirname "$0")" && pwd)

mapfile -t modules < <(
  find "${directory}" -name go.mod -type f \
    -not -path "*/.git/*" \
    -print |
  sed 's|/go.mod$||' |
  sort
)

root_module="${directory}"
sub_modules=()
for module in "${modules[@]}"; do
  if [[ "${module}" == "${root_module}" ]]; then
    continue
  fi
  sub_modules+=("${module}")
done

for module in "${sub_modules[@]}"; do
  echo "tidy: ${module#${directory}/}"
  cd "${module}"
  go mod tidy
done

if [[ -f "${root_module}/go.mod" ]]; then
  echo "tidy: ./"
  cd "${root_module}"
  go mod tidy
fi

#!/usr/bin/env bash
set -euo pipefail

if [[ $# -lt 1 || -z "${1// }" ]]; then
  echo "usage: ./release.sh <version> (e.g. 0.0.1 or v0.0.1)" >&2
  exit 1
fi

version="${1#v}"
expected="v${version}"

root="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

mapfile -t go_mods < <(find "$root" -name go.mod -print0 | xargs -0 -n1 printf '%s\n' | grep -vE '/\.git(/|$)')

update_go_mod() {
  local file="$1"
  local tmp
  tmp="$(mktemp)"
  awk -v expected="$expected" '
    BEGIN { in_require=0 }
    /^[[:space:]]*require[[:space:]]*\([[:space:]]*$/ { in_require=1; print; next }
    in_require && /^[[:space:]]*\)[[:space:]]*$/ { in_require=0; print; next }
    /^[[:space:]]*require[[:space:]]+github.com\/byteweap\/wukong[[:space:]]+\S+/ {
      sub(/github.com\/byteweap\/wukong[[:space:]]+\S+/, "github.com/byteweap/wukong " expected)
      print
      next
    }
    in_require && /^[[:space:]]*github.com\/byteweap\/wukong[[:space:]]+\S+/ {
      sub(/github.com\/byteweap\/wukong[[:space:]]+\S+/, "github.com/byteweap/wukong " expected)
      print
      next
    }
    { print }
  ' "$file" > "$tmp"
  if ! cmp -s "$file" "$tmp"; then
    mv "$tmp" "$file"
  else
    rm -f "$tmp"
  fi
}

for gomod in "${go_mods[@]}"; do
  update_go_mod "$gomod"
done

bad_deps=()
for gomod in "${go_mods[@]}"; do
  in_require=0
  while IFS= read -r line; do
    if [[ "$line" =~ ^[[:space:]]*require[[:space:]]*\([[:space:]]*$ ]]; then
      in_require=1
      continue
    fi
    if [[ $in_require -eq 1 && "$line" =~ ^[[:space:]]*\)[[:space:]]*$ ]]; then
      in_require=0
      continue
    fi
    if [[ "$line" =~ ^[[:space:]]*require[[:space:]]+github.com/byteweap/wukong[[:space:]]+([^[:space:]]+) ]]; then
      ver="${BASH_REMATCH[1]}"
      if [[ "$ver" != "$expected" ]]; then
        bad_deps+=("${gomod}: ${ver}")
      fi
      continue
    fi
    if [[ $in_require -eq 1 && "$line" =~ ^[[:space:]]*github.com/byteweap/wukong[[:space:]]+([^[:space:]]+) ]]; then
      ver="${BASH_REMATCH[1]}"
      if [[ "$ver" != "$expected" ]]; then
        bad_deps+=("${gomod}: ${ver}")
      fi
      continue
    fi
  done < "$gomod"
done

if [[ ${#bad_deps[@]} -gt 0 ]]; then
  echo "dependency version mismatch, expected ${expected}:" >&2
  printf '%s\n' "${bad_deps[@]}" >&2
  exit 1
fi

git add -A
if [[ -n "$(git status --porcelain)" ]]; then
  git commit -m "release: v${version}"
  branch="$(git rev-parse --abbrev-ref HEAD)"
  git push origin "$branch"
fi

mapfile -t modules < <(printf '%s\n' "${go_mods[@]}" | xargs -n1 dirname | sort -u)

tags=()
for module in "${modules[@]}"; do
  if [[ "$module" == "$root" ]]; then
    tags+=("v${version}")
  else
    rel="${module#${root}/}"
    rel="${rel//\\//}"
    tags+=("${rel}/v${version}")
  fi
done

for tag in "${tags[@]}"; do
  git tag "$tag"
done

git push origin --tags

echo "tagged: ${tags[*]}"

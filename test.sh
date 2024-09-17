#!/bin/bash

set -euo pipefail

# Build the go-bcov binary
go build -o bin/go-bcov *.go

# Define an array of repositories and their branches
declare -A repos=(
  ["https://github.com/jellydator/ttlcache"]="v3.3.0"
  ["https://github.com/rs/zerolog"]="v1.33.0"
  ["https://github.com/expr-lang/expr"]="v1.16.9"
)

process_repo() {
  local repo="$1"
  local branch="$2"
  local action="$3"

  tmp_dir="$(mktemp -d)"
  trap 'rm -rf -- "${tmp_dir:-?}"' EXIT INT TERM HUP

  git clone --depth 1 --branch "$branch" "$repo" "${tmp_dir}"

  cp bin/go-bcov "${tmp_dir}"
  cp "testdata/$(basename "$repo")_coverage.out" "${tmp_dir}/coverage.out" 2>/dev/null || true

  pushd "${tmp_dir}"

  [[ ! -f coverage.out ]] &&
    go test -coverprofile coverage.out -covermode count ./...

  ./go-bcov -format sonar-cover-report <coverage.out >coverage.xml

  popd

  if [[ $action == "check" ]]; then
    go run tools/xmldiff.go "testdata/$(basename "$repo")_sonar-cover-report.xml" "${tmp_dir}/coverage.xml"
  elif [[ $action == "generate" ]]; then
    cp "${tmp_dir}/coverage.xml" "testdata/$(basename "$repo")_sonar-cover-report.xml"
    cp "${tmp_dir}/coverage.out" "testdata/$(basename "$repo")_coverage.out"
  fi
}

check() {
  for repo in "${!repos[@]}"; do
    process_repo "$repo" "${repos[$repo]}" "check"
  done
}

generate_test_data() {
  for repo in "${!repos[@]}"; do
    process_repo "$repo" "${repos[$repo]}" "generate"
  done
}

if [[ ${1:-} == "--generate-test-data" ]]; then
  generate_test_data
else
  check
fi

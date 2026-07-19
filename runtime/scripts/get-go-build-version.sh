#!/usr/bin/env bash

set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

{
  printf '%s\0' go.mod go.sum
  find cmd internal \
    -type f \
    -name '*.go' \
    ! -name '*_test.go' \
    -print0
} | LC_ALL=C sort -z | xargs -0 sha256sum | sha256sum | \
  cut -c1-12 | sed 's/^/dev-go-/'

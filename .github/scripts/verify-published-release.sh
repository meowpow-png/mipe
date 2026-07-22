#!/usr/bin/env bash

set -euo pipefail

if [[ ! -f "$RELEASE_SOURCE_MANIFEST" ]]; then
  echo "Release source manifest does not exist: $RELEASE_SOURCE_MANIFEST" >&2
  exit 1
fi

targets=(
  runtime
  codex
  codex-java
  codex-web
  claude
  claude-java
  claude-web
)

for target in "${targets[@]}"; do
  repository="$(jq -er --arg target "$target" '.images[$target].repository | select(type == "string" and length > 0)' "$RELEASE_SOURCE_MANIFEST")"
  expected_digest="$(jq -er --arg target "$target" '.images[$target].digest | select(type == "string" and test("^sha256:[0-9a-f]{64}$"))' "$RELEASE_SOURCE_MANIFEST")"

  for image_tag in "v${RELEASE_VERSION}" latest; do
    actual_digest="$(
      docker buildx imagetools inspect "${repository}:${image_tag}" --format '{{json .Manifest}}' \
        | jq -er '.digest | select(type == "string" and test("^sha256:[0-9a-f]{64}$"))'
    )"
    if [[ "$actual_digest" != "$expected_digest" ]]; then
      echo "${repository}:${image_tag} resolves to $actual_digest, not $expected_digest." >&2
      exit 1
    fi
  done
done

gh release view "$RELEASE_TAG" --repo "$GITHUB_REPOSITORY" >/dev/null

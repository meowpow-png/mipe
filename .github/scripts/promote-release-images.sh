#!/usr/bin/env bash

set -euo pipefail

release_image_tag="v${RELEASE_VERSION}"

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

repositories=()
digests=()

for target in "${targets[@]}"; do
  repositories+=("$(jq -er --arg target "$target" '.images[$target].repository | select(type == "string" and length > 0)' "$RELEASE_SOURCE_MANIFEST")")
  digests+=("$(jq -er --arg target "$target" '.images[$target].digest | select(type == "string" and test("^sha256:[0-9a-f]{64}$"))' "$RELEASE_SOURCE_MANIFEST")")
done

version_tags_exist=()

for index in "${!repositories[@]}"; do
  repository="${repositories[$index]}"
  digest="${digests[$index]}"
  source_reference="${repository}@${digest}"
  version_reference="${repository}:${release_image_tag}"

  inspect_output="$(mktemp)"
  if docker buildx imagetools inspect "$version_reference" --format '{{json .Manifest}}' >"$inspect_output" 2>&1; then
    existing_digest="$(jq -er '.digest | select(type == "string" and test("^sha256:[0-9a-f]{64}$"))' "$inspect_output")"
    rm -f "$inspect_output"

    if [[ "$existing_digest" != "$digest" ]]; then
      echo "Release tag $version_reference already points to $existing_digest, not $digest." >&2
      exit 1
    fi

    echo "$version_reference already points to $digest."
    version_tags_exist[$index]=true
  else
    if ! grep -Eqi 'manifest unknown|name unknown|not found' "$inspect_output"; then
      cat "$inspect_output" >&2
      rm -f "$inspect_output"
      exit 1
    fi
    rm -f "$inspect_output"
    version_tags_exist[$index]=false
  fi
done

for index in "${!repositories[@]}"; do
  repository="${repositories[$index]}"
  digest="${digests[$index]}"
  source_reference="${repository}@${digest}"
  version_reference="${repository}:${release_image_tag}"

  if [[ "${version_tags_exist[$index]}" == false ]]; then
    docker buildx imagetools create \
      --tag "$version_reference" \
      "$source_reference"
  fi

  docker buildx imagetools create \
    --tag "${repository}:latest" \
    "$source_reference"
done

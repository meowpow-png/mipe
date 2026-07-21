#!/usr/bin/env bash

set -euo pipefail

artifact_dir="$(mktemp -d)"
trap 'rm -rf "$artifact_dir"' EXIT

gh run download "$CI_RUN_ID" \
  --repo "$GITHUB_REPOSITORY" \
  --name runtime-release-source \
  --dir "$artifact_dir"

manifest_path="$(find "$artifact_dir" -type f -name release-source.json -print -quit)"
if [[ -z "$manifest_path" ]]; then
  echo "CI run $CI_RUN_ID did not provide a runtime release-source manifest." >&2
  exit 1
fi

jq -e \
  --arg source_sha "$SOURCE_SHA" \
  --arg ci_run_id "$CI_RUN_ID" \
  '
    .schema_version == 1
    and .source_sha == $source_sha
    and (.ci_run_id | tostring) == $ci_run_id
  ' \
  "$manifest_path" >/dev/null

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
  output_name="${target//-/_}"
  repository="$(jq -er --arg target "$target" '.images[$target].repository | select(type == "string" and length > 0)' "$manifest_path")"
  digest="$(jq -er --arg target "$target" '.images[$target].digest | select(type == "string" and test("^sha256:[0-9a-f]{64}$"))' "$manifest_path")"

  {
    echo "${output_name}_repository=$repository"
    echo "${output_name}_digest=$digest"
  } >> "$GITHUB_OUTPUT"
done

build_version="$(jq -er '.build_version | select(type == "string" and length > 0)' "$manifest_path")"
echo "build_version=$build_version" >> "$GITHUB_OUTPUT"

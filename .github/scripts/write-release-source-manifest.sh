#!/usr/bin/env bash

set -euo pipefail

if [[ -n "$IMAGE_TAGS" ]]; then
  image_tag="${IMAGE_TAGS%%,*}"
  published_tags="$IMAGE_TAGS"
else
  image_tag="$BUILD_VERSION"
  published_tags="dev-latest,$BUILD_VERSION"
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

repositories=(
  "$IMAGE_PREFIX"
  "$IMAGE_PREFIX-codex"
  "$IMAGE_PREFIX-codex-java"
  "$IMAGE_PREFIX-codex-web"
  "$IMAGE_PREFIX-claude"
  "$IMAGE_PREFIX-claude-java"
  "$IMAGE_PREFIX-claude-web"
)

digests=(
  "${RUNTIME_DIGEST:-}"
  "${CODEX_DIGEST:-}"
  "${CODEX_JAVA_DIGEST:-}"
  "${CODEX_WEB_DIGEST:-}"
  "${CLAUDE_DIGEST:-}"
  "${CLAUDE_JAVA_DIGEST:-}"
  "${CLAUDE_WEB_DIGEST:-}"
)

for index in "${!targets[@]}"; do
  digest="${digests[$index]}"
  if [[ ! "$digest" =~ ^sha256:[0-9a-f]{64}$ ]]; then
    echo "Missing or invalid Buildx digest for ${targets[$index]} (${repositories[$index]}:${image_tag})." >&2
    exit 1
  fi
done

jq -n \
  --arg source_sha "$SOURCE_SHA" \
  --arg ci_run_id "$CI_RUN_ID" \
  --arg build_version "$BUILD_VERSION" \
  --arg image_tags "$published_tags" \
  --arg runtime_repository "${repositories[0]}" \
  --arg runtime_digest "${digests[0]}" \
  --arg codex_repository "${repositories[1]}" \
  --arg codex_digest "${digests[1]}" \
  --arg codex_java_repository "${repositories[2]}" \
  --arg codex_java_digest "${digests[2]}" \
  --arg codex_web_repository "${repositories[3]}" \
  --arg codex_web_digest "${digests[3]}" \
  --arg claude_repository "${repositories[4]}" \
  --arg claude_digest "${digests[4]}" \
  --arg claude_java_repository "${repositories[5]}" \
  --arg claude_java_digest "${digests[5]}" \
  --arg claude_web_repository "${repositories[6]}" \
  --arg claude_web_digest "${digests[6]}" \
  '
    {
      schema_version: 1,
      source_sha: $source_sha,
      ci_run_id: $ci_run_id,
      build_version: $build_version,
      image_tags: $image_tags,
      images: {
        runtime: {repository: $runtime_repository, digest: $runtime_digest},
        codex: {repository: $codex_repository, digest: $codex_digest},
        "codex-java": {repository: $codex_java_repository, digest: $codex_java_digest},
        "codex-web": {repository: $codex_web_repository, digest: $codex_web_digest},
        claude: {repository: $claude_repository, digest: $claude_digest},
        "claude-java": {repository: $claude_java_repository, digest: $claude_java_digest},
        "claude-web": {repository: $claude_web_repository, digest: $claude_web_digest}
      }
    }
  '

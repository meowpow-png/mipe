#!/usr/bin/env bash

set -euo pipefail

if [[ -n "$IMAGE_TAGS" ]]; then
  image_tag="${IMAGE_TAGS%%,*}"
  published_tags="$IMAGE_TAGS"
else
  image_tag="$BUILD_VERSION"
  published_tags="dev-latest,$BUILD_VERSION"
fi

repositories=(
  "$IMAGE_PREFIX"
  "$IMAGE_PREFIX-codex"
  "$IMAGE_PREFIX-codex-java"
  "$IMAGE_PREFIX-codex-web"
  "$IMAGE_PREFIX-claude"
  "$IMAGE_PREFIX-claude-java"
  "$IMAGE_PREFIX-claude-web"
)

if [[ -z "${BUILD_METADATA:-}" ]]; then
  echo "Buildx metadata is empty." >&2
  exit 1
fi

jq -e \
  --arg source_sha "$SOURCE_SHA" \
  --arg ci_run_id "$CI_RUN_ID" \
  --arg build_version "$BUILD_VERSION" \
  --arg image_tags "$published_tags" \
  --arg runtime_repository "${repositories[0]}" \
  --arg codex_repository "${repositories[1]}" \
  --arg codex_java_repository "${repositories[2]}" \
  --arg codex_web_repository "${repositories[3]}" \
  --arg claude_repository "${repositories[4]}" \
  --arg claude_java_repository "${repositories[5]}" \
  --arg claude_web_repository "${repositories[6]}" \
  '
    def image($target; $repository):
      .[$target]["containerimage.digest"]
      | select(type == "string" and test("^sha256:[0-9a-f]{64}$"))
      | {repository: $repository, digest: .};

    {
      schema_version: 1,
      source_sha: $source_sha,
      ci_run_id: $ci_run_id,
      build_version: $build_version,
      image_tags: $image_tags,
      images: {
        runtime: image("runtime"; $runtime_repository),
        codex: image("codex"; $codex_repository),
        "codex-java": image("codex-java"; $codex_java_repository),
        "codex-web": image("codex-web"; $codex_web_repository),
        claude: image("claude"; $claude_repository),
        "claude-java": image("claude-java"; $claude_java_repository),
        "claude-web": image("claude-web"; $claude_web_repository)
      }
    }
  ' <<<"$BUILD_METADATA"

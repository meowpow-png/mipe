#!/usr/bin/env bash

set -euo pipefail

jq -e \
  --arg source_sha "$SOURCE_SHA" \
  --arg ci_run_id "$CI_RUN_ID" \
  --arg build_version "$BUILD_VERSION" \
  --arg image_prefix "$IMAGE_PREFIX" \
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
      images: {
        runtime: image("runtime"; $image_prefix),
        codex: image("codex"; $image_prefix + "-codex"),
        "codex-java": image("codex-java"; $image_prefix + "-codex-java"),
        "codex-web": image("codex-web"; $image_prefix + "-codex-web"),
        claude: image("claude"; $image_prefix + "-claude"),
        "claude-java": image("claude-java"; $image_prefix + "-claude-java"),
        "claude-web": image("claude-web"; $image_prefix + "-claude-web")
      }
    }
  ' <<<"$BUILD_METADATA"

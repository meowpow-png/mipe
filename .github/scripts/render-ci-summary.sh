#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(git rev-parse --show-toplevel)"
RUNTIME_DIR="${ROOT_DIR}/runtime"
TEMPLATE="${ROOT_DIR}/.github/templates/ci-summary.md"

eval "$(
  docker buildx bake \
    --file "${RUNTIME_DIR}/docker-bake.hcl" \
    --print \
  | jq -r '
      .target
      | map_values(.args)
      | values
      | add
      | with_entries(select(.key | IN(
          "NODE_VERSION",
          "GIT_VERSION",
          "CODEX_VERSION",
          "CLAUDE_VERSION",
          "TEMURIN_21_JDK_VERSION",
          "CHROMIUM_VERSION",
          "PLAYWRIGHT_MCP_VERSION"
      )))
      | to_entries[]
      | "export \(.key)=\(.value | @sh)"
    '
)"

if [[ -z "${BUILD_METADATA:-}" ]]; then
  printf 'Missing required build metadata output.\n' >&2
  exit 1
fi

extract_digest() {
  local target="$1"
  local variable="$2"
  local digest

  if ! digest="$(
    jq -er \
      --arg target "$target" \
      '.[$target]["containerimage.digest"] // empty
       | select(type == "string" and test("^sha256:[0-9a-f]{64}$"))' \
      <<<"${BUILD_METADATA}"
  )"; then
    printf 'Missing or invalid image digest for Bake target %s.\n' "$target" >&2
    exit 1
  fi

  printf -v "$variable" '%s' "$digest"
  export "$variable"
}

targets=(
  runtime
  codex
  codex-java
  codex-web
  claude
  claude-java
  claude-web
)

variables=(
  RUNTIME_DIGEST
  CODEX_DIGEST
  CODEX_JAVA_DIGEST
  CODEX_WEB_DIGEST
  CLAUDE_DIGEST
  CLAUDE_JAVA_DIGEST
  CLAUDE_WEB_DIGEST
)

for index in "${!targets[@]}"; do
  extract_digest "${targets[$index]}" "${variables[$index]}"
done

required=(
  BRANCH
  COMMIT
  VERSION

  UNIT_TESTS
  INTEGRATION_TESTS
  COVERAGE

  NODE_VERSION
  GIT_VERSION
  CODEX_VERSION
  CLAUDE_VERSION
  TEMURIN_21_JDK_VERSION
  CHROMIUM_VERSION
  PLAYWRIGHT_MCP_VERSION

  RUNTIME_DIGEST
  CODEX_DIGEST
  CODEX_JAVA_DIGEST
  CODEX_WEB_DIGEST
  CLAUDE_DIGEST
  CLAUDE_JAVA_DIGEST
  CLAUDE_WEB_DIGEST

  OWNER
)

missing=()

for var in "${required[@]}"; do
  if [[ -z "${!var:-}" ]]; then
    missing+=("$var")
  fi
done

if ((${#missing[@]})); then
  printf 'Missing required environment variables:\n' >&2
  printf '  %s\n' "${missing[@]}" >&2
  exit 1
fi

envsubst < "$TEMPLATE"

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

status_marker() {
  case "${1,,}" in
    success|pass) printf '✅' ;;
    failure|fail) printf '❌' ;;
    skipped) printf '⏭️' ;;
    cancelled) printf '🚫' ;;
    *) printf '%s' "$1" ;;
  esac
}

UNIT_TESTS_DISPLAY="$(status_marker "$UNIT_TESTS")"
INTEGRATION_TESTS_DISPLAY="$(status_marker "$INTEGRATION_TESTS")"
export UNIT_TESTS_DISPLAY INTEGRATION_TESTS_DISPLAY

envsubst < "$TEMPLATE"

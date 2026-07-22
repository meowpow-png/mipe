#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(git rev-parse --show-toplevel)"
TEMPLATE="${ROOT_DIR}/.github/templates/coverage-summary.md"

required=(
  MODULE
  BRANCH
  COMMIT
  COVERAGE
  CODECOV_STATUS
)

missing=()

for var in "${required[@]}"; do
  if [[ -z "${!var:-}" ]]; then
    missing+=("$var")
  fi
done

if ((${#missing[@]})); then
  printf 'Missing required coverage-summary environment variables:\n' >&2
  printf '  %s\n' "${missing[@]}" >&2
  exit 1
fi

case "${CODECOV_STATUS,,}" in
  success) CODECOV_STATUS_DISPLAY='✅' ;;
  failure) CODECOV_STATUS_DISPLAY='❌' ;;
  *) CODECOV_STATUS_DISPLAY="${CODECOV_STATUS}" ;;
esac

export CODECOV_STATUS_DISPLAY
envsubst < "${TEMPLATE}"

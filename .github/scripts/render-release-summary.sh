#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(git rev-parse --show-toplevel)"
TEMPLATE="${ROOT_DIR}/.github/templates/release-summary.md"

required=(
  RELEASE_TAG
  RELEASE_VERSION
  SOURCE_SHA
  SOURCE_CANDIDATE_TAG
  SOURCE_RC_RUN_ID
  REPOSITORY

  VALIDATE_STATUS
  RESOLVE_STATUS
  PROMOTE_STATUS
  PUBLISH_STATUS
  VERIFY_STATUS
)

missing=()

for variable in "${required[@]}"; do
  if [[ -z "${!variable:-}" ]]; then
    missing+=("$variable")
  fi
done

if ((${#missing[@]})); then
  printf 'Missing required release-summary variables:\n' >&2
  printf '  %s\n' "${missing[@]}" >&2
  exit 1
fi

envsubst < "$TEMPLATE"

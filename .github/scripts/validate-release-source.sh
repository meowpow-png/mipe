#!/usr/bin/env bash

set -euo pipefail

if [[ ! "$RELEASE_TAG" =~ ^runtime/v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]; then
  echo "Release tag must match runtime/vMAJOR.MINOR.PATCH: $RELEASE_TAG" >&2
  exit 1
fi

release_version="${RELEASE_TAG#runtime/v}"
source_sha="$(git rev-parse "${RELEASE_TAG}^{commit}")"

git fetch --no-tags origin +refs/heads/dev:refs/remotes/origin/dev
if ! git merge-base --is-ancestor "$source_sha" origin/dev; then
  echo "Release tag $RELEASE_TAG must point to a commit reachable from dev." >&2
  exit 1
fi

changelog_pattern="^## \[${release_version//./\\.}\] - [0-9]{4}-[0-9]{2}-[0-9]{2}$"
changelog_entries="$(grep -Ec "$changelog_pattern" runtime/CHANGELOG.md || true)"
if [[ "$changelog_entries" -ne 1 ]]; then
  echo "runtime/CHANGELOG.md must contain exactly one dated entry for $release_version." >&2
  exit 1
fi

ci_run="$({
  gh api --paginate \
    "/repos/${GITHUB_REPOSITORY}/actions/workflows/ci.yml/runs?head_sha=${source_sha}&status=completed&per_page=100" \
    | jq -s '
        [.[].workflow_runs[]
         | select(.event == "push" or .event == "workflow_dispatch")]
        | sort_by(.run_started_at)
        | last // empty
      '
})"

if [[ -z "$ci_run" ]]; then
  echo "No completed CI run exists for $source_sha. Run CI for this commit, then retry the release." >&2
  exit 1
fi

ci_run_id="$(jq -r '.id' <<<"$ci_run")"
ci_conclusion="$(jq -r '.conclusion' <<<"$ci_run")"
if [[ "$ci_conclusion" != "success" ]]; then
  echo "Latest CI run $ci_run_id for $source_sha concluded with $ci_conclusion. Run CI successfully, then retry the release." >&2
  exit 1
fi

{
  echo "release_tag=$RELEASE_TAG"
  echo "release_version=$release_version"
  echo "source_sha=$source_sha"
  echo "ci_run_id=$ci_run_id"
} >> "$GITHUB_OUTPUT"

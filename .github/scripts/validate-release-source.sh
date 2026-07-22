#!/usr/bin/env bash

set -euo pipefail

if [[ ! "$RELEASE_TAG" =~ ^runtime/v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]; then
  echo "Release tag must match runtime/vMAJOR.MINOR.PATCH: $RELEASE_TAG" >&2
  exit 1
fi

release_version="${RELEASE_TAG#runtime/v}"
release_branch="release/runtime-v${release_version}"
changelog_version="runtime-${release_version}"
source_sha="$(git rev-parse "${RELEASE_TAG}^{commit}")"

if ! git fetch --no-tags origin "+refs/heads/${release_branch}:refs/remotes/origin/${release_branch}"; then
  echo "Release branch $release_branch does not exist." >&2
  exit 1
fi

if ! git merge-base --is-ancestor "$source_sha" "origin/${release_branch}"; then
  echo "Release tag $RELEASE_TAG must point to a commit reachable from $release_branch." >&2
  exit 1
fi

changelog_pattern="^## \[${changelog_version//./\\.}\] - [0-9]{4}-[0-9]{2}-[0-9]{2}$"
changelog_entries="$(grep -Ec "$changelog_pattern" CHANGELOG.md || true)"
if [[ "$changelog_entries" -ne 1 ]]; then
  echo "CHANGELOG.md must contain exactly one dated entry for $changelog_version." >&2
  exit 1
fi

candidate_run="$({
  gh api --paginate \
    "/repos/${GITHUB_REPOSITORY}/actions/workflows/runtime-rc.yml/runs?head_sha=${source_sha}&status=completed&per_page=100" \
    | jq -s '
        [.[].workflow_runs[]
         | select(.event == "push" and .conclusion == "success")]
        | sort_by(.run_started_at)
        | last // empty
      '
})"

if [[ -z "$candidate_run" ]]; then
  echo "No successful runtime release-candidate run exists for $source_sha. Publish an RC for this commit, then retry the release." >&2
  exit 1
fi

candidate_run_id="$(jq -r '.id' <<<"$candidate_run")"
candidate_tag="$(jq -r '.head_branch' <<<"$candidate_run")"

{
  echo "release_tag=$RELEASE_TAG"
  echo "release_version=$release_version"
  echo "source_sha=$source_sha"
  echo "candidate_run_id=$candidate_run_id"
  echo "candidate_tag=$candidate_tag"
} >> "$GITHUB_OUTPUT"

#!/usr/bin/env bash

set -euo pipefail

if [[ ! "$RUNTIME_CANDIDATE_TAG" =~ ^runtime/v((0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)-rc\.(0|[1-9][0-9]*))$ ]]; then
  echo "Runtime candidate tag must match runtime/vMAJOR.MINOR.PATCH-rc.N: $RUNTIME_CANDIDATE_TAG" >&2
  exit 1
fi

candidate_version="${BASH_REMATCH[1]}"
image_tag="v${candidate_version}"
release_version="${candidate_version%-rc.*}"
release_branch="release/runtime-v${release_version}"
changelog_version="runtime-${release_version}"
source_sha="$(git rev-parse "${RUNTIME_CANDIDATE_TAG}^{commit}")"

if ! git fetch --no-tags origin "+refs/heads/${release_branch}:refs/remotes/origin/${release_branch}"; then
  echo "Release branch $release_branch does not exist." >&2
  exit 1
fi

if ! git merge-base --is-ancestor "$source_sha" "origin/${release_branch}"; then
  echo "Runtime candidate tag $RUNTIME_CANDIDATE_TAG must point to a commit reachable from $release_branch." >&2
  exit 1
fi

changelog_pattern="^## \[${changelog_version//./\\.}\] - [0-9]{4}-[0-9]{2}-[0-9]{2}$"
changelog_entries="$(grep -Ec "$changelog_pattern" CHANGELOG.md || true)"
if [[ "$changelog_entries" -ne 1 ]]; then
  echo "CHANGELOG.md must contain exactly one dated entry for $changelog_version." >&2
  exit 1
fi

{
  echo "candidate_version=$candidate_version"
  echo "image_tag=$image_tag"
  echo "release_version=$release_version"
  echo "source_sha=$source_sha"
} >> "$GITHUB_OUTPUT"

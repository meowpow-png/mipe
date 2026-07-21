#!/usr/bin/env bash

set -euo pipefail

changelog_file="${CHANGELOG_FILE:-CHANGELOG.md}"
changelog_heading="## [runtime-${RELEASE_VERSION}] - "
release_notes="$(mktemp)"
trap 'rm -f "$release_notes"' EXIT

awk -v heading="$changelog_heading" '
  index($0, heading) == 1 {
    found = 1
  }
  found && /^## / && index($0, heading) != 1 {
    exit
  }
  found {
    print
  }
  END {
    if (!found) {
      exit 1
    }
  }
' "$changelog_file" > "$release_notes"

release_title="Runtime v${RELEASE_VERSION}"

if gh release view "$RELEASE_TAG" --repo "$GITHUB_REPOSITORY" >/dev/null 2>&1; then
  gh release edit "$RELEASE_TAG" \
    --repo "$GITHUB_REPOSITORY" \
    --title "$release_title" \
    --notes-file "$release_notes"
else
  gh release create "$RELEASE_TAG" \
    --repo "$GITHUB_REPOSITORY" \
    --title "$release_title" \
    --notes-file "$release_notes"
fi

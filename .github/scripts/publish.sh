#!/usr/bin/env bash
set -euo pipefail

owner="${GHCR_OWNER,,}"
version="${VERSION}"

publish_image() {
    local image="$1"
    local suffix=""

    if [[ -n "${image}" ]]; then
        suffix="-${image}"
    fi
    local local_image="mipe-runtime${suffix}:local"
    local remote_image="ghcr.io/${owner}/mipe-runtime${suffix}"

    docker tag "${local_image}" "${remote_image}:dev-latest"
    docker tag "${local_image}" "${remote_image}:${version}"

    docker image push --all-tags "${remote_image}"
}

for image in "" codex codex-java codex-web claude claude-java claude-web; do
    publish_image "${image}"
done

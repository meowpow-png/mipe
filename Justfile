set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

mod runtime "runtime/Justfile"

import 'just/docker.just'

# List available recipes
default:
    @just --list

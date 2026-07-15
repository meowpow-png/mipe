set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

mod pwdgen "examples/password-generator/Justfile"
mod runtime "runtime/Justfile"

import 'just/docker.just'

# List available recipes
default:
    @just --list

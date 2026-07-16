group "default" {
  targets = ["codex", "claude"]
}

target "runtime-base" {
  context = "."
  dockerfile = "docker/runtime/base/Dockerfile"
  tags = ["mipe-runtime-base:latest"]
}

target "codex" {
  context = "."
  dockerfile = "docker/runtime/codex/Dockerfile"
  tags = ["mipe-runtime-codex:latest"]

  contexts = {
    runtime = "target:runtime-base"
  }
}

target "claude" {
  context = "."
  dockerfile = "docker/runtime/claude/Dockerfile"
  tags = ["mipe-runtime-claude:latest"]

  contexts = {
    runtime = "target:runtime-base"
  }
}

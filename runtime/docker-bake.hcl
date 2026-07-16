group "default" {
  targets = ["codex", "claude"]
}

target "runtime-base" {
  context = "."
  dockerfile = "docker/base/Dockerfile"
  tags = ["mipe-runtime-base:latest"]
}

target "codex" {
  context = "."
  dockerfile = "docker/codex/Dockerfile"
  tags = ["mipe-runtime-codex:latest"]

  contexts = {
    runtime = "target:runtime-base"
  }
}

target "claude" {
  context = "."
  dockerfile = "docker/claude/Dockerfile"
  tags = ["mipe-runtime-claude:latest"]

  contexts = {
    runtime = "target:runtime-base"
  }
}

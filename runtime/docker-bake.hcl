group "default" {
  targets = ["codex", "claude"]
}

target "runtime" {
  context = "."
  dockerfile = "docker/base/Dockerfile"
  tags = ["mipe-runtime:latest"]
}

target "codex" {
  context = "."
  dockerfile = "docker/codex/Dockerfile"
  tags = ["mipe-runtime-codex:latest"]

  contexts = {
    runtime = "target:runtime"
  }
}

target "claude" {
  context = "."
  dockerfile = "docker/claude/Dockerfile"
  tags = ["mipe-runtime-claude:latest"]

  contexts = {
    runtime = "target:runtime"
  }
}

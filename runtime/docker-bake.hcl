group "default" {
  targets = [
    "test",
    "codex",
    "claude",
    "codex-java",
    "claude-java"
  ]
}

target "runtime-base" {
  context = "."
  dockerfile = "docker/runtime/base/Dockerfile"
  tags = ["mipe-runtime-base:latest"]
}

target "test" {
  context = "."
  dockerfile = "docker/runtime/test/Dockerfile"
  tags = ["mipe-runtime-test:latest"]

  args = {
    LOCAL_UID = "1000"
    LOCAL_GID = "1000"
  }
  contexts = {
    runtime = "target:runtime-base"
  }
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

target "codex-java" {
  context = "."
  dockerfile = "docker/toolchain/java/codex/Dockerfile"
  tags = ["mipe-runtime-codex-java:latest"]

  contexts = {
    runtime = "target:codex"
  }
}

target "claude-java" {
  context = "."
  dockerfile = "docker/toolchain/java/claude/Dockerfile"
  tags = ["mipe-runtime-claude-java:latest"]

  contexts = {
    runtime = "target:claude"
  }
}

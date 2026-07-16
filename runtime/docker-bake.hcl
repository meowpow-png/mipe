group "default" {
  targets = [
    "test",
    "codex",
    "claude",
    "codex-java",
    "claude-java"
  ]
}

variable "CODEX_VERSION" {
  default = "0.144.5"
}

variable "CLAUDE_VERSION" {
  default = "2.1.211"
}

variable "NODE_VERSION" {
  default = "22.23.1"
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

  args = {
    CODEX_VERSION = CODEX_VERSION
  }

  contexts = {
    runtime = "target:node-base"
  }
}

target "claude" {
  context = "."
  dockerfile = "docker/runtime/claude/Dockerfile"
  tags = ["mipe-runtime-claude:latest"]

  args = {
    CLAUDE_VERSION = CLAUDE_VERSION
  }

  contexts = {
    runtime = "target:node-base"
  }
}

target "node-base" {
  context = "."
  dockerfile = "docker/runtime/node/Dockerfile"

  args = {
    NODE_VERSION = NODE_VERSION
  }

  contexts = {
    runtime = "target:runtime-base"
  }
}

target "java-base" {
  context = "."
  dockerfile = "docker/toolchain/java/base/Dockerfile"

  contexts = {
    runtime = "target:node-base"
  }
}

target "codex-java" {
  context = "."
  dockerfile = "docker/toolchain/java/codex/Dockerfile"
  tags = ["mipe-runtime-codex-java:latest"]

  args = {
    CODEX_VERSION = CODEX_VERSION
  }

  contexts = {
    runtime = "target:java-base"
  }
}

target "claude-java" {
  context = "."
  dockerfile = "docker/toolchain/java/claude/Dockerfile"
  tags = ["mipe-runtime-claude-java:latest"]

  args = {
    CLAUDE_VERSION = CLAUDE_VERSION
  }

  contexts = {
    runtime = "target:java-base"
  }
}

group "default" {
  targets = [
    "test",
    "codex",
    "claude",
    "codex-java",
    "claude-java",
    "codex-web",
    "claude-web"
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

variable "PLAYWRIGHT_MCP_VERSION" {
  default = "0.0.78"
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
  dockerfile = "docker/toolchain/java/Dockerfile"

  contexts = {
    runtime = "target:node-base"
  }
}

target "web-base" {
  context = "."
  dockerfile = "docker/toolchain/web/Dockerfile"

  args = {
    PLAYWRIGHT_MCP_VERSION = PLAYWRIGHT_MCP_VERSION
  }

  contexts = {
    runtime = "target:node-base"
  }
}

target "codex" {
  context = "."
  dockerfile = "docker/runtime/codex/Dockerfile"
  target = "default"
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
  target = "default"
  tags = ["mipe-runtime-claude:latest"]

  args = {
    CLAUDE_VERSION = CLAUDE_VERSION
  }

  contexts = {
    runtime = "target:node-base"
  }
}

target "codex-java" {
  context = "."
  dockerfile = "docker/runtime/codex/Dockerfile"
  target = "default"
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
  dockerfile = "docker/runtime/claude/Dockerfile"
  target = "default"
  tags = ["mipe-runtime-claude-java:latest"]

  args = {
    CLAUDE_VERSION = CLAUDE_VERSION
  }

  contexts = {
    runtime = "target:java-base"
  }
}

target "codex-web" {
  context = "."
  dockerfile = "docker/runtime/codex/Dockerfile"
  target = "web"
  tags = ["mipe-runtime-codex-web:latest"]

  args = {
    CODEX_VERSION = CODEX_VERSION
  }

  contexts = {
    runtime = "target:web-base"
  }
}

target "claude-web" {
  context = "."
  dockerfile = "docker/runtime/claude/Dockerfile"
  target = "web"
  tags = ["mipe-runtime-claude-web:latest"]

  args = {
    CLAUDE_VERSION = CLAUDE_VERSION
  }

  contexts = {
    runtime = "target:web-base"
  }
}

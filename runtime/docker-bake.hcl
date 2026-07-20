group "default" {
  targets = [
    "runtime",
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
  default = "0.144.6"
}

variable "CLAUDE_VERSION" {
  default = "2.1.216"
}

variable "NODE_VERSION" {
  default = "22.23.1"
}

variable "PLAYWRIGHT_MCP_VERSION" {
  default = "0.0.78"
}

variable "VERSION" {
  default = "dev"
}

variable "SOURCE_DATE_EPOCH" {
  default = "1784448885"
}

target "runtime" {
  context = "."
  dockerfile = "docker/Dockerfile"
  target = "runtime"
  tags = ["mipe-runtime:local"]

  args = {
    VERSION            = VERSION
    SOURCE_DATE_EPOCH  = SOURCE_DATE_EPOCH
  }
}

target "test" {
  context = "."
  dockerfile = "docker/Dockerfile"
  target = "test"
  tags = ["mipe-runtime-test:local"]

  args = {
    VERSION           = VERSION
    SOURCE_DATE_EPOCH = SOURCE_DATE_EPOCH
    LOCAL_UID         = "1000"
    LOCAL_GID         = "1000"
  }
}

target "node-base" {
  context = "."
  dockerfile = "docker/Dockerfile"
  target = "node-base"

  args = {
    NODE_VERSION = NODE_VERSION
  }
}

target "java-base" {
  context = "."
  dockerfile = "docker/Dockerfile"
  target = "java-base"

  args = {
    VERSION           = VERSION
    SOURCE_DATE_EPOCH = SOURCE_DATE_EPOCH
    NODE_VERSION      = NODE_VERSION
  }
}

target "web-base" {
  context = "."
  dockerfile = "docker/Dockerfile"
  target = "web-base"

  args = {
    VERSION                = VERSION
    SOURCE_DATE_EPOCH      = SOURCE_DATE_EPOCH
    NODE_VERSION           = NODE_VERSION
    PLAYWRIGHT_MCP_VERSION = PLAYWRIGHT_MCP_VERSION
  }
}

target "codex" {
  context = "."
  dockerfile = "docker/Dockerfile"
  target = "codex"
  tags = ["mipe-runtime-codex:local"]

  args = {
    VERSION                = VERSION
    SOURCE_DATE_EPOCH      = SOURCE_DATE_EPOCH
    NODE_VERSION           = NODE_VERSION
    CODEX_VERSION          = CODEX_VERSION
  }
}

target "claude" {
  context = "."
  dockerfile = "docker/Dockerfile"
  target = "claude"
  tags = ["mipe-runtime-claude:local"]

  args = {
    VERSION                = VERSION
    SOURCE_DATE_EPOCH      = SOURCE_DATE_EPOCH
    NODE_VERSION           = NODE_VERSION
    CLAUDE_VERSION         = CLAUDE_VERSION
  }
}

target "codex-java" {
  context = "."
  dockerfile = "docker/Dockerfile"
  target = "codex-java"
  tags = ["mipe-runtime-codex-java:local"]

  args = {
    VERSION           = VERSION
    SOURCE_DATE_EPOCH = SOURCE_DATE_EPOCH
    NODE_VERSION      = NODE_VERSION
    CODEX_VERSION     = CODEX_VERSION
  }
}

target "claude-java" {
  context = "."
  dockerfile = "docker/Dockerfile"
  target = "claude-java"
  tags = ["mipe-runtime-claude-java:local"]

  args = {
    VERSION           = VERSION
    SOURCE_DATE_EPOCH = SOURCE_DATE_EPOCH
    NODE_VERSION      = NODE_VERSION
    CLAUDE_VERSION = CLAUDE_VERSION
  }
}

target "codex-web" {
  context = "."
  dockerfile = "docker/Dockerfile"
  target = "codex-web"
  tags = ["mipe-runtime-codex-web:local"]

  args = {
    VERSION           = VERSION
    SOURCE_DATE_EPOCH = SOURCE_DATE_EPOCH
    NODE_VERSION      = NODE_VERSION
    CODEX_VERSION     = CODEX_VERSION
    PLAYWRIGHT_MCP_VERSION = PLAYWRIGHT_MCP_VERSION
  }
}

target "claude-web" {
  context = "."
  dockerfile = "docker/Dockerfile"
  target = "claude-web"
  tags = ["mipe-runtime-claude-web:local"]

  args = {
    VERSION           = VERSION
    SOURCE_DATE_EPOCH = SOURCE_DATE_EPOCH
    NODE_VERSION      = NODE_VERSION
    CLAUDE_VERSION = CLAUDE_VERSION
    PLAYWRIGHT_MCP_VERSION = PLAYWRIGHT_MCP_VERSION
  }
}

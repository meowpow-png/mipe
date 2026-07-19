variable "IMAGE_PREFIX" {
  default = ""
}

target "runtime" {
  tags = [
    "${IMAGE_PREFIX}:dev-latest",
    "${IMAGE_PREFIX}:${VERSION}"
  ]
}

target "codex" {
  tags = [
    "${IMAGE_PREFIX}-codex:dev-latest",
    "${IMAGE_PREFIX}-codex:${VERSION}"
  ]
}

target "codex-java" {
  tags = [
    "${IMAGE_PREFIX}-codex-java:dev-latest",
    "${IMAGE_PREFIX}-codex-java:${VERSION}"
  ]
}

target "codex-web" {
  tags = [
    "${IMAGE_PREFIX}-codex-web:dev-latest",
    "${IMAGE_PREFIX}-codex-web:${VERSION}"
  ]
}

target "claude" {
  tags = [
    "${IMAGE_PREFIX}-claude:dev-latest",
    "${IMAGE_PREFIX}-claude:${VERSION}"
  ]
}

target "claude-java" {
  tags = [
    "${IMAGE_PREFIX}-claude-java:dev-latest",
    "${IMAGE_PREFIX}-claude-java:${VERSION}"
  ]
}

target "claude-web" {
  tags = [
    "${IMAGE_PREFIX}-claude-web:dev-latest",
    "${IMAGE_PREFIX}-claude-web:${VERSION}"
  ]
}

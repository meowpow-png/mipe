variable "IMAGE_PREFIX" {
  default = ""
}

variable "IMAGE_TAGS" {
  default = ""
}

target "runtime" {
  tags = IMAGE_TAGS == "" ? [
    "${IMAGE_PREFIX}:dev-latest",
    "${IMAGE_PREFIX}:${VERSION}"
  ] : [for tag in split(",", IMAGE_TAGS) : "${IMAGE_PREFIX}:${tag}"]
}

target "codex" {
  tags = IMAGE_TAGS == "" ? [
    "${IMAGE_PREFIX}-codex:dev-latest",
    "${IMAGE_PREFIX}-codex:${VERSION}"
  ] : [for tag in split(",", IMAGE_TAGS) : "${IMAGE_PREFIX}-codex:${tag}"]
}

target "codex-java" {
  tags = IMAGE_TAGS == "" ? [
    "${IMAGE_PREFIX}-codex-java:dev-latest",
    "${IMAGE_PREFIX}-codex-java:${VERSION}"
  ] : [for tag in split(",", IMAGE_TAGS) : "${IMAGE_PREFIX}-codex-java:${tag}"]
}

target "codex-web" {
  tags = IMAGE_TAGS == "" ? [
    "${IMAGE_PREFIX}-codex-web:dev-latest",
    "${IMAGE_PREFIX}-codex-web:${VERSION}"
  ] : [for tag in split(",", IMAGE_TAGS) : "${IMAGE_PREFIX}-codex-web:${tag}"]
}

target "claude" {
  tags = IMAGE_TAGS == "" ? [
    "${IMAGE_PREFIX}-claude:dev-latest",
    "${IMAGE_PREFIX}-claude:${VERSION}"
  ] : [for tag in split(",", IMAGE_TAGS) : "${IMAGE_PREFIX}-claude:${tag}"]
}

target "claude-java" {
  tags = IMAGE_TAGS == "" ? [
    "${IMAGE_PREFIX}-claude-java:dev-latest",
    "${IMAGE_PREFIX}-claude-java:${VERSION}"
  ] : [for tag in split(",", IMAGE_TAGS) : "${IMAGE_PREFIX}-claude-java:${tag}"]
}

target "claude-web" {
  tags = IMAGE_TAGS == "" ? [
    "${IMAGE_PREFIX}-claude-web:dev-latest",
    "${IMAGE_PREFIX}-claude-web:${VERSION}"
  ] : [for tag in split(",", IMAGE_TAGS) : "${IMAGE_PREFIX}-claude-web:${tag}"]
}

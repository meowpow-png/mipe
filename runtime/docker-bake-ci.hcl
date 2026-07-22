variable "IMAGE_PREFIX" {
  default = ""
}

variable "IMAGE_TAGS" {
  default = ""
}

variable "CACHE_MODE" {
  default = "gha"
}

function "cache_from" {
  params = [scope]
  result = CACHE_MODE == "gha" ? ["type=gha,scope=${scope}"] : []
}

function "cache_to" {
  params = [scope]
  result = CACHE_MODE == "gha" ? ["type=gha,scope=${scope},mode=max"] : []
}

target "runtime" {
  tags = IMAGE_TAGS == "" ? [
    "${IMAGE_PREFIX}:dev-latest",
    "${IMAGE_PREFIX}:${VERSION}"
  ] : [for tag in split(",", IMAGE_TAGS) : "${IMAGE_PREFIX}:${tag}"]
  cache-from = cache_from("runtime")
  cache-to = cache_to("runtime")
}

target "test" {
  cache-from = cache_from("runtime")
}

target "codex" {
  tags = IMAGE_TAGS == "" ? [
    "${IMAGE_PREFIX}-codex:dev-latest",
    "${IMAGE_PREFIX}-codex:${VERSION}"
  ] : [for tag in split(",", IMAGE_TAGS) : "${IMAGE_PREFIX}-codex:${tag}"]
  cache-from = concat(cache_from("node-base"), cache_from("runtime"))
}

target "codex-java" {
  tags = IMAGE_TAGS == "" ? [
    "${IMAGE_PREFIX}-codex-java:dev-latest",
    "${IMAGE_PREFIX}-codex-java:${VERSION}"
  ] : [for tag in split(",", IMAGE_TAGS) : "${IMAGE_PREFIX}-codex-java:${tag}"]
  cache-from = concat(cache_from("java-base"), cache_from("node-base"), cache_from("runtime"))
}

target "codex-web" {
  tags = IMAGE_TAGS == "" ? [
    "${IMAGE_PREFIX}-codex-web:dev-latest",
    "${IMAGE_PREFIX}-codex-web:${VERSION}"
  ] : [for tag in split(",", IMAGE_TAGS) : "${IMAGE_PREFIX}-codex-web:${tag}"]
  cache-from = concat(cache_from("web-base"), cache_from("node-base"), cache_from("runtime"))
}

target "claude" {
  tags = IMAGE_TAGS == "" ? [
    "${IMAGE_PREFIX}-claude:dev-latest",
    "${IMAGE_PREFIX}-claude:${VERSION}"
  ] : [for tag in split(",", IMAGE_TAGS) : "${IMAGE_PREFIX}-claude:${tag}"]
  cache-from = concat(cache_from("node-base"), cache_from("runtime"))
}

target "claude-java" {
  tags = IMAGE_TAGS == "" ? [
    "${IMAGE_PREFIX}-claude-java:dev-latest",
    "${IMAGE_PREFIX}-claude-java:${VERSION}"
  ] : [for tag in split(",", IMAGE_TAGS) : "${IMAGE_PREFIX}-claude-java:${tag}"]
  cache-from = concat(cache_from("java-base"), cache_from("node-base"), cache_from("runtime"))
}

target "claude-web" {
  tags = IMAGE_TAGS == "" ? [
    "${IMAGE_PREFIX}-claude-web:dev-latest",
    "${IMAGE_PREFIX}-claude-web:${VERSION}"
  ] : [for tag in split(",", IMAGE_TAGS) : "${IMAGE_PREFIX}-claude-web:${tag}"]
  cache-from = concat(cache_from("web-base"), cache_from("node-base"), cache_from("runtime"))
}

target "node-base" {
  cache-from = concat(cache_from("node-base"), cache_from("runtime"))
  cache-to = cache_to("node-base")
}

target "java-base" {
  cache-from = concat(cache_from("java-base"), cache_from("node-base"), cache_from("runtime"))
  cache-to = cache_to("java-base")
}

target "web-base" {
  cache-from = concat(cache_from("web-base"), cache_from("node-base"), cache_from("runtime"))
  cache-to = cache_to("web-base")
}

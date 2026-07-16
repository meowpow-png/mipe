group "default" {
  targets = ["codex"]
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

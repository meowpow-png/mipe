# Create and push a runtime release candidate
[arg("version", help="Release version without the v prefix")]
[arg("candidate", help="Release candidate number")]
release-rc version candidate:
    @git tag -a runtime/v{{version}}-rc.{{candidate}} -m "Runtime v{{version}} RC {{candidate}}"
    @git push --atomic origin dev refs/tags/runtime/v{{version}}-rc.{{candidate}}

# Create and push a stable runtime release
[arg("version", help="Release version without the v prefix")]
[arg("candidate", help="Release candidate number")]
release version candidate:
    @git tag -a runtime/v{{version}} -m "Runtime v{{version}}" runtime/v{{version}}-rc.{{candidate}}^{commit}
    @git push origin refs/tags/runtime/v{{version}}

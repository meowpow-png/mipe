# Runtime Build Summary

## Build

| Property        | Value        |
|-----------------|--------------|
| Branch          | `${BRANCH}`  |
| Commit          | `${COMMIT}`  |
| Runtime Version | `${VERSION}` |

## Verification

| Check             | Result               |
|-------------------|----------------------|
| Unit Tests        | ${UNIT_TESTS}        |
| Integration Tests | ${INTEGRATION_TESTS} |
| Coverage          | `${COVERAGE}`        |

## Published Images

| Image         | Image Tags          | Digest                  |
|---------------|---------------------|-------------------------|
| `runtime`     | `${PUBLISHED_TAGS}` | `${RUNTIME_DIGEST}`     |
| `codex`       | `${PUBLISHED_TAGS}` | `${CODEX_DIGEST}`       |
| `codex-java`  | `${PUBLISHED_TAGS}` | `${CODEX_JAVA_DIGEST}`  |
| `codex-web`   | `${PUBLISHED_TAGS}` | `${CODEX_WEB_DIGEST}`   |
| `claude`      | `${PUBLISHED_TAGS}` | `${CLAUDE_DIGEST}`      |
| `claude-java` | `${PUBLISHED_TAGS}` | `${CLAUDE_JAVA_DIGEST}` |
| `claude-web`  | `${PUBLISHED_TAGS}` | `${CLAUDE_WEB_DIGEST}`  |

## Toolchain

| Component      | Version                     |
|----------------|-----------------------------|
| Node.js        | `${NODE_VERSION}`           |
| Git            | `${GIT_VERSION}`            |
| Codex CLI      | `${CODEX_VERSION}`          |
| Claude Code    | `${CLAUDE_VERSION}`         |
| Temurin JDK    | `${TEMURIN_21_JDK_VERSION}` |
| Chromium       | `${CHROMIUM_VERSION}`       |
| Playwright MCP | `${PLAYWRIGHT_MCP_VERSION}` |

## OCI Image References

```text
ghcr.io/${OWNER}/mipe-runtime@${RUNTIME_DIGEST}
ghcr.io/${OWNER}/mipe-runtime-codex@${CODEX_DIGEST}
ghcr.io/${OWNER}/mipe-runtime-codex-java@${CODEX_JAVA_DIGEST}
ghcr.io/${OWNER}/mipe-runtime-codex-web@${CODEX_WEB_DIGEST}
ghcr.io/${OWNER}/mipe-runtime-claude@${CLAUDE_DIGEST}
ghcr.io/${OWNER}/mipe-runtime-claude-java@${CLAUDE_JAVA_DIGEST}
ghcr.io/${OWNER}/mipe-runtime-claude-web@${CLAUDE_WEB_DIGEST}
```

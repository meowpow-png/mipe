# Runtime Release Summary

## Release

| Property      | Value                                                                                    |
|---------------|------------------------------------------------------------------------------------------|
| Tag           | `${RELEASE_TAG}`                                                                         |
| Version       | `${RELEASE_VERSION}`                                                                     |
| Source commit | `${SOURCE_SHA}`                                                                          |
| Source CI run | [${SOURCE_CI_RUN_ID}](https://github.com/${REPOSITORY}/actions/runs/${SOURCE_CI_RUN_ID}) |

## Release Stages

| Stage                             | Status             |
|-----------------------------------|--------------------|
| Validate release tag and source   | ${VALIDATE_STATUS} |
| Resolve development image digests | ${RESOLVE_STATUS}  |
| Promote release images            | ${PROMOTE_STATUS}  |
| Generate release attestations     | ${ATTEST_STATUS}   |
| Create release manifest           | ${MANIFEST_STATUS} |
| Publish GitHub release            | ${PUBLISH_STATUS}  |
| Verify published release          | ${VERIFY_STATUS}   |

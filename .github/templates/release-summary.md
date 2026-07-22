# Runtime Release Summary

## Release

| Property      | Value                                                                                    |
|---------------|------------------------------------------------------------------------------------------|
| Tag           | `${RELEASE_TAG}`                                                                         |
| Version       | `${RELEASE_VERSION}`                                                                     |
| Source commit | `${SOURCE_SHA}`                                                                          |
| Source RC tag | `${SOURCE_CANDIDATE_TAG}`                                                                |
| Source RC run | [${SOURCE_RC_RUN_ID}](https://github.com/${REPOSITORY}/actions/runs/${SOURCE_RC_RUN_ID}) |

## Release Stages

| Stage                             | Status                     |
|-----------------------------------|----------------------------|
| Validate release tag and source   | ${VALIDATE_STATUS_DISPLAY} |
| Resolve release-candidate digests | ${RESOLVE_STATUS_DISPLAY}  |
| Promote release images            | ${PROMOTE_STATUS_DISPLAY}  |
| Publish GitHub release            | ${PUBLISH_STATUS_DISPLAY}  |
| Verify published release          | ${VERIFY_STATUS_DISPLAY}   |

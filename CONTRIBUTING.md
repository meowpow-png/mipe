
## Welcome

Thanks for considering a contribution.

Pull requests are welcome. For larger changes, consider opening an issue first so the direction can be discussed before implementation.

When appropriate, include relevant tests and documentation updates.

## Planning

The roadmap tracks the evolution of Mipe through platform milestones rather than individual component releases. Each milestone should deliver a complete, usable capability that builds upon the previous one and advances the platform toward its long-term vision.

Architectural questions are expected, but they should be answered through the implementation of working software rather than isolated research. Every milestone should leave the project in a state that developers can meaningfully adopt and build upon.

## Components

Mipe consists of independently releasable components:

- `mipe-runtime`
- `mipe-api`
- `mipe-web`

Each component is versioned independently using Semantic Versioning. Components should only be released when they change, regardless of the current roadmap milestone.

## Versioning

All components follow Semantic Versioning.

- Increment patch version for backwards-compatible bug fixes
- Increment minor version for backwards-compatible functionality
- Increment major version for breaking changes

Major versions represent a platform compatibility generation. Components within the same major version are expected to be compatible with one another.

## Releases

Components are released independently as they evolve. A release includes only components that have changed and is assigned a Semantic Version. Container image tags and digests identify the published artifacts for that release but are not component versions.

A roadmap milestone is considered complete once all required component releases are available. Components unaffected by a milestone do not receive a new release.

### Runtime

Runtime releases are prepared on `dev`. Start by committing the completed root changelog entry for the version being released:

```text
## [runtime-X.Y.Z] - YYYY-MM-DD
```

Create an annotated release candidate tag on that commit, then push the branch and tag together:

```bash
git tag -a runtime/vX.Y.Z-rc.1 -m "Runtime vX.Y.Z RC 1"
git push --atomic origin dev refs/tags/runtime/vX.Y.Z-rc.1
```

Wait for the release candidate to pass. It builds and tests the images that will become the stable release. Then create the stable tag on the exact candidate commit and push it:

```bash
git tag -a runtime/vX.Y.Z -m "Runtime vX.Y.Z" runtime/vX.Y.Z-rc.1^{commit}
git push origin refs/tags/runtime/vX.Y.Z
```

If the candidate needs a fix, make a new commit on `dev` and create a new candidate tag such as `runtime/vX.Y.Z-rc.2`. Do not move or reuse existing tags.

See [Continuous Integration](docs/CI.md#releases) for the release pipeline and published image behavior.

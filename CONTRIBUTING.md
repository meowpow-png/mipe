## Releasing

### Planning

The roadmap tracks the evolution of Mipe through platform milestones rather than individual component releases. Each milestone should deliver a complete, usable capability that builds upon the previous one and advances the platform toward its long-term vision.

Architectural questions are expected, but they should be answered through the implementation of working software rather than isolated research. Every milestone should leave the project in a state that developers can meaningfully adopt and build upon.

### Components

Mipe consists of independently releasable components:

- `mipe-runtime`
- `mipe-api`
- `mipe-web`

Each component is versioned independently using Semantic Versioning. Components should only be released when they change, regardless of the current roadmap milestone.

### Versioning

All components follow Semantic Versioning.

- Increment patch version for backwards-compatible bug fixes
- Increment minor version for backwards-compatible functionality
- Increment major version for breaking changes

Major versions represent a platform compatibility generation. Components within the same major version are expected to be compatible with one another.

### Strategy

Components are released independently as they evolve. A release should only include components that have changed. Each release is assigned a Semantic Version.

A roadmap milestone is considered complete once all required component releases are available. Components that are unaffected by a milestone do not receive a new release.

## Publishing

### Building

Every successful CI build on `dev` branch publishes development images to GitHub Container Registry. Failed builds are never published.

Published builds provide immediate access to latest development state. They are intended for testing and evaluation, are not considered releases, and do not receive Semantic Versions.

### Tagging

Each published build is assigned two container image tags:

- `dev-latest` — points to the latest successful build from the `dev` branch
- `dev-<hash>` — identifies the build produced from a specific set of Go build inputs

The `dev-latest` tag provides a convenient moving reference for tracking ongoing development, while `dev-<hash>` provides a stable reference for reproducible environments. The hash is computed from `go.mod`, `go.sum`, and production Go source files.

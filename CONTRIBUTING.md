## Releasing

### Planning

The roadmap tracks the evolution of Mipe through platform milestones rather than component releases. Each milestone should deliver a complete, usable capability that builds upon the previous one and advances the platform toward its long-term vision.

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

Components are released independently as they evolve. A release should only include components that have changed.

A roadmap milestone is considered complete once all required component releases are available. Components that are unaffected by a milestone do not receive a new release.

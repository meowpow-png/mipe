# Note: gosu Usage and Runtime Limitations

## Context

While making Mipe runnable directly on a host, I reviewed where the runtime depends on `gosu` and whether that dependency is necessary for every execution path.

## Observations

Mipe uses `gosu` in three bootstrap paths: workspace validation, dependency initialization, and final command execution. The requested UID and GID come from the resolved runtime configuration.

Mipe changes ownership separately with Go's `os.Chown`. `gosu` is not involved in the recursive ownership update.

When Mipe already runs with the configured effective UID and GID, invoking `gosu` is redundant. The runtime can execute the command directly and therefore does not require `gosu` to be installed in this case.

## Analysis

The runtime now compares the configured UID and GID with the current effective process credentials. It bypasses `gosu` only when both values match. If either value differs, the existing `gosu UID:GID` path remains in use.

Matching process credentials do not guarantee that `UserHome` is correctly owned. Mipe checks the ownership of the complete home directory tree separately and only calls `os.Chown` when an entry has a different UID or GID. An unprivileged host process cannot repair incorrect ownership itself.

The approach remains Unix-specific because both credential handling and file ownership use Unix identity semantics. It does not provide a Windows execution path.

## Conclusions

`gosu` is not required when Mipe is launched as the target user and group. It remains required when Mipe must transition to a different UID or GID.

The `gosu` dependency is therefore conditional rather than universal. The direct-execution path removes the dependency for local host usage without reimplementing `gosu`'s privilege-dropping behavior.

Ownership preparation and process identity are separate concerns. Skipping `gosu` does not remove the need to verify or, when privileged, update ownership.

## Next Steps

- Keep `gosu` available in runtime images for configurations where the target UID or GID differs from the process credentials
- When running Mipe locally, configure the target UID and GID to match the host user and ensure `UserHome` is already correctly owned

No further changes are required unless Mipe needs to support credential switching without an external binary.

## References

- [Bootstrap Configuration and Workspace Resolution](002-note-bootstrap-configuration-and-workspace.md)

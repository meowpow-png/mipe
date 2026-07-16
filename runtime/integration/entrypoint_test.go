package integration_test

import "testing"

func TestEntrypointRequiresDeveloperIdentity(t *testing.T) {
	for _, missing := range []string{"LOCAL_UID", "LOCAL_GID"} {
		t.Run(missing, func(t *testing.T) {
			environment := defaultEnvironment()
			spec := containerSpec{
				environment:        environment,
				replaceEnvironment: true,
				command:            []string{"true"},
			}
			delete(environment, missing)
			result := runContainer(t, spec)
			result.requireFailure(t)
			result.requireOutput(t, missing+" is required")
		})
	}
}

func TestEntrypointRejectsConflictingExistingIdentity(t *testing.T) {
	spec := containerSpec{
		command: []string{"usermod", "--uid", "1001", "dev"},
	}
	name := createContainer(t, spec)

	first := startContainer(t, name)
	first.requireSuccess(t)
	second := startContainer(t, name)
	second.requireFailure(t)
	second.requireOutput(t, "expected 1000")
}

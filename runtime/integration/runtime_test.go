package integration_test

import "testing"

func TestRuntimeHappyPath(t *testing.T) {
	result := runContainer(t, containerSpec{
		command: mipeCommand(defaultConfigPath(), `
			test "$(id -u)" = 1000
			test "$(id -g)" = 1000
			test "$PWD" = /workspace
			test "$HOME" = /home/dev
			test "$AGENT_HOME" = /home/dev/.mipe-agent
			test "$RUNTIME_HOME" = /opt/mipe
			test -z "${USER_HOME:-}"
			test -f "$AGENT_HOME/AGENTS.md"
			test -f /workspace/.test
			echo runtime-happy-path
		`),
	})
	result.requireSuccess(t)
	result.requireOutput(t, "runtime-happy-path")
}

func TestRuntimeWithoutInitializationScript(t *testing.T) {
	finalScript := `
		test ! -e /workspace/.test
		echo initialization-skipped
	`
	mipeCommand := mipeCommand(defaultConfigPath(), finalScript)
	result := runContainer(t, containerSpec{
		command: rootSetup(`rm -f /workspace/.mipe/init/setup.sh`, mipeCommand),
	})
	result.requireSuccess(t)
	result.requireOutput(t, "initialization-skipped")
}

func TestRuntimeStopsWhenInitializationFails(t *testing.T) {
	result := runContainer(t, containerSpec{
		files: map[string]string{
			"/tmp/setup.sh": "setup_project() { echo init-sentinel >&2; return 23; }\n",
		},
		command: mipeCommand(defaultConfigPath(), `echo final-command-ran`),
	})
	result.requireFailure(t)
	result.requireOutput(t, "init-sentinel")
	result.rejectOutput(t, "final-command-ran")
}

func TestRuntimeWithoutAgentHome(t *testing.T) {
	config := defaultConfig()
	config.AgentHome = ""
	result := runWithConfig(t, config, mipeCommand(defaultConfigPath(), `
		test -z "${AGENT_HOME:-}"
		test -f /workspace/.test
		echo agent-home-skipped
	`))
	result.requireSuccess(t)
	result.requireOutput(t, "agent-home-skipped")
}

func TestPreparationPreservesStateAndUpdatesOwnership(t *testing.T) {
	script := `mkdir -p /home/dev/.mipe-agent/cache; touch /home/dev/.mipe-agent/cache/state`
	finalScript := `
		test -w /home/dev/.mipe-agent/cache/state
		test -f /home/dev/.mipe-agent/cache/state
		test -f /home/dev/.mipe-agent/AGENTS.md
		test "$(stat -c %u:%g /home/dev/.mipe-agent/cache/state)" = 1000:1000
		echo state-preserved
	`
	mipeCommand := mipeCommand(defaultConfigPath(), finalScript)
	result := runContainer(t, containerSpec{
		command: rootSetup(script, mipeCommand),
	})
	result.requireSuccess(t)
	result.requireOutput(t, "state-preserved")
}

func TestFinalCommandExitCodeIsPropagated(t *testing.T) {
	result := runContainer(t, containerSpec{command: mipeCommand(defaultConfigPath(), `exit 42`)})
	if result.exitCode != 42 {
		t.Fatalf("exit code = %d, want 42\noutput:\n%s", result.exitCode, result.output)
	}
}

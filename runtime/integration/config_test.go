package integration_test

import "testing"

func TestEnvironmentOverridesFileConfiguration(t *testing.T) {
	result := runContainer(t, containerSpec{
		environment: map[string]string{
			"USER_HOME":    "/home/override",
			"AGENT_HOME":   "/home/override/agent",
			"RUNTIME_HOME": "/opt/mipe",
			"WORKSPACE":    "/override-workspace",
			"LOCAL_UID":    "2000",
			"LOCAL_GID":    "2001",
		},
		command: rootSetup(`install -d -o 2000 -g 2001 /home/override /override-workspace`, mipeCommand(defaultConfigPath(), `
			test "$(id -u)" = 2000
			test "$(id -g)" = 2001
			test "$HOME" = /home/override
			test "$AGENT_HOME" = /home/override/agent
			test "$PWD" = /override-workspace
			test -f /home/override/agent/AGENTS.md
			echo environment-overrides
		`)),
	})
	result.requireSuccess(t)
	result.requireOutput(t, "environment-overrides")
}

func TestInvalidConfiguration(t *testing.T) {
	tests := []struct {
		name       string
		contents   string
		command    []string
		wantOutput string
	}{
		{
			name:     "malformed JSON",
			contents: `{`, command: []string{"mipe", "--config", defaultConfigPath(), "true"},
			wantOutput: "configuration file error",
		},
		{
			name:       "missing workspace",
			contents:   encodedConfig(t, editConfig(func(config *runtimeConfig) { config.Workspace = "" })),
			command:    []string{"mipe", "--config", defaultConfigPath(), "true"},
			wantOutput: "configuration missing required value",
		},
		{
			name:       "nonnumeric UID",
			contents:   encodedConfig(t, editConfig(func(config *runtimeConfig) { config.LocalUID = "invalid" })),
			command:    []string{"env", "-u", "LOCAL_UID", "-u", "LOCAL_GID", "mipe", "--config", defaultConfigPath(), "true"},
			wantOutput: "configuration contains invalid value",
		},
		{
			name:       "missing command",
			contents:   encodedConfig(t, defaultConfig()),
			command:    []string{"mipe", "--config", defaultConfigPath()},
			wantOutput: "configuration missing required value",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := runContainer(t, containerSpec{
				files:   map[string]string{defaultConfigPath(): test.contents},
				command: test.command,
			})
			result.requireFailure(t)
			result.requireOutput(t, test.wantOutput)
		})
	}
}

func TestInvalidWorkspace(t *testing.T) {
	tests := []struct {
		name      string
		workspace string
		setup     string
	}{
		{
			name:      "missing",
			workspace: "/missing-workspace",
		},
		{
			name:      "file",
			workspace: "/workspace-file",
			setup:     `touch /workspace-file`,
		},
		{
			name:      "not writable",
			workspace: "/locked-workspace",
			setup:     `install -d -m 0700 -o root -g root /locked-workspace`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			config := defaultConfig()
			config.Workspace = test.workspace
			command := mipeCommand(defaultConfigPath(), `echo invalid-workspace-final-command`)
			if test.setup != "" {
				command = rootSetup(test.setup, command)
			}
			result := runWithConfig(t, config, command)
			result.requireFailure(t)
			result.requireOutput(t, "not a writable directory")
			result.rejectOutput(t, "invalid-workspace-final-command")
		})
	}
}

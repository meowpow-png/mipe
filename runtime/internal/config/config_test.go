package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNew_DerivesAgentHomeAndPreservesCommand(t *testing.T) {
	t.Parallel()

	command := []string{"bash", "-lc", "echo ok"}
	cfg := New(Environment{
		AgentName:   "codex",
		Home:        "/home/user",
		RuntimeHome: "/runtime",
		Workspace:   "/workspace",
		LocalUID:    "1000",
		LocalGID:    "1001",
	}, command)

	if cfg.AgentHome != filepath.Join("/home/user", ".codex") {
		t.Fatalf("AgentHome = %q, want /home/user/.codex", cfg.AgentHome)
	}
	if !reflect.DeepEqual(cfg.Command, command) {
		t.Fatalf("Command = %#v, want %#v", cfg.Command, command)
	}
}

func TestNew_LeavesAgentHomeEmptyWithoutHomeOrAgentName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		env  Environment
	}{
		{name: "missing home", env: Environment{AgentName: "codex"}},
		{name: "missing agent", env: Environment{Home: "/home/user"}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := New(tt.env, nil)
			if cfg.AgentHome != "" {
				t.Fatalf("AgentHome = %q, want empty", cfg.AgentHome)
			}
		})
	}
}

func TestLoad_LoadsFileAndAppliesEnvironmentOverrides(t *testing.T) {
	t.Setenv("AGENT_NAME", "process-agent")
	t.Setenv("HOME", "/home/process")
	t.Setenv("LOCAL_UID", "2000")

	path := filepath.Join(t.TempDir(), "config.json")
	content := `{
		"agent_name": "file-agent",
		"home": "/home/file",
		"runtime_home": "/runtime/file",
		"workspace": "/workspace/file",
		"local_uid": "1000",
		"local_gid": "1001"
	}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load([]string{"--config", path, "bash", "-lc", "echo ok"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.AgentName != "process-agent" {
		t.Fatalf("AgentName = %q, want process-agent", cfg.AgentName)
	}
	if cfg.LocalUID != "2000" {
		t.Fatalf("LocalUID = %q, want 2000", cfg.LocalUID)
	}
	if cfg.LocalGID != "1001" {
		t.Fatalf("LocalGID = %q, want 1001", cfg.LocalGID)
	}
	if cfg.AgentHome != filepath.Join("/home/process", ".process-agent") {
		t.Fatalf("AgentHome = %q, want /home/process/.process-agent", cfg.AgentHome)
	}
	if want := []string{"bash", "-lc", "echo ok"}; !reflect.DeepEqual(cfg.Command, want) {
		t.Fatalf("Command = %#v, want %#v", cfg.Command, want)
	}
}

func TestLoad_LoadsDefaultConfigFromRuntimeHome(t *testing.T) {
	runtimeHome := t.TempDir()
	t.Setenv("RUNTIME_HOME", runtimeHome)
	t.Setenv("HOME", "/home/file")

	path := filepath.Join(runtimeHome, "config.json")
	content := `{
		"agent_name": "file-agent",
		"home": "/home/file",
		"workspace": "/workspace/file",
		"local_uid": "1000",
		"local_gid": "1001"
	}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load([]string{"bash", "-lc", "echo ok"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.AgentName != "file-agent" {
		t.Fatalf("AgentName = %q, want file-agent", cfg.AgentName)
	}
	if cfg.RuntimeHome != runtimeHome {
		t.Fatalf("RuntimeHome = %q, want %q", cfg.RuntimeHome, runtimeHome)
	}
	if cfg.AgentHome != filepath.Join("/home/file", ".file-agent") {
		t.Fatalf("AgentHome = %q, want /home/file/.file-agent", cfg.AgentHome)
	}
	if want := []string{"bash", "-lc", "echo ok"}; !reflect.DeepEqual(cfg.Command, want) {
		t.Fatalf("Command = %#v, want %#v", cfg.Command, want)
	}
}

func TestLoad_ExplicitConfigOverridesDefaultConfigPath(t *testing.T) {
	runtimeHome := t.TempDir()
	t.Setenv("RUNTIME_HOME", runtimeHome)

	defaultContent := `{
		"agent_name": "default-agent",
		"home": "/home/default",
		"workspace": "/workspace/default",
		"local_uid": "1000",
		"local_gid": "1001"
	}`
	if err := os.WriteFile(filepath.Join(runtimeHome, "config.json"), []byte(defaultContent), 0o600); err != nil {
		t.Fatalf("write default config: %v", err)
	}

	explicitPath := filepath.Join(t.TempDir(), "config.json")
	explicitContent := `{
		"agent_name": "explicit-agent",
		"home": "/home/explicit",
		"workspace": "/workspace/explicit",
		"local_uid": "2000",
		"local_gid": "2001"
	}`
	if err := os.WriteFile(explicitPath, []byte(explicitContent), 0o600); err != nil {
		t.Fatalf("write explicit config: %v", err)
	}
	cfg, err := Load([]string{"--config", explicitPath, "bash"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.AgentName != "explicit-agent" {
		t.Fatalf("AgentName = %q, want explicit-agent", cfg.AgentName)
	}
	if cfg.LocalUID != "2000" || cfg.LocalGID != "2001" {
		t.Fatalf("uid/gid = %q/%q, want 2000/2001", cfg.LocalUID, cfg.LocalGID)
	}
}

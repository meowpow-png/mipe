package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestNew_PreservesConfiguredAgentHomeAndCommand(t *testing.T) {
	t.Parallel()

	command := []string{"bash", "-lc", "echo ok"}
	cfg := New(Environment{
		AgentName:   "codex",
		Home:        "/home/user",
		AgentHome:   "/agent/home",
		RuntimeHome: "/runtime",
		Workspace:   "/workspace",
		LocalUID:    "1000",
		LocalGID:    "1001",
	}, command)

	if cfg.AgentName != "codex" {
		t.Fatalf("AgentName = %q, want codex", cfg.AgentName)
	}
	if cfg.AgentHome != "/agent/home" {
		t.Fatalf("AgentHome = %q, want /agent/home", cfg.AgentHome)
	}
	if !reflect.DeepEqual(cfg.Command, command) {
		t.Fatalf("Command = %#v, want %#v", cfg.Command, command)
	}
}

func TestNew_LeavesAgentHomeEmptyWhenUnset(t *testing.T) {
	t.Parallel()

	cfg := New(Environment{Home: "/home/user"}, nil)
	if cfg.AgentHome != "" {
		t.Fatalf("AgentHome = %q, want empty", cfg.AgentHome)
	}
}

func TestLoad_LoadsFileAndAppliesEnvironmentOverrides(t *testing.T) {
	t.Setenv("AGENT_NAME", "process-agent")
	t.Setenv("HOME", "/home/process")
	t.Setenv("AGENT_HOME", "/agent/process")
	t.Setenv("LOCAL_UID", "2000")

	path := filepath.Join(t.TempDir(), "config.json")
	content := `{
		"agent_name": "file-agent",
		"home": "/home/file",
		"agent_home": "/agent/file",
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
	if cfg.AgentHome != "/agent/process" {
		t.Fatalf("AgentHome = %q, want /agent/process", cfg.AgentHome)
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
		"agent_home": "/agent/file",
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
	if cfg.RuntimeHome != runtimeHome {
		t.Fatalf("RuntimeHome = %q, want %q", cfg.RuntimeHome, runtimeHome)
	}
	if cfg.AgentName != "file-agent" {
		t.Fatalf("AgentName = %q, want file-agent", cfg.AgentName)
	}
	if cfg.AgentHome != "/agent/file" {
		t.Fatalf("AgentHome = %q, want /agent/file", cfg.AgentHome)
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
		"agent_home": "/agent/default",
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
		"agent_home": "/agent/explicit",
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
	if cfg.AgentHome != "/agent/explicit" {
		t.Fatalf("AgentHome = %q, want /agent/explicit", cfg.AgentHome)
	}
	if cfg.AgentName != "explicit-agent" {
		t.Fatalf("AgentName = %q, want explicit-agent", cfg.AgentName)
	}
	if cfg.LocalUID != "2000" || cfg.LocalGID != "2001" {
		t.Fatalf("uid/gid = %q/%q, want 2000/2001", cfg.LocalUID, cfg.LocalGID)
	}
}

func TestLoad_SetsDebugFromFlag(t *testing.T) {
	cfg, err := Load([]string{"-debug", "bash"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.Debug {
		t.Fatal("Debug = false, want true")
	}
}

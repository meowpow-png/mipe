package config

import (
	"errors"
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
		UserHome:    "/home/user",
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

	cfg := New(Environment{UserHome: "/home/user"}, nil)
	if cfg.AgentHome != "" {
		t.Fatalf("AgentHome = %q, want empty", cfg.AgentHome)
	}
}

func TestLoad_LoadsFileAndAppliesEnvironmentOverrides(t *testing.T) {
	t.Setenv("AGENT_NAME", "process-agent")
	t.Setenv("HOME", "/home/ignored")
	t.Setenv("USER_HOME", "/home/process")
	t.Setenv("AGENT_HOME", "/agent/process")
	t.Setenv("LOCAL_UID", "2000")

	path := filepath.Join(t.TempDir(), "config.json")
	content := `{
		"agent_name": "file-agent",
		"user_home": "/home/file",
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
	if cfg.UserHome != "/home/process" {
		t.Fatalf("UserHome = %q, want /home/process", cfg.UserHome)
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

func TestConfigPath_UsesFixedDefaultIndependentOfRuntimeHome(t *testing.T) {
	t.Setenv("RUNTIME_HOME", "/ignored/runtime")

	if got := configPath(""); got != "/opt/mipe/config/config.json" {
		t.Fatalf("configPath() = %q, want /opt/mipe/config/config.json", got)
	}
}

func TestLoad_UsesExplicitConfigPath(t *testing.T) {
	explicitPath := filepath.Join(t.TempDir(), "config.json")
	explicitContent := `{
		"agent_name": "explicit-agent",
		"user_home": "/home/explicit",
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
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load([]string{"--config", path, "-debug", "bash"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.Debug {
		t.Fatal("Debug = false, want true")
	}
}

func TestLoad_VersionSkipsConfigFileAndHonorsDebug(t *testing.T) {
	t.Run("flag", func(t *testing.T) {
		cfg, err := Load([]string{"--version", "--debug", "--config", "/missing/config.json"})
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if !cfg.Version || !cfg.Debug {
			t.Fatalf("Config = %#v, want version and debug enabled", cfg)
		}
	})

	t.Run("environment", func(t *testing.T) {
		t.Setenv("MIPE_DEBUG", "true")

		cfg, err := Load([]string{"-v", "--config", "/missing/config.json"})
		if err != nil {
			t.Fatalf("Load() error = %v", err)
		}
		if !cfg.Version || !cfg.Debug {
			t.Fatalf("Config = %#v, want version and debug enabled", cfg)
		}
	})
}

func TestLoad_SetsDebugFromEnvironment(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("MIPE_DEBUG", "true")
	cfg, err := Load([]string{"--config", path, "bash"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.Debug {
		t.Fatal("Debug = false, want true")
	}
}

func TestLoad_CombinesDebugFlagAndEnvironment(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("MIPE_DEBUG", "false")
	cfg, err := Load([]string{"--config", path, "--debug", "bash"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if !cfg.Debug {
		t.Fatal("Debug = false, want true")
	}
}

func TestLoad_RejectsInvalidDebugEnvironmentValue(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	t.Setenv("MIPE_DEBUG", "enabled")
	_, err := Load([]string{"--config", path, "bash"})
	var invalidErr *InvalidValueError
	if !errors.As(err, &invalidErr) {
		t.Fatalf("Load() error = %T, want *InvalidValueError", err)
	}
	if invalidErr.Field != "MIPE_DEBUG" || invalidErr.Reason != "boolean" {
		t.Fatalf("InvalidValueError = %#v, want MIPE_DEBUG boolean", invalidErr)
	}
}

func TestLoad_SetsDefaultLogFormat(t *testing.T) {
	t.Setenv("MIPE_LOG_FORMAT", "")

	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load([]string{"--config", path, "bash"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.LogFormat != LogFormatConsole {
		t.Fatalf("LogFormat = %q, want %q", cfg.LogFormat, LogFormatConsole)
	}
}

func TestLoad_SetsLogFormatFromEnvironment(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	for _, format := range []string{LogFormatConsole, LogFormatJSON} {
		t.Run(format, func(t *testing.T) {
			t.Setenv("MIPE_LOG_FORMAT", format)
			cfg, err := Load([]string{"--config", path, "bash"})
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}
			if cfg.LogFormat != format {
				t.Fatalf("LogFormat = %q, want %q", cfg.LogFormat, format)
			}
		})
	}
}

func TestLoad_ProcessEnvironmentOverridesFileLogFormat(t *testing.T) {
	t.Setenv("MIPE_LOG_FORMAT", LogFormatConsole)

	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{"environment":{"MIPE_LOG_FORMAT":"json"}}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	cfg, err := Load([]string{"--config", path, "bash"})
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.LogFormat != LogFormatConsole {
		t.Fatalf("LogFormat = %q, want %q", cfg.LogFormat, LogFormatConsole)
	}
}

func TestLoad_RejectsInvalidLogFormat(t *testing.T) {
	t.Setenv("MIPE_LOG_FORMAT", "text")

	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, []byte(`{}`), 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}
	_, err := Load([]string{"--config", path, "bash"})
	var invalidErr *InvalidValueError
	if !errors.As(err, &invalidErr) {
		t.Fatalf("Load() error = %T, want *InvalidValueError", err)
	}
	if invalidErr.Field != "MIPE_LOG_FORMAT" || invalidErr.Reason != "must be console or json" {
		t.Fatalf("InvalidValueError = %#v, want MIPE_LOG_FORMAT must be console or json", invalidErr)
	}
}

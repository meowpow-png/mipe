package config

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestParseFlags_ReadsConfigAndCommand(t *testing.T) {
	t.Parallel()

	flags, err := ParseFlags([]string{"--config", "/config.json", "-debug", "bash", "-lc", "echo ok"})
	if err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}
	if flags.ConfigPath != "/config.json" {
		t.Fatalf("ConfigPath = %q, want /config.json", flags.ConfigPath)
	}
	if !flags.Debug {
		t.Fatal("Debug = false, want true")
	}
	if want := []string{"bash", "-lc", "echo ok"}; !reflect.DeepEqual(flags.Command, want) {
		t.Fatalf("Command = %#v, want %#v", flags.Command, want)
	}
}

func TestParseFlags_ReadsVersionAliases(t *testing.T) {
	t.Parallel()

	for _, arg := range []string{"--version", "-v"} {
		t.Run(arg, func(t *testing.T) {
			t.Parallel()

			flags, err := ParseFlags([]string{arg})
			if err != nil {
				t.Fatalf("ParseFlags() error = %v", err)
			}
			if !flags.Version {
				t.Fatal("Version = false, want true")
			}
		})
	}
}

func TestParseFlags_ReturnsFlagErrorForInvalidFlags(t *testing.T) {
	t.Parallel()

	_, err := ParseFlags([]string{"--unknown"})
	if err == nil {
		t.Fatal("ParseFlags() error = nil, want error")
	}
	if _, ok := errors.AsType[*FlagError](err); !ok {
		t.Fatalf("ParseFlags() error = %T, want *FlagError", err)
	}
}

func TestLoadFile_LoadsConfigFileValues(t *testing.T) {
	t.Parallel()

	t.Run("empty path", func(t *testing.T) {
		t.Parallel()

		values, err := LoadFile("")
		if err != nil {
			t.Fatalf("LoadFile() error = %v", err)
		}
		if values != nil {
			t.Fatalf("LoadFile() = %#v, want nil", values)
		}
	})

	t.Run("valid JSON", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "config.json")
		content := `{
			"environment": {"EXTRA": "value", "AGENT_NAME": "nested-agent", "AGENT_HOME": "/agent/nested"},
			"agent_name": "file-agent",
			"user_home": "/home/agent",
			"agent_home": "/agent/home",
			"runtime_home": "/runtime",
			"workspace": "/workspace",
			"local_uid": "1000",
			"local_gid": "1001"
		}`
		if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
			t.Fatalf("write config: %v", err)
		}
		values, err := LoadFile(path)
		if err != nil {
			t.Fatalf("LoadFile() error = %v", err)
		}
		want := map[string]string{
			"EXTRA":        "value",
			"AGENT_NAME":   "file-agent",
			"USER_HOME":    "/home/agent",
			"AGENT_HOME":   "/agent/home",
			"RUNTIME_HOME": "/runtime",
			"WORKSPACE":    "/workspace",
			"LOCAL_UID":    "1000",
			"LOCAL_GID":    "1001",
		}
		if !reflect.DeepEqual(values, want) {
			t.Fatalf("LoadFile() = %#v, want %#v", values, want)
		}
	})

	t.Run("open error", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "missing.json")

		_, err := LoadFile(path)
		var fileErr *FileError
		if !errors.As(err, &fileErr) {
			t.Fatalf("LoadFile() error = %T, want *FileError", err)
		}
		if fileErr.Path != path || fileErr.Operation != "open" {
			t.Fatalf("FileError = %#v, want path %q operation open", fileErr, path)
		}
	})

	t.Run("parse error", func(t *testing.T) {
		t.Parallel()

		path := filepath.Join(t.TempDir(), "config.json")
		if err := os.WriteFile(path, []byte("{"), 0o600); err != nil {
			t.Fatalf("write config: %v", err)
		}
		_, err := LoadFile(path)
		var fileErr *FileError
		if !errors.As(err, &fileErr) {
			t.Fatalf("LoadFile() error = %T, want *FileError", err)
		}
		if fileErr.Path != path || fileErr.Operation != "parse" {
			t.Fatalf("FileError = %#v, want path %q operation parse", fileErr, path)
		}
	})
}

func TestFileConfigEnvironmentValues_MergesNestedAndTopLevelValues(t *testing.T) {
	t.Parallel()

	cfg := fileConfig{
		Environment: map[string]string{
			"AGENT_NAME": "nested-agent",
			"AGENT_HOME": "/agent/nested",
			"KEEP":       "yes",
		},
		AgentName:   "top-agent",
		UserHome:    "/home/top",
		AgentHome:   "/agent/top",
		RuntimeHome: "",
		Workspace:   "/workspace",
		LocalUID:    "1000",
		LocalGID:    "1001",
	}
	values := cfg.EnvironmentValues()
	want := map[string]string{
		"AGENT_NAME": "top-agent",
		"AGENT_HOME": "/agent/top",
		"KEEP":       "yes",
		"USER_HOME":  "/home/top",
		"WORKSPACE":  "/workspace",
		"LOCAL_UID":  "1000",
		"LOCAL_GID":  "1001",
	}
	if !reflect.DeepEqual(values, want) {
		t.Fatalf("EnvironmentValues() = %#v, want %#v", values, want)
	}
}

func TestSetIfPresent_WritesOnlyNonEmptyValues(t *testing.T) {
	t.Parallel()

	values := map[string]string{"EXISTING": "keep"}

	setIfPresent(values, "EMPTY", "")
	setIfPresent(values, "VALUE", "present")

	if _, ok := values["EMPTY"]; ok {
		t.Fatal("setIfPresent() wrote empty value")
	}
	if got, want := values["VALUE"], "present"; got != want {
		t.Fatalf("VALUE = %q, want %q", got, want)
	}
}

package validation

import (
	"errors"
	"testing"

	"github.com/meowpow-png/mipe/runtime/internal/config"
)

func validConfig() config.Config {
	return config.Config{
		AgentName:   "codex",
		Home:        "/home/user",
		AgentHome:   "/home/user/.codex",
		RuntimeHome: "/runtime",
		Workspace:   "/workspace",
		LocalUID:    "1000",
		LocalGID:    "1001",
		Command:     []string{"bash"},
	}
}

func TestConfig_AcceptsValidConfiguration(t *testing.T) {
	t.Parallel()

	if err := Config(validConfig()); err != nil {
		t.Fatalf("Config() error = %v, want nil", err)
	}
}

func TestConfig_ReturnsMissingValueErrorForRequiredFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		field string
		edit  func(*config.Config)
	}{
		{name: "agent name", field: "agent_name", edit: func(cfg *config.Config) { cfg.AgentName = "" }},
		{name: "agent home", field: "agent_home", edit: func(cfg *config.Config) { cfg.AgentHome = "" }},
		{name: "home", field: "home", edit: func(cfg *config.Config) { cfg.Home = "" }},
		{name: "runtime home", field: "runtime_home", edit: func(cfg *config.Config) { cfg.RuntimeHome = "" }},
		{name: "workspace", field: "workspace", edit: func(cfg *config.Config) { cfg.Workspace = "" }},
		{name: "local uid", field: "local_uid", edit: func(cfg *config.Config) { cfg.LocalUID = "" }},
		{name: "local gid", field: "local_gid", edit: func(cfg *config.Config) { cfg.LocalGID = "" }},
		{name: "command", field: "command", edit: func(cfg *config.Config) { cfg.Command = nil }},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := validConfig()
			tt.edit(&cfg)

			err := Config(cfg)
			var missing *config.MissingValueError
			if !errors.As(err, &missing) {
				t.Fatalf("Config() error = %T, want *MissingValueError", err)
			}
			if missing.Field != tt.field {
				t.Fatalf("Field = %q, want %q", missing.Field, tt.field)
			}
		})
	}
}

func TestConfig_ReturnsInvalidValueErrorForNonNumericOwnership(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		field string
		edit  func(*config.Config)
	}{
		{name: "local uid", field: "local_uid", edit: func(cfg *config.Config) { cfg.LocalUID = "abc" }},
		{name: "local gid", field: "local_gid", edit: func(cfg *config.Config) { cfg.LocalGID = "abc" }},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			cfg := validConfig()
			tt.edit(&cfg)

			err := Config(cfg)
			var invalid *config.InvalidValueError
			if !errors.As(err, &invalid) {
				t.Fatalf("Config() error = %T, want *InvalidValueError", err)
			}
			if invalid.Field != tt.field {
				t.Fatalf("Field = %q, want %q", invalid.Field, tt.field)
			}
		})
	}
}

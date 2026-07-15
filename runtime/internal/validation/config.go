package validation

import (
	"fmt"
	"strconv"

	"github.com/meowpow-png/mipe/runtime/internal/config"
)

// Config validates bootstrap configuration
func Config(cfg config.Config) error {
	required := []struct {
		name  string
		value string
	}{
		{name: "agent_name", value: cfg.AgentName},
		{name: "agent_home", value: cfg.AgentHome},
		{name: "home", value: cfg.Home},
		{name: "runtime_home", value: cfg.RuntimeHome},
		{name: "workspace", value: cfg.Workspace},
		{name: "local_uid", value: cfg.LocalUID},
		{name: "local_gid", value: cfg.LocalGID},
	}
	for _, field := range required {
		if field.value == "" {
			return fmt.Errorf("required configuration value %s is missing", field.name)
		}
	}
	if len(cfg.Command) == 0 {
		return fmt.Errorf("required configuration value command is missing")
	}
	if _, err := strconv.Atoi(cfg.LocalUID); err != nil {
		return fmt.Errorf("local_uid must be a numeric user id: %w", err)
	}
	if _, err := strconv.Atoi(cfg.LocalGID); err != nil {
		return fmt.Errorf("local_gid must be a numeric group id: %w", err)
	}
	return nil
}

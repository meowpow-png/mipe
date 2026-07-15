package validation

import (
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
		{name: "home", value: cfg.Home},
		{name: "runtime_home", value: cfg.RuntimeHome},
		{name: "workspace", value: cfg.Workspace},
		{name: "local_uid", value: cfg.LocalUID},
		{name: "local_gid", value: cfg.LocalGID},
	}
	for _, field := range required {
		if field.value == "" {
			return &config.MissingValueError{Field: field.name}
		}
	}
	if len(cfg.Command) == 0 {
		return &config.MissingValueError{Field: "command"}
	}
	if _, err := strconv.Atoi(cfg.LocalUID); err != nil {
		return &config.InvalidValueError{Field: "local_uid", Reason: "must be a numeric user id", Err: err}
	}
	if _, err := strconv.Atoi(cfg.LocalGID); err != nil {
		return &config.InvalidValueError{Field: "local_gid", Reason: "must be a numeric group id", Err: err}
	}
	return nil
}

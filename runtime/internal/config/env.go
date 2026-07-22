package config

import (
	"maps"
	"os"
	"strconv"
	"strings"
)

type Environment struct {
	Values      map[string]string
	AgentName   string
	UserHome    string
	AgentHome   string
	RuntimeHome string
	Workspace   string
	LocalUID    string
	LocalGID    string
}

// LoadEnvironment loads bootstrap configuration from environment variables
func LoadEnvironment(defaults map[string]string) Environment {
	values := make(map[string]string)

	maps.Copy(values, defaults)
	for _, entry := range os.Environ() {
		key, value, found := strings.Cut(entry, "=")
		if !found {
			continue
		}
		values[key] = value
	}
	if values["LOCAL_UID"] == "" {
		values["LOCAL_UID"] = "1000"
	}
	if values["LOCAL_GID"] == "" {
		values["LOCAL_GID"] = "1000"
	}
	agentName := values["AGENT_NAME"]
	userHome := values["USER_HOME"]
	agentHome := values["AGENT_HOME"]
	runtimeHome := values["RUNTIME_HOME"]
	workspace := values["WORKSPACE"]

	return Environment{
		Values:      values,
		AgentName:   agentName,
		UserHome:    userHome,
		AgentHome:   agentHome,
		RuntimeHome: runtimeHome,
		Workspace:   workspace,
		LocalUID:    values["LOCAL_UID"],
		LocalGID:    values["LOCAL_GID"],
	}
}

func debugEnabled(values map[string]string) (bool, error) {
	value, ok := values["MIPE_DEBUG"]
	if !ok || value == "" {
		return false, nil
	}
	debug, err := strconv.ParseBool(value)
	if err != nil {
		return false, &InvalidValueError{Field: "MIPE_DEBUG", Reason: "boolean", Err: err}
	}
	return debug, nil
}

func logFormat(values map[string]string) (string, error) {
	value := values["MIPE_LOG_FORMAT"]
	if value == "" {
		return LogFormatConsole, nil
	}
	switch value {
	case LogFormatConsole, LogFormatJSON:
		return value, nil
	default:
		return "", &InvalidValueError{Field: "MIPE_LOG_FORMAT", Reason: "must be console or json"}
	}
}

package config

import (
	"os"
	"strings"
)

type Environment struct {
	Values      map[string]string
	AgentName   string
	Home        string
	RuntimeHome string
	Workspace   string
	LocalUID    string
	LocalGID    string
}

// LoadEnvironment loads bootstrap configuration from environment variables
func LoadEnvironment(defaults map[string]string) Environment {
	values := make(map[string]string)

	for key, value := range defaults {
		values[key] = value
	}
	for _, entry := range os.Environ() {
		key, value, found := strings.Cut(entry, "=")
		if !found {
			continue
		}

		values[key] = value
	}
	agentName := values["AGENT_NAME"]
	home := values["HOME"]
	runtimeHome := values["RUNTIME_HOME"]
	workspace := values["WORKSPACE"]

	return Environment{
		Values:      values,
		AgentName:   agentName,
		Home:        home,
		RuntimeHome: runtimeHome,
		Workspace:   workspace,
		LocalUID:    values["LOCAL_UID"],
		LocalGID:    values["LOCAL_GID"],
	}
}

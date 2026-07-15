package config

import (
	"fmt"
	"path/filepath"
)

type Config struct {
	AgentName   string
	Home        string
	AgentHome   string
	RuntimeHome string
	Workspace   string
	LocalUID    string
	LocalGID    string
	Command     []string
}

// Load loads bootstrap configuration from flags, file, and environment
func Load(args []string) (Config, error) {
	flags, err := ParseFlags(args)
	if err != nil {
		return Config{}, err
	}
	values, err := LoadFile(flags.ConfigPath)
	if err != nil {
		return Config{}, fmt.Errorf("load config file: %w", err)
	}
	return New(LoadEnvironment(values), flags.Command), nil
}

// New constructs a bootstrap configuration
func New(env Environment, command []string) Config {
	agentHome := ""
	if env.Home != "" && env.AgentName != "" {
		agentHome = filepath.Join(env.Home, "."+env.AgentName)
	}
	return Config{
		AgentName:   env.AgentName,
		Home:        env.Home,
		AgentHome:   agentHome,
		RuntimeHome: env.RuntimeHome,
		Workspace:   env.Workspace,
		LocalUID:    env.LocalUID,
		LocalGID:    env.LocalGID,
		Command:     command,
	}
}

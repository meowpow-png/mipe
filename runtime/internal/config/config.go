package config

import (
	"os"
	"path/filepath"
)

const defaultConfigFile = "config.json"

type Config struct {
	AgentName   string
	Home        string
	AgentHome   string
	RuntimeHome string
	Workspace   string
	LocalUID    string
	LocalGID    string
	Debug       bool
	Command     []string
}

// Load loads bootstrap configuration from flags, file, and environment
func Load(args []string) (Config, error) {
	flags, err := ParseFlags(args)
	if err != nil {
		return Config{}, err
	}
	values, err := LoadFile(configPath(flags.ConfigPath))
	if err != nil {
		return Config{}, err
	}
	cfg := New(LoadEnvironment(values), flags.Command)
	cfg.Debug = flags.Debug
	return cfg, nil
}

func configPath(path string) string {
	if path != "" {
		return path
	}
	runtimeHome := os.Getenv("RUNTIME_HOME")
	if runtimeHome == "" {
		return ""
	}
	return filepath.Join(runtimeHome, defaultConfigFile)
}

// New constructs a bootstrap configuration
func New(env Environment, command []string) Config {
	return Config{
		AgentName:   env.AgentName,
		Home:        env.Home,
		AgentHome:   env.AgentHome,
		RuntimeHome: env.RuntimeHome,
		Workspace:   env.Workspace,
		LocalUID:    env.LocalUID,
		LocalGID:    env.LocalGID,
		Command:     command,
	}
}

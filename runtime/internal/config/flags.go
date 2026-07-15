package config

import (
	"encoding/json"
	"flag"
	"io"
	"os"
)

type Flags struct {
	ConfigPath string
	Debug      bool
	Command    []string
}

type fileConfig struct {
	Environment map[string]string `json:"environment"`
	AgentName   string            `json:"agent_name"`
	Home        string            `json:"home"`
	RuntimeHome string            `json:"runtime_home"`
	Workspace   string            `json:"workspace"`
	LocalUID    string            `json:"local_uid"`
	LocalGID    string            `json:"local_gid"`
}

// ParseFlags parses bootstrap command-line flags
func ParseFlags(args []string) (Flags, error) {
	flags := flag.NewFlagSet("mipe", flag.ContinueOnError)
	flags.SetOutput(io.Discard)

	var parsed Flags
	flags.StringVar(&parsed.ConfigPath, "config", "", "path to bootstrap config file")
	flags.BoolVar(&parsed.Debug, "debug", false, "enable debug logging")

	if err := flags.Parse(args); err != nil {
		return Flags{}, &FlagError{Err: err}
	}
	parsed.Command = flags.Args()

	return parsed, nil
}

// LoadFile loads bootstrap configuration from a file
func LoadFile(path string) (map[string]string, error) {
	if path == "" {
		return nil, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, &FileError{Path: path, Operation: "open", Err: err}
	}
	//goland:noinspection GoUnhandledErrorResult
	defer file.Close()

	var cfg fileConfig
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, &FileError{Path: path, Operation: "parse", Err: err}
	}
	return cfg.EnvironmentValues(), nil
}

// EnvironmentValues converts file configuration into environment variables
func (cfg fileConfig) EnvironmentValues() map[string]string {
	values := make(map[string]string)

	for key, value := range cfg.Environment {
		values[key] = value
	}
	setIfPresent(values, "AGENT_NAME", cfg.AgentName)
	setIfPresent(values, "HOME", cfg.Home)
	setIfPresent(values, "RUNTIME_HOME", cfg.RuntimeHome)
	setIfPresent(values, "WORKSPACE", cfg.Workspace)
	setIfPresent(values, "LOCAL_UID", cfg.LocalUID)
	setIfPresent(values, "LOCAL_GID", cfg.LocalGID)

	return values
}

// setIfPresent adds a value when it is present
func setIfPresent(values map[string]string, key string, value string) {
	if value == "" {
		return
	}
	values[key] = value
}

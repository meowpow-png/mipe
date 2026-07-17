package config

const defaultConfigPath = "/opt/mipe/config.json"

const (
	LogFormatConsole = "console"
	LogFormatJSON    = "json"
)

type Config struct {
	AgentName   string
	UserHome    string
	AgentHome   string
	RuntimeHome string
	Workspace   string
	LocalUID    string
	LocalGID    string
	Debug       bool
	Version     bool
	LogFormat   string
	Command     []string
}

// Load loads bootstrap configuration from flags, file, and environment
func Load(args []string) (Config, error) {
	flags, err := ParseFlags(args)
	if err != nil {
		return Config{}, err
	}
	if flags.Version {
		env := LoadEnvironment(nil)
		debug, err := debugEnabled(env.Values)
		if err != nil {
			return Config{}, err
		}
		return Config{Debug: flags.Debug || debug, Version: true}, nil
	}
	values, err := LoadFile(configPath(flags.ConfigPath))
	if err != nil {
		return Config{}, err
	}
	env := LoadEnvironment(values)
	debug, err := debugEnabled(env.Values)
	if err != nil {
		return Config{}, err
	}
	logFormat, err := logFormat(env.Values)
	if err != nil {
		return Config{}, err
	}
	cfg := New(env, flags.Command)
	cfg.Debug = flags.Debug || debug
	cfg.LogFormat = logFormat
	return cfg, nil
}

func configPath(path string) string {
	if path != "" {
		return path
	}
	return defaultConfigPath
}

// New constructs a bootstrap configuration
func New(env Environment, command []string) Config {
	return Config{
		AgentName:   env.AgentName,
		UserHome:    env.UserHome,
		AgentHome:   env.AgentHome,
		RuntimeHome: env.RuntimeHome,
		Workspace:   env.Workspace,
		LocalUID:    env.LocalUID,
		LocalGID:    env.LocalGID,
		Command:     command,
	}
}

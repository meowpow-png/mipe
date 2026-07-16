package phase

import "github.com/meowpow-png/mipe/runtime/internal/config"

func runtimeEnvironment(cfg config.Config, extra ...string) []string {
	env := []string{
		"HOME=" + cfg.UserHome,
		"RUNTIME_HOME=" + cfg.RuntimeHome,
	}
	if cfg.AgentHome != "" {
		env = append(env, "AGENT_HOME="+cfg.AgentHome)
	}
	env = append(env, extra...)
	return env
}

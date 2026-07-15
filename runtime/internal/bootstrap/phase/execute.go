package phase

import (
	"fmt"
	"strings"

	"github.com/meowpow-png/mipe/runtime/internal/bootstrap/process"
	"github.com/meowpow-png/mipe/runtime/internal/config"
	"go.uber.org/zap"
)

// Execute executes the requested process
func Execute(cfg config.Config, logger *zap.Logger) error {
	logger.Info(
		"executing requested process",
		zap.String("user", fmt.Sprintf("%s:%s", cfg.LocalUID, cfg.LocalGID)),
		zap.Strings("command", cfg.Command),
	)
	_ = logger.Sync()

	args := append([]string{
		fmt.Sprintf("%s:%s", cfg.LocalUID, cfg.LocalGID),
		"env",
		"HOME=" + cfg.Home,
		agentHomeEnvironment(cfg),
		"RUNTIME_HOME=" + cfg.RuntimeHome,
	}, cfg.Command...)

	return process.Exec("gosu", args...)
}

func agentHomeEnvironment(cfg config.Config) string {
	name := strings.ToUpper(strings.ReplaceAll(cfg.AgentName, "-", "_"))

	return name + "_HOME=" + cfg.AgentHome
}

package phase

import (
	"fmt"

	"github.com/meowpow-png/mipe/runtime/internal/bootstrap/process"
	"github.com/meowpow-png/mipe/runtime/internal/config"
	"go.uber.org/zap"
)

var execProcess = process.ExecInDir

// Execute executes the requested process
func Execute(cfg config.Config, logger *zap.Logger) error {
	logger.Debug(
		"executing requested process",
		zap.String("user", fmt.Sprintf("%s:%s", cfg.LocalUID, cfg.LocalGID)),
		zap.Strings("command", cfg.Command),
	)
	_ = logger.Sync()

	args := runtimeEnvironment(cfg)
	args = append(args, cfg.Command...)

	name, args, err := commandAsUser(cfg.LocalUID, cfg.LocalGID, "env", args...)
	if err != nil {
		return err
	}
	return execProcess(cfg.Workspace, name, args...)
}

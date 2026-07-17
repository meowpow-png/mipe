package phase

import (
	"context"
	"fmt"

	"github.com/meowpow-png/mipe/runtime/internal/config"
	"github.com/meowpow-png/mipe/runtime/internal/validation"
	"go.uber.org/zap"
)

var workspaceWritable = checkWorkspaceWritable

// Validate validates runtime configuration
func Validate(cfg config.Config, logger *zap.Logger) error {
	logger.Debug("validating configuration")
	if err := validation.Config(cfg); err != nil {
		return err
	}
	if err := workspaceWritable(cfg); err != nil {
		return err
	}
	logger.Debug("configuration validated")

	return nil
}

func checkWorkspaceWritable(cfg config.Config) error {
	user := fmt.Sprintf("%s:%s", cfg.LocalUID, cfg.LocalGID)
	if err := runProcess(context.Background(), "gosu", user, "test", "-d", cfg.Workspace, "-a", "-w", cfg.Workspace); err != nil {
		return fmt.Errorf("workspace %q is not a writable directory for %s: %w", cfg.Workspace, user, err)
	}
	return nil
}

package phase

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/meowpow-png/mipe/runtime/internal/bootstrap/process"
	"github.com/meowpow-png/mipe/runtime/internal/config"
	"go.uber.org/zap"
)

var (
	statFile        = os.Stat
	runProcess      = process.Run
	runProcessInDir = process.RunInDir
)

// Initialize initializes the project
func Initialize(ctx context.Context, cfg config.Config, logger *zap.Logger) error {
	script := filepath.Join(cfg.Workspace, ".mipe", "init", "setup.sh")

	logger.Debug("checking project setup",
		zap.String("script", script),
	)
	if _, err := statFile(script); err != nil {
		if os.IsNotExist(err) {
			logger.Debug("project setup skipped", zap.String("script", script))
			return nil
		}
		return fmt.Errorf("check project setup script %q: %w", script, err)
	}
	logger.Debug("project setup started",
		zap.String("script", script),
	)
	args := runtimeEnvironment(cfg,
		"WORKSPACE="+cfg.Workspace,
		"SETUP_SCRIPT="+script,
	)
	args = append(args,
		"bash",
		"-c",
		`set -euo pipefail; source "$SETUP_SCRIPT"; setup_project`,
	)
	name, args, err := commandAsUser(cfg.LocalUID, cfg.LocalGID, "env", args...)
	if err != nil {
		return err
	}
	if err := runProcessInDir(ctx, cfg.Workspace, name, args...); err != nil {
		return fmt.Errorf("run project setup: %w", err)
	}
	logger.Debug("project setup completed",
		zap.String("script", script),
	)
	return nil
}

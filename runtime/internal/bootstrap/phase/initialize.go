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
	script := filepath.Join(cfg.Workspace, ".mipe", "init", "dependencies.sh")

	logger.Debug("checking project dependency initialization",
		zap.String("script", script),
	)
	if _, err := statFile(script); err != nil {
		if os.IsNotExist(err) {
			logger.Debug("project dependency initialization skipped", zap.String("script", script))
			return nil
		}
		return fmt.Errorf("check project dependency initialization script %q: %w", script, err)
	}
	logger.Debug("project dependency initialization started",
		zap.String("script", script),
	)
	args := []string{
		fmt.Sprintf("%s:%s", cfg.LocalUID, cfg.LocalGID),
		"env",
	}
	args = append(args, runtimeEnvironment(cfg,
		"WORKSPACE="+cfg.Workspace,
		"DEPENDENCIES_SCRIPT="+script,
	)...)
	args = append(args,
		"bash",
		"-c",
		`set -euo pipefail; source "$DEPENDENCIES_SCRIPT"; install_dependencies`,
	)
	if err := runProcessInDir(ctx, cfg.Workspace, "gosu", args...); err != nil {
		return fmt.Errorf("initialize project dependencies: %w", err)
	}
	logger.Debug("project dependency initialization completed",
		zap.String("script", script),
	)
	return nil
}

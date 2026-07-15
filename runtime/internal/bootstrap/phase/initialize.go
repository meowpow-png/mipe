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
	statFile   = os.Stat
	runProcess = process.Run
)

// Initialize initializes the project
func Initialize(ctx context.Context, cfg config.Config, logger *zap.Logger) error {
	script := filepath.Join(cfg.Workspace, ".codex", "init", "dependencies.sh")

	logger.Info("checking project dependency initialization",
		zap.String("script", script),
	)
	if _, err := statFile(script); err != nil {
		if os.IsNotExist(err) {
			logger.Info("project dependency initialization skipped", zap.String("script", script))
			return nil
		}
		return fmt.Errorf("check project dependency initialization script %q: %w", script, err)
	}
	logger.Info("project dependency initialization started",
		zap.String("script", script),
	)
	if err := runProcess(
		ctx,
		"gosu",
		fmt.Sprintf("%s:%s", cfg.LocalUID, cfg.LocalGID),
		"env",
		"HOME="+cfg.Home,
		"CODEX_HOME="+cfg.AgentHome,
		"RUNTIME_HOME="+cfg.RuntimeHome,
		"WORKSPACE="+cfg.Workspace,
		"DEPENDENCIES_SCRIPT="+script,
		"bash",
		"-c",
		`set -euo pipefail; source "$DEPENDENCIES_SCRIPT"; install_dependencies`,
	); err != nil {
		return fmt.Errorf("initialize project dependencies: %w", err)
	}
	logger.Info("project dependency initialization completed",
		zap.String("script", script),
	)
	return nil
}

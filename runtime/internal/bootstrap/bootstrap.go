package bootstrap

import (
	"context"

	"github.com/meowpow-png/mipe/runtime/internal/bootstrap/phase"
	"github.com/meowpow-png/mipe/runtime/internal/config"
	"go.uber.org/zap"
)

func Run(ctx context.Context, cfg config.Config, logger *zap.Logger) error {
	logger.Info("validation started")
	if err := phase.Validate(cfg, logger); err != nil {
		return err
	}
	logger.Info("validation completed")

	logger.Info("preparation started")
	if err := phase.Prepare(cfg, logger); err != nil {
		return err
	}
	logger.Info("preparation completed")

	logger.Info("initialization started")
	if err := phase.Initialize(ctx, cfg, logger); err != nil {
		return err
	}
	logger.Info("initialization completed")

	logger.Info("execution started")
	if err := phase.Execute(cfg, logger); err != nil {
		return err
	}
	logger.Info("execution completed")

	return nil
}

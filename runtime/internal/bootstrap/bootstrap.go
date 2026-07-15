package bootstrap

import (
	"context"

	"github.com/meowpow-png/mipe/runtime/internal/bootstrap/phase"
	"github.com/meowpow-png/mipe/runtime/internal/config"
	"go.uber.org/zap"
)

type phases struct {
	validate   func(config.Config, *zap.Logger) error
	prepare    func(config.Config, *zap.Logger) error
	initialize func(context.Context, config.Config, *zap.Logger) error
	execute    func(config.Config, *zap.Logger) error
}

var defaultPhases = phases{
	validate:   phase.Validate,
	prepare:    phase.Prepare,
	initialize: phase.Initialize,
	execute:    phase.Execute,
}

func Run(ctx context.Context, cfg config.Config, logger *zap.Logger) error {
	return run(ctx, cfg, logger, defaultPhases)
}

func run(ctx context.Context, cfg config.Config, logger *zap.Logger, phases phases) error {
	logger.Info("validation started")
	if err := phases.validate(cfg, logger); err != nil {
		return err
	}
	logger.Info("validation completed")

	logger.Info("preparation started")
	if err := phases.prepare(cfg, logger); err != nil {
		return err
	}
	logger.Info("preparation completed")

	logger.Info("initialization started")
	if err := phases.initialize(ctx, cfg, logger); err != nil {
		return err
	}
	logger.Info("initialization completed")

	logger.Info("execution started")
	if err := phases.execute(cfg, logger); err != nil {
		return err
	}
	logger.Info("execution completed")

	return nil
}

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
	logger.Info("Welcome to Mipe")
	logger.Debug("bootstrap configuration",
		zap.String("agent_name", cfg.AgentName),
		zap.String("user_home", cfg.UserHome),
		zap.String("agent_home", cfg.AgentHome),
		zap.String("runtime_home", cfg.RuntimeHome),
		zap.String("workspace", cfg.Workspace),
		zap.String("local_uid", cfg.LocalUID),
		zap.String("local_gid", cfg.LocalGID),
		zap.Strings("command", cfg.Command),
	)
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

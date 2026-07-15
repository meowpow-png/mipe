package main

import (
	"context"
	"errors"
	"os"

	"github.com/meowpow-png/mipe/runtime/internal/bootstrap"
	"github.com/meowpow-png/mipe/runtime/internal/config"
	"go.uber.org/zap"
)

func main() {
	logger, err := newLogger(false)
	if err != nil {
		panic(err)
	}
	cfg, err := config.Load(os.Args[1:])
	if err != nil {
		logConfigError(logger, err)
		_ = logger.Sync()
		os.Exit(1)
	}
	logger, err = newLogger(cfg.Debug)
	if err != nil {
		panic(err)
	}
	if err := bootstrap.Run(context.Background(), cfg, logger); err != nil {
		if isConfigError(err) {
			logConfigError(logger, err)
			_ = logger.Sync()
			os.Exit(1)
		}
		logger.Error("bootstrap failed", zap.Error(err))
		_ = logger.Sync()
		os.Exit(1)
	}
	_ = logger.Sync()
}

func newLogger(debug bool) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	if debug {
		cfg.Level.SetLevel(zap.DebugLevel)
	}
	return cfg.Build()
}

func logConfigError(logger *zap.Logger, err error) {
	if _, ok := errors.AsType[*config.FlagError](err); ok {
		logger.Error("configuration flags error", zap.Error(err))
		return
	}
	if fileErr, ok := errors.AsType[*config.FileError](err); ok {
		logger.Error("configuration file error",
			zap.String("path", fileErr.Path),
			zap.String("operation", fileErr.Operation),
			zap.Error(err),
		)
		return
	}
	if missingErr, ok := errors.AsType[*config.MissingValueError](err); ok {
		logger.Error("configuration missing required value",
			zap.String("field", missingErr.Field),
			zap.Error(err),
		)
		return
	}
	if invalidErr, ok := errors.AsType[*config.InvalidValueError](err); ok {
		logger.Error("configuration contains invalid value",
			zap.String("field", invalidErr.Field),
			zap.String("reason", invalidErr.Reason),
			zap.Error(err),
		)
		return
	}
	logger.Error("configuration error", zap.Error(err))
}

func isConfigError(err error) bool {
	var flagErr *config.FlagError
	var fileErr *config.FileError
	var missingErr *config.MissingValueError
	var invalidErr *config.InvalidValueError

	return errors.As(err, &flagErr) ||
		errors.As(err, &fileErr) ||
		errors.As(err, &missingErr) ||
		errors.As(err, &invalidErr)
}

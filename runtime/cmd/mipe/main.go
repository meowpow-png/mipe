package main

import (
	"context"
	"os"

	"github.com/meowpow-png/mipe/runtime/internal/bootstrap"
	"github.com/meowpow-png/mipe/runtime/internal/config"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	cfg, err := config.Load(os.Args[1:])
	if err != nil {
		logger.Error("load config", zap.Error(err))
		_ = logger.Sync()
		os.Exit(1)
	}
	if err := bootstrap.Run(context.Background(), cfg, logger); err != nil {
		logger.Error("bootstrap failed", zap.Error(err))
		_ = logger.Sync()
		os.Exit(1)
	}
	_ = logger.Sync()
}

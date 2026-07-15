package phase

import (
	"github.com/meowpow-png/mipe/runtime/internal/config"
	"github.com/meowpow-png/mipe/runtime/internal/validation"
	"go.uber.org/zap"
)

// Validate validates runtime configuration
func Validate(cfg config.Config, logger *zap.Logger) error {
	logger.Info("validating configuration")
	if err := validation.Config(cfg); err != nil {
		return err
	}
	logger.Info("configuration validated")

	return nil
}

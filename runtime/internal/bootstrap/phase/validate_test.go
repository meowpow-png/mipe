package phase

import (
	"errors"
	"testing"

	"github.com/meowpow-png/mipe/runtime/internal/config"
	"go.uber.org/zap"
)

func TestValidate_AcceptsValidConfiguration(t *testing.T) {
	t.Parallel()

	if err := Validate(testConfig(), zap.NewNop()); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
}

func TestValidate_ReturnsConfigurationValidationError(t *testing.T) {
	t.Parallel()

	cfg := testConfig()
	cfg.Workspace = ""

	err := Validate(cfg, zap.NewNop())
	var missing *config.MissingValueError
	if !errors.As(err, &missing) {
		t.Fatalf("Validate() error = %T, want *MissingValueError", err)
	}
	if missing.Field != "workspace" {
		t.Fatalf("Field = %q, want workspace", missing.Field)
	}
}

package phase

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/meowpow-png/mipe/runtime/internal/config"
	"go.uber.org/zap"
)

func TestValidate_AcceptsValidConfiguration(t *testing.T) {
	restore := replaceValidateSeams(t)
	defer restore()

	var got config.Config
	workspaceWritable = func(cfg config.Config) error {
		got = cfg
		return nil
	}

	if err := Validate(testConfig(), zap.NewNop()); err != nil {
		t.Fatalf("Validate() error = %v, want nil", err)
	}
	if got.Workspace != "/workspace" {
		t.Fatalf("workspaceWritable() cfg = %#v, want workspace /workspace", got)
	}
}

func TestValidate_ReturnsConfigurationValidationError(t *testing.T) {
	restore := replaceValidateSeams(t)
	defer restore()

	cfg := testConfig()
	cfg.Workspace = ""
	workspaceWritable = func(config.Config) error {
		t.Fatal("workspaceWritable was called for invalid config")
		return nil
	}

	err := Validate(cfg, zap.NewNop())
	var missing *config.MissingValueError
	if !errors.As(err, &missing) {
		t.Fatalf("Validate() error = %T, want *MissingValueError", err)
	}
	if missing.Field != "workspace" {
		t.Fatalf("Field = %q, want workspace", missing.Field)
	}
}

func TestValidate_ReturnsWorkspaceWritableError(t *testing.T) {
	restore := replaceValidateSeams(t)
	defer restore()

	sentinel := errors.New("permission denied")
	workspaceWritable = func(config.Config) error {
		return sentinel
	}

	err := Validate(testConfig(), zap.NewNop())
	if !errors.Is(err, sentinel) {
		t.Fatalf("Validate() error = %v, want sentinel", err)
	}
}

func TestCheckWorkspaceWritable_RunsAsConfiguredUser(t *testing.T) {
	restore := replaceInitializeSeams(t)
	defer restore()

	var gotName string
	var gotArgs []string
	runProcess = func(ctx context.Context, name string, args ...string) error {
		gotName = name
		gotArgs = append([]string(nil), args...)
		return nil
	}

	if err := checkWorkspaceWritable(testConfig()); err != nil {
		t.Fatalf("checkWorkspaceWritable() error = %v", err)
	}
	if gotName != "gosu" {
		t.Fatalf("name = %q, want gosu", gotName)
	}
	want := []string{"1000:1001", "test", "-d", "/workspace", "-a", "-w", "/workspace"}
	if !reflect.DeepEqual(gotArgs, want) {
		t.Fatalf("args = %#v, want %#v", gotArgs, want)
	}
}

func TestCheckWorkspaceWritable_ReturnsClearError(t *testing.T) {
	restore := replaceInitializeSeams(t)
	defer restore()

	sentinel := errors.New("test failed")
	runProcess = func(ctx context.Context, name string, args ...string) error {
		return sentinel
	}

	err := checkWorkspaceWritable(testConfig())
	if !errors.Is(err, sentinel) {
		t.Fatalf("checkWorkspaceWritable() error = %v, want sentinel", err)
	}
	if !strings.Contains(err.Error(), `workspace "/workspace" is not a writable directory for 1000:1001`) {
		t.Fatalf("error = %q, want writable workspace message", err)
	}
}

func replaceValidateSeams(t *testing.T) func() {
	t.Helper()

	originalWorkspaceWritable := workspaceWritable

	return func() {
		workspaceWritable = originalWorkspaceWritable
	}
}

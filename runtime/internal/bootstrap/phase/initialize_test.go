package phase

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestInitialize_SkipsMissingDependencyScript(t *testing.T) {
	restore := replaceInitializeSeams(t)
	defer restore()

	statFile = func(path string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}
	runProcessInDir = func(ctx context.Context, dir string, name string, args ...string) error {
		t.Fatal("runProcess was called for missing script")
		return nil
	}
	if err := Initialize(context.Background(), testConfig(), zap.NewNop()); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
}

func TestInitialize_ReturnsStatError(t *testing.T) {
	restore := replaceInitializeSeams(t)
	defer restore()

	sentinel := errors.New("stat failed")
	statFile = func(path string) (os.FileInfo, error) {
		return nil, sentinel
	}
	err := Initialize(context.Background(), testConfig(), zap.NewNop())
	if !errors.Is(err, sentinel) {
		t.Fatalf("Initialize() error = %v, want sentinel", err)
	}
}

func TestInitialize_RunsExistingDependencyScript(t *testing.T) {
	restore := replaceInitializeSeams(t)
	defer restore()

	statFile = func(path string) (os.FileInfo, error) {
		return nil, nil
	}
	var gotName string
	var gotDir string
	var gotArgs []string
	runProcessInDir = func(ctx context.Context, dir string, name string, args ...string) error {
		gotDir = dir
		gotName = name
		gotArgs = append([]string(nil), args...)
		return nil
	}
	if err := Initialize(context.Background(), testConfig(), zap.NewNop()); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}
	if gotName != "gosu" {
		t.Fatalf("name = %q, want gosu", gotName)
	}
	if gotDir != "/workspace" {
		t.Fatalf("dir = %q, want /workspace", gotDir)
	}
	wantArgs := []string{
		"1000:1001",
		"env",
		"HOME=/home/user",
		"RUNTIME_HOME=/runtime",
		"AGENT_HOME=/agent/home",
		"WORKSPACE=/workspace",
		"DEPENDENCIES_SCRIPT=/workspace/.mipe/init/dependencies.sh",
		"bash",
		"-c",
		`set -euo pipefail; source "$DEPENDENCIES_SCRIPT"; install_dependencies`,
	}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("args = %#v, want %#v", gotArgs, wantArgs)
	}
}

func TestInitialize_ReturnsProcessRunError(t *testing.T) {
	restore := replaceInitializeSeams(t)
	defer restore()

	sentinel := errors.New("run failed")
	statFile = func(path string) (os.FileInfo, error) {
		return nil, nil
	}
	runProcessInDir = func(ctx context.Context, dir string, name string, args ...string) error {
		return sentinel
	}
	err := Initialize(context.Background(), testConfig(), zap.NewNop())
	if !errors.Is(err, sentinel) {
		t.Fatalf("Initialize() error = %v, want sentinel", err)
	}
}

func replaceInitializeSeams(t *testing.T) func() {
	t.Helper()

	originalStat := statFile
	originalRun := runProcess
	originalRunInDir := runProcessInDir

	return func() {
		statFile = originalStat
		runProcess = originalRun
		runProcessInDir = originalRunInDir
	}
}

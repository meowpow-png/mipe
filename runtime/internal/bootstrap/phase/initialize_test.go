package phase

import (
	"context"
	"errors"
	"os"
	"reflect"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestInitialize_SkipsMissingSetupScript(t *testing.T) {
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

func TestInitialize_LogsSkippedSetupScriptAtDebug(t *testing.T) {
	restore := replaceInitializeSeams(t)
	defer restore()

	statFile = func(path string) (os.FileInfo, error) {
		return nil, os.ErrNotExist
	}
	core, logs := observer.New(zapcore.DebugLevel)

	if err := Initialize(context.Background(), testConfig(), zap.New(core)); err != nil {
		t.Fatalf("Initialize() error = %v", err)
	}

	entries := logs.All()
	if len(entries) != 2 {
		t.Fatalf("log entry count = %d, want 2", len(entries))
	}
	for _, entry := range entries {
		if entry.Level != zapcore.DebugLevel {
			t.Fatalf("log level = %s, want DEBUG", entry.Level)
		}
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

func TestInitialize_RunsExistingSetupScript(t *testing.T) {
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
		"SETUP_SCRIPT=/workspace/.mipe/init/setup.sh",
		"bash",
		"-c",
		`set -euo pipefail; source "$SETUP_SCRIPT"; setup_project`,
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
	originalIdentity := currentIdentity
	currentIdentity = func() (int, int) { return 0, 0 }

	return func() {
		statFile = originalStat
		runProcess = originalRun
		runProcessInDir = originalRunInDir
		currentIdentity = originalIdentity
	}
}

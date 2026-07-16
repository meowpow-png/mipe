package process

import (
	"context"
	"errors"
	"os/exec"
	"reflect"
	"testing"
)

func TestRun_RunsCommandSuccessfully(t *testing.T) {
	t.Parallel()

	if err := Run(context.Background(), "go", "version"); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestRun_ReturnsCommandError(t *testing.T) {
	t.Parallel()

	err := Run(context.Background(), "go", "tool", "does-not-exist")
	if err == nil {
		t.Fatal("Run() error = nil, want error")
	}
	if _, ok := errors.AsType[*exec.ExitError](err); !ok {
		t.Fatalf("Run() error = %T, want *exec.ExitError", err)
	}
}

func TestRun_ReturnsContextCancellationError(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := Run(ctx, "go", "version")
	if err == nil {
		t.Fatal("Run() error = nil, want error")
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Run() error = %v, want context.Canceled", err)
	}
}

func TestExec_ReturnsLookupErrorForMissingExecutable(t *testing.T) {
	t.Parallel()

	err := Exec("__mipe_missing_executable__")
	if err == nil {
		t.Fatal("Exec() error = nil, want error")
	}
	if _, ok := errors.AsType[*exec.Error](err); !ok {
		t.Fatalf("Exec() error = %T, want *exec.Error", err)
	}
}

func TestExec_ReplacesProcessWithResolvedPath(t *testing.T) {
	restore := replaceExecSeams(t)
	defer restore()

	lookPath = func(name string) (string, error) {
		if name != "gosu" {
			t.Fatalf("lookPath name = %q, want gosu", name)
		}
		return "/usr/bin/gosu", nil
	}
	environ = func() []string {
		return []string{"A=B"}
	}
	var gotPath string
	var gotArgv []string
	var gotEnv []string
	execProcess = func(path string, argv []string, envv []string) error {
		gotPath = path
		gotArgv = append([]string(nil), argv...)
		gotEnv = append([]string(nil), envv...)
		return nil
	}

	if err := Exec("gosu", "1000:1001", "env", "bash"); err != nil {
		t.Fatalf("Exec() error = %v", err)
	}
	if gotPath != "/usr/bin/gosu" {
		t.Fatalf("path = %q, want /usr/bin/gosu", gotPath)
	}
	if want := []string{"/usr/bin/gosu", "1000:1001", "env", "bash"}; !reflect.DeepEqual(gotArgv, want) {
		t.Fatalf("argv = %#v, want %#v", gotArgv, want)
	}
	if want := []string{"A=B"}; !reflect.DeepEqual(gotEnv, want) {
		t.Fatalf("env = %#v, want %#v", gotEnv, want)
	}
}

func TestExec_ReturnsExecError(t *testing.T) {
	restore := replaceExecSeams(t)
	defer restore()

	sentinel := errors.New("exec failed")
	lookPath = func(name string) (string, error) {
		return "/usr/bin/gosu", nil
	}
	execProcess = func(path string, argv []string, envv []string) error {
		return sentinel
	}
	if err := Exec("gosu"); !errors.Is(err, sentinel) {
		t.Fatalf("Exec() error = %v, want sentinel", err)
	}
}

func TestExecInDir_ChangesDirectoryBeforeExec(t *testing.T) {
	restore := replaceExecSeams(t)
	defer restore()

	lookPath = func(name string) (string, error) { return "/usr/bin/gosu", nil }
	var gotDir string
	changeDir = func(dir string) error {
		gotDir = dir
		return nil
	}
	execProcess = func(path string, argv []string, envv []string) error { return nil }

	if err := ExecInDir("/work", "gosu"); err != nil {
		t.Fatalf("ExecInDir() error = %v", err)
	}
	if gotDir != "/work" {
		t.Fatalf("directory = %q, want /work", gotDir)
	}
}

func TestExecInDir_ReturnsChangeDirectoryError(t *testing.T) {
	restore := replaceExecSeams(t)
	defer restore()

	lookPath = func(name string) (string, error) { return "/usr/bin/gosu", nil }
	sentinel := errors.New("missing workspace")
	changeDir = func(dir string) error { return sentinel }
	execProcess = func(path string, argv []string, envv []string) error {
		t.Fatal("execProcess called after change directory failed")
		return nil
	}

	err := ExecInDir("/work", "gosu")
	if !errors.Is(err, sentinel) {
		t.Fatalf("ExecInDir() error = %v, want sentinel", err)
	}
}

func replaceExecSeams(t *testing.T) func() {
	t.Helper()

	originalLookPath := lookPath
	originalEnviron := environ
	originalChangeDir := changeDir
	originalExecProcess := execProcess

	return func() {
		lookPath = originalLookPath
		environ = originalEnviron
		changeDir = originalChangeDir
		execProcess = originalExecProcess
	}
}

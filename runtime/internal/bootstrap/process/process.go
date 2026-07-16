package process

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

var (
	lookPath    = exec.LookPath
	environ     = os.Environ
	changeDir   = os.Chdir
	execProcess = syscall.Exec
)

// Run runs a process and forwards its output to the current process
func Run(ctx context.Context, name string, args ...string) error {
	return RunInDir(ctx, "", name, args...)
}

// RunInDir runs a process in dir and forwards its output to the current process.
func RunInDir(ctx context.Context, dir string, name string, args ...string) error {
	command := exec.CommandContext(ctx, name, args...)
	command.Dir = dir

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

// Exec replaces the current process
func Exec(name string, args ...string) error {
	return ExecInDir("", name, args...)
}

// ExecInDir replaces the current process after changing to dir.
func ExecInDir(dir string, name string, args ...string) error {
	path, err := lookPath(name)
	if err != nil {
		return fmt.Errorf("locate %s: %w", name, err)
	}
	argv := append([]string{path}, args...)
	if dir != "" {
		if err := changeDir(dir); err != nil {
			return fmt.Errorf("change working directory to %q: %w", dir, err)
		}
	}

	return execProcess(path, argv, environ())
}

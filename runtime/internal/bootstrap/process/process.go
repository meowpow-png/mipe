package process

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// Run runs a process and forwards its output to the current process
func Run(ctx context.Context, name string, args ...string) error {
	command := exec.CommandContext(ctx, name, args...)

	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	return command.Run()
}

// Exec replaces the current process
func Exec(name string, args ...string) error {
	path, err := exec.LookPath(name)
	if err != nil {
		return fmt.Errorf("locate %s: %w", name, err)
	}
	argv := append([]string{path}, args...)

	return syscall.Exec(path, argv, os.Environ())
}

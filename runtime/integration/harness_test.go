package integration_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync/atomic"
	"testing"
)

const testImage = "mipe-runtime-test:latest"

const defaultInitializationScript = `#!/usr/bin/env bash
set -euo pipefail

install_dependencies() {
	touch "$WORKSPACE/.test"
}
`

var containerSequence atomic.Uint64

func TestMain(m *testing.M) {
	if os.Getenv("MIPE_INTEGRATION") != "1" {
		os.Exit(0)
	}
	command := exec.Command(
		"docker",
		"buildx",
		"bake",
		"--load",
		"--provenance=false",
		"--sbom=false",
		"test",
	)
	command.Dir = projectRoot()
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "build integration image: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}

func projectRoot() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("locate integration test source")
	}
	return filepath.Dir(filepath.Dir(file))
}

type runtimeConfig struct {
	AgentName   string `json:"agent_name"`
	UserHome    string `json:"user_home"`
	AgentHome   string `json:"agent_home,omitempty"`
	RuntimeHome string `json:"runtime_home"`
	Workspace   string `json:"workspace"`
	LocalUID    string `json:"local_uid"`
	LocalGID    string `json:"local_gid"`
}

type containerSpec struct {
	environment        map[string]string
	replaceEnvironment bool
	files              map[string]string
	command            []string
}

type containerResult struct {
	exitCode int
	output   string
}

func defaultConfig() runtimeConfig {
	return runtimeConfig{
		AgentName: "test-agent", UserHome: "/home/dev", AgentHome: "/home/dev/.mipe-agent",
		RuntimeHome: "/opt/mipe", Workspace: "/workspace", LocalUID: "1000", LocalGID: "1000",
	}
}

func editConfig(edit func(*runtimeConfig)) runtimeConfig {
	config := defaultConfig()
	edit(&config)
	return config
}

func encodedConfig(t *testing.T, config runtimeConfig) string {
	t.Helper()
	contents, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("encode config: %v", err)
	}
	return string(contents)
}

func defaultConfigPath() string { return "/tmp/config.json" }

func defaultEnvironment() map[string]string {
	return map[string]string{"LOCAL_UID": "1000", "LOCAL_GID": "1000"}
}

func mipeCommand(configPath, finalScript string) []string {
	return []string{"mipe", "--config", configPath, "bash", "-ceu", finalScript}
}

func rootSetup(script string, command []string) []string {
	quoted := make([]string, len(command))
	for index, argument := range command {
		quoted[index] = shellQuote(argument)
	}
	return []string{"bash", "-ceu", script + "\nexec " + strings.Join(quoted, " ")}
}

func shellQuote(value string) string {
	return "'" + strings.ReplaceAll(value, "'", `'"'"'`) + "'"
}

func runWithConfig(t *testing.T, config runtimeConfig, command []string) containerResult {
	t.Helper()
	return runContainer(t, containerSpec{
		files:   map[string]string{defaultConfigPath(): encodedConfig(t, config)},
		command: command,
	})
}

func runContainer(t *testing.T, spec containerSpec) containerResult {
	t.Helper()
	name := createContainer(t, spec)
	return startContainer(t, name)
}

func createContainer(t *testing.T, spec containerSpec) string {
	t.Helper()
	pid := strconv.Itoa(os.Getpid())
	tag := strconv.FormatUint(containerSequence.Add(1), 10)
	name := "mipe-runtime-integration-" + pid + "-" + tag
	files := mergeValues(map[string]string{
		defaultConfigPath():    encodedConfig(t, defaultConfig()),
		"/tmp/dependencies.sh": defaultInitializationScript,
	}, spec.files)
	spec.command = rootSetup(`
		install -d -m 0755 -o 1000 -g 1000 /workspace/.mipe/init
		install -m 0644 /tmp/dependencies.sh /workspace/.mipe/init/dependencies.sh
	`, spec.command)
	args := []string{"create", "--name", name}
	environment := spec.environment
	if !spec.replaceEnvironment {
		environment = mergeValues(defaultEnvironment(), environment)
	}
	for key, value := range environment {
		args = append(args, "--env", key+"="+value)
	}
	args = append(args, testImage)
	args = append(args, spec.command...)
	output, err := exec.Command("docker", args...).CombinedOutput()
	if err != nil {
		t.Fatalf("docker %s: %v\n%s", strings.Join(args, " "), err, output)
	}
	t.Cleanup(func() { removeContainer(t, name) })
	for destination, contents := range files {
		copyFileToContainer(t, name, destination, contents)
	}
	return name
}

func copyFileToContainer(t *testing.T, container, destination, contents string) {
	t.Helper()
	directory := t.TempDir()
	source := filepath.Join(directory, "fixture")
	if err := os.WriteFile(source, []byte(contents), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	output, err := exec.Command("docker", "cp", source, container+":"+destination).CombinedOutput()
	if err != nil {
		t.Fatalf("copy fixture to %s: %v\n%s", destination, err, output)
	}
}

func startContainer(t *testing.T, name string) containerResult {
	t.Helper()
	command := exec.Command("docker", "start", "--attach", name)
	var output bytes.Buffer
	command.Stdout = &output
	command.Stderr = &output
	err := command.Run()
	exitCode := 0
	if err != nil {
		var exitError *exec.ExitError
		if !errors.As(err, &exitError) {
			t.Fatalf("start container %s: %v\n%s", name, err, output.String())
		}
		exitCode = exitError.ExitCode()
	}
	return containerResult{exitCode: exitCode, output: output.String()}
}

func removeContainer(t *testing.T, name string) {
	t.Helper()
	if output, err := exec.Command("docker", "rm", "--force", name).CombinedOutput(); err != nil {
		t.Errorf("remove container %s: %v\n%s", name, err, output)
	}
}

func mergeValues(base, override map[string]string) map[string]string {
	merged := make(map[string]string, len(base)+len(override))
	for key, value := range base {
		merged[key] = value
	}
	for key, value := range override {
		merged[key] = value
	}
	return merged
}

func (result containerResult) requireSuccess(t *testing.T) {
	t.Helper()
	if result.exitCode != 0 {
		t.Fatalf("exit code = %d, want 0\noutput:\n%s", result.exitCode, result.output)
	}
}

func (result containerResult) requireFailure(t *testing.T) {
	t.Helper()
	if result.exitCode == 0 {
		t.Fatalf("exit code = 0, want nonzero\noutput:\n%s", result.output)
	}
}

func (result containerResult) requireOutput(t *testing.T, value string) {
	t.Helper()
	if !strings.Contains(result.output, value) {
		t.Fatalf("output does not contain %q:\n%s", value, result.output)
	}
}

func (result containerResult) rejectOutput(t *testing.T, value string) {
	t.Helper()
	if strings.Contains(result.output, value) {
		t.Fatalf("output unexpectedly contains %q:\n%s", value, result.output)
	}
}

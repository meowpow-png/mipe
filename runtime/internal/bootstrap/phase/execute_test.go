package phase

import (
	"errors"
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestRuntimeEnvironment_ExportsUserHomeAsHomeWithoutExposingUserHome(t *testing.T) {
	t.Parallel()

	cfg := testConfig()

	want := []string{
		"HOME=/home/user",
		"RUNTIME_HOME=/runtime",
		"AGENT_HOME=/agent/home",
	}
	if got := runtimeEnvironment(cfg); !reflect.DeepEqual(got, want) {
		t.Fatalf("runtimeEnvironment() = %#v, want %#v", got, want)
	}
}

func TestRuntimeEnvironment_SkipsAgentHomeWhenUnset(t *testing.T) {
	t.Parallel()

	cfg := testConfig()
	cfg.AgentHome = ""

	want := []string{
		"HOME=/home/user",
		"RUNTIME_HOME=/runtime",
	}
	if got := runtimeEnvironment(cfg); !reflect.DeepEqual(got, want) {
		t.Fatalf("runtimeEnvironment() = %#v, want %#v", got, want)
	}
}

func TestExecute_BuildsGosuExecInvocation(t *testing.T) {
	originalExec := execProcess
	defer func() { execProcess = originalExec }()

	var gotName string
	var gotArgs []string
	execProcess = func(name string, args ...string) error {
		gotName = name
		gotArgs = append([]string(nil), args...)
		return nil
	}
	if err := Execute(testConfig(), zap.NewNop()); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if gotName != "gosu" {
		t.Fatalf("name = %q, want gosu", gotName)
	}
	wantArgs := []string{
		"1000:1001",
		"env",
		"HOME=/home/user",
		"RUNTIME_HOME=/runtime",
		"AGENT_HOME=/agent/home",
		"bash",
	}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Fatalf("args = %#v, want %#v", gotArgs, wantArgs)
	}
}

func TestExecute_ReturnsExecError(t *testing.T) {
	originalExec := execProcess
	defer func() { execProcess = originalExec }()

	sentinel := errors.New("exec failed")
	execProcess = func(name string, args ...string) error {
		return sentinel
	}
	if err := Execute(testConfig(), zap.NewNop()); !errors.Is(err, sentinel) {
		t.Fatalf("Execute() error = %v, want sentinel", err)
	}
}

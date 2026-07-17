package bootstrap

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/meowpow-png/mipe/runtime/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestRun_ExecutesPhasesInOrder(t *testing.T) {
	t.Parallel()

	var calls []string
	cfg := config.Config{UserHome: "/home/user"}
	testPhases := phases{
		validate: func(got config.Config, logger *zap.Logger) error {
			calls = append(calls, "validate")
			if got.UserHome != cfg.UserHome {
				t.Fatalf("cfg = %#v, want %#v", got, cfg)
			}
			return nil
		},
		prepare: func(config.Config, *zap.Logger) error {
			calls = append(calls, "prepare")
			return nil
		},
		initialize: func(context.Context, config.Config, *zap.Logger) error {
			calls = append(calls, "initialize")
			return nil
		},
		execute: func(config.Config, *zap.Logger) error {
			calls = append(calls, "execute")
			return nil
		},
	}

	if err := run(context.Background(), cfg, zap.NewNop(), testPhases); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	want := []string{"validate", "prepare", "initialize", "execute"}
	if !reflect.DeepEqual(calls, want) {
		t.Fatalf("calls = %#v, want %#v", calls, want)
	}
}

func TestRun_StopsOnFirstPhaseError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		failCall  string
		wantCalls []string
	}{
		{name: "validate", failCall: "validate", wantCalls: []string{"validate"}},
		{name: "prepare", failCall: "prepare", wantCalls: []string{"validate", "prepare"}},
		{name: "initialize", failCall: "initialize", wantCalls: []string{"validate", "prepare", "initialize"}},
		{name: "execute", failCall: "execute", wantCalls: []string{"validate", "prepare", "initialize", "execute"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			sentinel := errors.New("boom")
			var calls []string
			fail := func(call string) error {
				calls = append(calls, call)
				if call == tt.failCall {
					return sentinel
				}
				return nil
			}
			testPhases := phases{
				validate: func(config.Config, *zap.Logger) error {
					return fail("validate")
				},
				prepare: func(config.Config, *zap.Logger) error {
					return fail("prepare")
				},
				initialize: func(context.Context, config.Config, *zap.Logger) error {
					return fail("initialize")
				},
				execute: func(config.Config, *zap.Logger) error {
					return fail("execute")
				},
			}
			err := run(context.Background(), config.Config{}, zap.NewNop(), testPhases)
			if !errors.Is(err, sentinel) {
				t.Fatalf("run() error = %v, want sentinel", err)
			}
			if !reflect.DeepEqual(calls, tt.wantCalls) {
				t.Fatalf("calls = %#v, want %#v", calls, tt.wantCalls)
			}
		})
	}
}

func TestRun_LogsLifecycleAtInfoAndConfigurationAtDebug(t *testing.T) {
	t.Parallel()

	core, logs := observer.New(zapcore.DebugLevel)
	cfg := config.Config{
		AgentName: "test-agent",
		LogFormat: config.LogFormatJSON,
	}
	noOpPhases := phases{
		validate:   func(config.Config, *zap.Logger) error { return nil },
		prepare:    func(config.Config, *zap.Logger) error { return nil },
		initialize: func(context.Context, config.Config, *zap.Logger) error { return nil },
		execute:    func(config.Config, *zap.Logger) error { return nil },
	}

	if err := run(context.Background(), cfg, zap.New(core), noOpPhases); err != nil {
		t.Fatalf("run() error = %v", err)
	}

	var infoCount int
	var debugEntries []observer.LoggedEntry
	for _, entry := range logs.All() {
		if entry.Level == zapcore.InfoLevel {
			infoCount++
		}
		if entry.Level == zapcore.DebugLevel {
			debugEntries = append(debugEntries, entry)
		}
	}
	if infoCount != 9 {
		t.Fatalf("INFO entry count = %d, want 9", infoCount)
	}
	if len(debugEntries) != 1 {
		t.Fatalf("DEBUG entry count = %d, want 1", len(debugEntries))
	}
	if got := debugEntries[0].ContextMap()["log_format"]; got != config.LogFormatJSON {
		t.Fatalf("log_format = %#v, want %q", got, config.LogFormatJSON)
	}
}

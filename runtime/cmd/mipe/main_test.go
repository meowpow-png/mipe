package main

import (
	"errors"
	"testing"

	"github.com/meowpow-png/mipe/runtime/internal/build"
	"github.com/meowpow-png/mipe/runtime/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

func TestVersionOutput(t *testing.T) {
	want := "Mipe version " + build.Version + "\n"
	if got := versionOutput(); got != want {
		t.Fatalf("versionOutput() = %q, want %q", got, want)
	}
}

func TestIsConfigError_IdentifiesTypedConfigErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "flag", err: &config.FlagError{Err: errors.New("bad flag")}, want: true},
		{name: "file", err: &config.FileError{Path: "config.json", Operation: "open", Err: errors.New("missing")}, want: true},
		{name: "missing", err: &config.MissingValueError{Field: "workspace"}, want: true},
		{name: "invalid", err: &config.InvalidValueError{Field: "local_uid", Reason: "numeric"}, want: true},
		{name: "wrapped", err: errors.Join(errors.New("outer"), &config.MissingValueError{Field: "command"}), want: true},
		{name: "other", err: errors.New("other"), want: false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := isConfigError(tt.err); got != tt.want {
				t.Fatalf("isConfigError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogConfigError_LogsTypedFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		err    error
		fields map[string]string
	}{
		{
			name: "flag",
			err:  &config.FlagError{Err: errors.New("bad flag")},
		},
		{
			name:   "file",
			err:    &config.FileError{Path: "config.json", Operation: "open", Err: errors.New("missing")},
			fields: map[string]string{"path": "config.json", "operation": "open"},
		},
		{
			name:   "missing",
			err:    &config.MissingValueError{Field: "workspace"},
			fields: map[string]string{"field": "workspace"},
		},
		{
			name:   "invalid",
			err:    &config.InvalidValueError{Field: "local_uid", Reason: "numeric"},
			fields: map[string]string{"field": "local_uid", "reason": "numeric"},
		},
		{
			name: "fallback",
			err:  errors.New("other"),
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			core, logs := observer.New(zapcore.ErrorLevel)
			logger := zap.New(core)

			logConfigError(logger, tt.err)

			entries := logs.All()
			if len(entries) != 1 {
				t.Fatalf("log entries = %d, want 1", len(entries))
			}
			context := entries[0].ContextMap()
			for key, want := range tt.fields {
				if got := context[key]; got != want {
					t.Fatalf("field %s = %#v, want %q", key, got, want)
				}
			}
			if _, ok := context["error"]; !ok {
				t.Fatal("error field is missing")
			}
		})
	}
}

func TestNewLogger_ConfiguresDebugLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		debug bool
		want  bool
	}{
		{name: "production", debug: false, want: false},
		{name: "debug", debug: true, want: true},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger, err := newLogger(tt.debug)
			if err != nil {
				t.Fatalf("newLogger() error = %v", err)
			}
			if got := logger.Core().Enabled(zapcore.DebugLevel); got != tt.want {
				t.Fatalf("debug enabled = %v, want %v", got, tt.want)
			}
		})
	}
}

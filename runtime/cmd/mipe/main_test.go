package main

import (
	"errors"
	"testing"

	"github.com/meowpow-png/mipe/runtime/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

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

func TestLogConfigError_LogsDistinctMessagesAndFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		err     error
		message string
		fields  map[string]string
	}{
		{
			name:    "flag",
			err:     &config.FlagError{Err: errors.New("bad flag")},
			message: "configuration flags error",
		},
		{
			name:    "file",
			err:     &config.FileError{Path: "config.json", Operation: "open", Err: errors.New("missing")},
			message: "configuration file error",
			fields:  map[string]string{"path": "config.json", "operation": "open"},
		},
		{
			name:    "missing",
			err:     &config.MissingValueError{Field: "workspace"},
			message: "configuration missing required value",
			fields:  map[string]string{"field": "workspace"},
		},
		{
			name:    "invalid",
			err:     &config.InvalidValueError{Field: "local_uid", Reason: "numeric"},
			message: "configuration contains invalid value",
			fields:  map[string]string{"field": "local_uid", "reason": "numeric"},
		},
		{
			name:    "fallback",
			err:     errors.New("other"),
			message: "configuration error",
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
			if entries[0].Message != tt.message {
				t.Fatalf("message = %q, want %q", entries[0].Message, tt.message)
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

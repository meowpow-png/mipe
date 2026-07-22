package main

import (
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestConsoleEncoder_FormatsEntriesByModeAndLevel(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		debug      bool
		level      zapcore.Level
		wantPrefix string
	}{
		{name: "normal info", level: zapcore.InfoLevel, wantPrefix: "• "},
		{name: "normal error", level: zapcore.ErrorLevel, wantPrefix: "✖ "},
		{name: "debug info", debug: true, level: zapcore.InfoLevel, wantPrefix: "INFO  "},
		{name: "debug detail", debug: true, level: zapcore.DebugLevel, wantPrefix: "DEBUG "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			encoded, err := newConsoleEncoder(tt.debug).EncodeEntry(zapcore.Entry{
				Level:   tt.level,
				Message: "message",
				Time:    time.Date(2026, time.January, 1, 12, 0, 0, 0, time.UTC),
				Caller:  zapcore.EntryCaller{Defined: true, File: "source.go", Line: 1},
			}, []zapcore.Field{zap.String("path", "/workspace")})
			if err != nil {
				t.Fatalf("EncodeEntry() error = %v", err)
			}
			defer encoded.Free()

			output := encoded.String()
			if !strings.HasPrefix(output, tt.wantPrefix) {
				t.Fatalf("output = %q, want prefix %q", output, tt.wantPrefix)
			}
			if !strings.Contains(output, "(path=/workspace)") {
				t.Fatalf("output = %q, want parenthesized fields", output)
			}
			if strings.Contains(output, "2026") || strings.Contains(output, "source.go") {
				t.Fatalf("output = %q, want no timestamp or caller", output)
			}
		})
	}
}

func TestConsoleEncoder_FormatsFieldsInEmissionOrder(t *testing.T) {
	t.Parallel()

	encoded, err := newConsoleEncoder(true).EncodeEntry(zapcore.Entry{Message: "message"}, []zapcore.Field{
		zap.Int("uid", 1000),
		zap.Int("gid", 1001),
	})
	if err != nil {
		t.Fatalf("EncodeEntry() error = %v", err)
	}
	defer encoded.Free()

	if output := encoded.String(); !strings.Contains(output, "(uid=1000, gid=1001)") {
		t.Fatalf("output = %q, want fields in emission order", output)
	}
}

func TestConsoleEncoder_FormatsWarningWithStack(t *testing.T) {
	t.Parallel()

	encoded, err := newConsoleEncoder(false).EncodeEntry(zapcore.Entry{
		Level:   zapcore.WarnLevel,
		Message: "warning",
		Stack:   "stack trace",
	}, nil)
	if err != nil {
		t.Fatalf("EncodeEntry() error = %v", err)
	}
	defer encoded.Free()

	if got, want := encoded.String(), "⚠ Warning\nstack trace\n"; got != want {
		t.Fatalf("encoded output = %q, want %q", got, want)
	}
}

func TestConsoleEncoder_ClonePreservesMode(t *testing.T) {
	t.Parallel()

	original := newConsoleEncoder(true).(*consoleEncoder)
	clone := original.Clone().(*consoleEncoder)

	if !clone.debug {
		t.Fatal("clone debug mode = false, want true")
	}
	encoded, err := clone.EncodeEntry(zapcore.Entry{Level: zapcore.InfoLevel, Message: "message"}, nil)
	if err != nil {
		t.Fatalf("EncodeEntry() error = %v", err)
	}
	defer encoded.Free()

	if got, want := encoded.String(), "INFO  Message\n"; got != want {
		t.Fatalf("encoded output = %q, want %q", got, want)
	}
}

func TestConsoleHelpers_HandleEmptyValues(t *testing.T) {
	t.Parallel()

	if got := sentenceCase(""); got != "" {
		t.Fatalf("sentenceCase(\"\") = %q, want empty string", got)
	}

	output := consoleBufferPool.Get()
	defer output.Free()
	appendConsoleFields(output, nil)
	if got := output.String(); got != "" {
		t.Fatalf("appendConsoleFields(nil) wrote %q, want empty string", got)
	}
}

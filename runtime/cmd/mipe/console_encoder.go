package main

import (
	"fmt"
	"strings"
	"unicode"

	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

const (
	consoleEncodingNormal = "mipe-console"
	consoleEncodingDebug  = "mipe-console-debug"
)

var consoleBufferPool = buffer.NewPool()

func init() {
	_ = zap.RegisterEncoder(consoleEncodingNormal, func(zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return newConsoleEncoder(false), nil
	})
	_ = zap.RegisterEncoder(consoleEncodingDebug, func(zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return newConsoleEncoder(true), nil
	})
}

func consoleEncoding(debug bool) string {
	if debug {
		return consoleEncodingDebug
	}
	return consoleEncodingNormal
}

type consoleEncoder struct {
	zapcore.Encoder
	debug bool
}

func newConsoleEncoder(debug bool) zapcore.Encoder {
	return &consoleEncoder{
		Encoder: zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig()),
		debug:   debug,
	}
}

func (encoder *consoleEncoder) Clone() zapcore.Encoder {
	return &consoleEncoder{Encoder: encoder.Encoder.Clone(), debug: encoder.debug}
}

func (encoder *consoleEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	buffer := consoleBufferPool.Get()
	buffer.AppendString(consolePrefix(entry.Level, encoder.debug))
	buffer.AppendString(sentenceCase(entry.Message))
	appendConsoleFields(buffer, fields)
	buffer.AppendByte('\n')
	if entry.Stack != "" {
		buffer.AppendString(entry.Stack)
		buffer.AppendByte('\n')
	}
	return buffer, nil
}

func consolePrefix(level zapcore.Level, debug bool) string {
	if debug {
		return fmt.Sprintf("%-6s", strings.ToUpper(level.String()))
	}
	switch {
	case level >= zapcore.ErrorLevel:
		return "✖ "
	case level == zapcore.WarnLevel:
		return "⚠ "
	default:
		return "• "
	}
}

func sentenceCase(message string) string {
	for index, character := range message {
		return string(unicode.ToUpper(character)) + message[index+len(string(character)):]
	}
	return message
}

func appendConsoleFields(buffer *buffer.Buffer, fields []zapcore.Field) {
	if len(fields) == 0 {
		return
	}

	var rendered []string
	for _, field := range fields {
		valueEncoder := zapcore.NewMapObjectEncoder()
		field.AddTo(valueEncoder)
		value, ok := valueEncoder.Fields[field.Key]
		if !ok {
			continue
		}
		rendered = append(rendered, fmt.Sprintf("%s=%v", field.Key, value))
	}
	if len(rendered) == 0 {
		return
	}
	buffer.AppendString(" (")
	buffer.AppendString(strings.Join(rendered, ", "))
	buffer.AppendByte(')')
}

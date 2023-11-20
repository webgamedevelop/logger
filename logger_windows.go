package logger

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New return a new logr.Logger and its associated flush function.
func New(ctx context.Context, enab zapcore.Level, format string, cfg zapcore.EncoderConfig, writers ...zapcore.WriteSyncer) (logr.Logger, func()) {
	var wss []zapcore.WriteSyncer
	for _, writer := range writers {
		wss = append(wss, zapcore.Lock(writer))
	}

	encoder := zapcore.NewConsoleEncoder(cfg)
	if format == JSONLogFormat {
		encoder = zapcore.NewJSONEncoder(cfg)
	}
	return newLogger(ctx, enab, encoder, wss...)
}

// NewLogger creates a new logr.Logger and its associated flush function.
func newLogger(ctx context.Context, enabLevel zapcore.Level, encoder zapcore.Encoder, streams ...zapcore.WriteSyncer) (logr.Logger, func()) {
	var cores []zapcore.Core
	enab := zap.NewAtomicLevel()
	enab.SetLevel(enabLevel)
	for _, stream := range streams {
		core := zapcore.NewCore(encoder, stream, enab)
		cores = append(cores, core)
	}

	core := zapcore.NewTee(cores...)
	stackTraceLevel := zap.NewAtomicLevel()
	stackTraceLevel.SetLevel(zapcore.ErrorLevel)
	l := zap.New(core, zap.WithCaller(true), zap.AddStacktrace(stackTraceLevel))
	return zapr.NewLoggerWithOptions(l), func() { _ = l.Sync() }
}

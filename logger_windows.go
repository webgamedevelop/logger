package logger

import (
	"context"
	"os"
	"slices"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New return a new logr.Logger and its associated flush function.
func New(ctx context.Context, cfg zapcore.EncoderConfig, writers ...zapcore.WriteSyncer) (logr.Logger, func()) {
	f := defaultFile()
	if f.Filename != "" {
		writers = append(writers, f)
	}
	if stdout {
		if !slices.ContainsFunc(writers, func(syncer zapcore.WriteSyncer) bool {
			if syncer == os.Stdout {
				return true
			}
			return false
		}) {
			writers = append(writers, os.Stdout)
		}
	}
	if stderr {
		if !slices.ContainsFunc(writers, func(syncer zapcore.WriteSyncer) bool {
			if syncer == os.Stderr {
				return true
			}
			return false
		}) {
			writers = append(writers, os.Stderr)
		}
	}

	combinedWriteSyncer := zap.CombineWriteSyncers(writers...)
	writer = combinedWriteSyncer
	encoder := zapcore.NewConsoleEncoder(cfg)
	if format == jsonLogFormat {
		encoder = zapcore.NewJSONEncoder(cfg)
	}
	return newLogger(zapcore.Level(-enab), encoder, combinedWriteSyncer)
}

// NewLogger creates a new logr.Logger and its associated flush function.
func newLogger(enabLevel zapcore.Level, encoder zapcore.Encoder, writeSyncer zapcore.WriteSyncer) (logr.Logger, func()) {
	enab := zap.NewAtomicLevel()
	enab.SetLevel(enabLevel)
	core := zapcore.NewCore(encoder, writeSyncer, enab)
	stackTraceLevel := zap.NewAtomicLevel()
	stackTraceLevel.SetLevel(zapcore.ErrorLevel)
	l := zap.New(core, zap.WithCaller(true), zap.AddStacktrace(stackTraceLevel))
	return zapr.NewLoggerWithOptions(l), func() { _ = l.Sync() }
}

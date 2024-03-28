package logger

import (
	"context"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new logr.Logger and its associated flush function.
func newLogger(_ context.Context, enabLevel zapcore.Level, encoder zapcore.Encoder, writeSyncer zapcore.WriteSyncer) (logr.Logger, func(), error) {
	enab := zap.NewAtomicLevel()
	enab.SetLevel(enabLevel)
	core := zapcore.NewCore(encoder, writeSyncer, enab)
	stackTraceLevel := zap.NewAtomicLevel()
	stackTraceLevel.SetLevel(zapcore.ErrorLevel)
	l := zap.New(core, zap.WithCaller(true), zap.AddStacktrace(stackTraceLevel))
	return zapr.NewLoggerWithOptions(l), func() { _ = l.Sync() }, nil
}

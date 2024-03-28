//go:build !windows

package logger

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new logr.Logger and its associated flush function.
func newLogger(ctx context.Context, enabLevel zapcore.Level, encoder zapcore.Encoder, writeSyncer zapcore.WriteSyncer) (logr.Logger, func()) {
	enab := zap.NewAtomicLevel()
	enab.SetLevel(enabLevel)
	core := zapcore.NewCore(encoder, writeSyncer, enab)
	stackTraceLevel := zap.NewAtomicLevel()
	stackTraceLevel.SetLevel(zapcore.ErrorLevel)
	l := zap.New(core, zap.WithCaller(true), zap.AddStacktrace(stackTraceLevel))

	if ctx != nil {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGUSR1, syscall.SIGUSR2)
		go func(ctx context.Context) {
			defer signal.Stop(c)
			for {
				select {
				case sig := <-c:
					if sig == syscall.SIGUSR1 {
						if enab.Level() < zapcore.FatalLevel {
							enab.SetLevel(enab.Level() + 1)
						}
					} else {
						if enab.Level() > -127 {
							enab.SetLevel(enab.Level() - 1)
						}
					}
				case <-ctx.Done():
					return
				}
			}
		}(ctx)
	}

	return zapr.NewLoggerWithOptions(l), func() { _ = l.Sync() }
}

//go:build !windows

package logger

import (
	"context"
	"os"
	"os/signal"
	"slices"
	"syscall"

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

	var wss []zapcore.WriteSyncer
	for _, writer := range writers {
		wss = append(wss, zapcore.Lock(writer))
	}
	encoder := zapcore.NewConsoleEncoder(cfg)
	if format == jsonLogFormat {
		encoder = zapcore.NewJSONEncoder(cfg)
	}
	return newLogger(ctx, zapcore.Level(-enab), encoder, wss...)
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

package logger

import (
	"context"
	"os"
	"slices"

	"github.com/go-logr/logr"
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
	return newLogger(ctx, zapcore.Level(-enab), encoder, combinedWriteSyncer)
}

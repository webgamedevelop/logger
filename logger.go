package logger

import (
	"context"
	"os"
	"slices"
	"strconv"

	"github.com/go-logr/logr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/klog/v2"
)

// New return a new logr.Logger and its associated flush function.
func New(ctx context.Context, cfg zapcore.EncoderConfig, writers ...zapcore.WriteSyncer) (logr.Logger, func(), error) {
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

	klog.InitFlags(klogFS)
	if err := klogFS.Set("v", strconv.Itoa(klogV)); err != nil {
		return logr.Logger{}, nil, err
	}

	return newLogger(ctx, zapcore.Level(-enab), encoder, combinedWriteSyncer)
}

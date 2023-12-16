package logger

import (
	"go.uber.org/zap/zapcore"
)

// Supported klog formats
const (
	consoleLogFormat = "console"
	jsonLogFormat    = "json"
)

var DefaultEncoderConfig = zapcore.EncoderConfig{
	TimeKey:        "ts",
	MessageKey:     "msg",
	CallerKey:      "caller",
	LevelKey:       "level",
	NameKey:        "logger",
	FunctionKey:    zapcore.OmitKey,
	StacktraceKey:  "stacktrace",
	LineEnding:     zapcore.DefaultLineEnding,
	EncodeLevel:    zapcore.CapitalLevelEncoder,
	EncodeTime:     zapcore.ISO8601TimeEncoder,
	EncodeDuration: zapcore.StringDurationEncoder,
	EncodeCaller:   zapcore.ShortCallerEncoder,
}

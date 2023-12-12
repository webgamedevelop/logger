package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/spf13/pflag"
	"github.com/webgamedevelop/logger"
	"go.uber.org/zap/zapcore"
	"k8s.io/klog/v2"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var filteredKlogFlags flag.FlagSet
	FilterFlags(&filteredKlogFlags)
	pflag.CommandLine.AddGoFlagSet(&filteredKlogFlags)
	pflag.Parse()

	// write log to stdout/stderr
	logrLogger, flush := logger.New(ctx, zapcore.Level(0), "text", logger.DefaultEncoderConfig, os.Stdout)
	klog.SetLoggerWithOptions(logrLogger, klog.FlushLogger(flush))
	defer klog.Flush()

	for {
		klog.Info("info log test")
		klog.V(1).Info("v1 log test")
		klog.V(2).Info("v2 log test")
		klog.V(3).Info("v3 log test")
		fmt.Println()
		time.Sleep(time.Second)
	}
}

func FilterFlags(fs *flag.FlagSet) {
	var allFlags flag.FlagSet
	klog.InitFlags(&allFlags)
	if fs == nil {
		fs = flag.CommandLine
	}
	allFlags.VisitAll(func(f *flag.Flag) {
		switch f.Name {
		case "v", "vmodule":
			fs.Var(f.Value, f.Name, f.Usage)
		}
	})
}

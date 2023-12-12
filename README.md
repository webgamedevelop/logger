# logger

logger 是一个日志配置库，通过简单的函数完成 `uber-go/zap` 的配置
* 动态调整日志级别
* 支持同时写入文件和 stdout/stderr
* 支持日志文件 rotate

## examples
### 将日志写入 stdout
```go
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

	// 启动时指定命令行参数 -v=3
	// 通过 `kill -SIGUSR1 $PID` 减少日志输出
	// 通过 `kill -SIGUSR2 $PID` 打印更详细信息
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
```

### 将日志同时写入文件和 stdout
```go
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

	// 配置日志文件
	// logger.Size(100)       最大文件大小
	// logger.Age(7)          日志文件保存时间
	// logger.Backups(7)      日志文件备份数
	// logger.LocalTime(true) 使用 LocalTime
	f := logger.NewFile("write-file.log", logger.Size(100), logger.Age(7), logger.Backups(7), logger.LocalTime(true))
	logrLogger, flush := logger.New(ctx, zapcore.Level(0), "text", logger.DefaultEncoderConfig, os.Stdout, f)
	klog.SetLoggerWithOptions(logrLogger, klog.FlushLogger(flush))
	defer klog.Flush()

	// 启动时指定命令行参数 -v=3
	// 通过 `kill -SIGUSR1 $PID` 减少日志输出
	// 通过 `kill -SIGUSR2 $PID` 打印更详细信息
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
```
![效果演示](pic/img.png)

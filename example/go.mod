module example

go 1.21.7

require (
	github.com/spf13/pflag v1.0.5
	github.com/webgamedevelop/logger v0.0.0
	k8s.io/klog/v2 v2.120.1
)

require (
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/zapr v1.3.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.2.1 // indirect
)

replace github.com/webgamedevelop/logger => ../

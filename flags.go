package logger

import (
	"flag"
)

var (
	format string
	enab   int
	stdout bool
	stderr bool
	file   File
)

var commandLine flag.FlagSet

func init() {
	commandLine.BoolVar(&stdout, "logger-stdout", true, "Log to standard out")
	commandLine.BoolVar(&stderr, "logger-stderr", false, "Log to standard error")
	commandLine.StringVar(&format, "logger-encoder", "console", "Zap log encoding (one of 'json' or 'console')")
	commandLine.IntVar(&enab, "logger-level-enabler", 0, "LevelEnabler decides whether a given logging level is enabled when logging a message")

	commandLine.StringVar(&file.Filename, "logger-filename", "", "File to write logs to")
	commandLine.IntVar(&file.MaxSize, "logger-max-size", 100, "Maximum size in megabytes of the log file before it gets rotated")
	commandLine.IntVar(&file.MaxAge, "logger-max-age", 7, "Maximum number of days to retain old log files based on the timestamp encoded in their filename")
	commandLine.IntVar(&file.MaxBackups, "logger-max-backups", 7, "Maximum number of old log files to retain")
	commandLine.BoolVar(&file.LocalTime, "logger-local-time", true, "LocalTime determines if the time used for formatting the timestamps in backup files is the computer's local time")
	commandLine.BoolVar(&file.Compress, "logger-compress", false, "Compress determines if the rotated log files should be compressed using gzip")
}

// InitFlags is for explicitly initializing the flags.
// It may get called repeatedly for different flagSets, but not
// twice for the same one. May get called concurrently
// to other goroutines using logger. However, only some flags
// may get set concurrently.
func InitFlags(flagSet *flag.FlagSet) {
	if flagSet == nil {
		flagSet = flag.CommandLine
	}

	commandLine.VisitAll(func(f *flag.Flag) {
		flagSet.Var(f.Value, f.Name, f.Usage)
	})
}

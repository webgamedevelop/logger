package logger

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

type File struct {
	lumberjack.Logger
}

func (l *File) Sync() error {
	return nil
}

// defaultFile returns a *File that implements the zapcore.WriteSyncer interface.
// filename is the file to write logs to.
// maxSize is the maximum size in megabytes of the log file before it gets rotated.
// maxAge is the maximum number of days to retain old log files based on the timestamp encoded in their filename.
// maxBackups is the maximum number of old log files to retain.
// localTime determines if the time used for formatting the timestamps in backup files is the computer's local time.
// compress determines if the rotated log files should be compressed using gzip.
func defaultFile() *File {
	return &File{lumberjack.Logger{
		Filename:   file.Filename,
		MaxSize:    file.MaxSize,
		MaxAge:     file.MaxAge,
		MaxBackups: file.MaxBackups,
		LocalTime:  file.LocalTime,
		Compress:   file.Compress,
	}}
}

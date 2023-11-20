package logger

import "gopkg.in/natefinch/lumberjack.v2"

type File struct {
	lumberjack.Logger
}

func (l *File) Sync() error {
	return l.Logger.Close()
}

// NewFile returns a *File that implements the zapcore.WriteSyncer interface.
// filename is the file to write logs to.
// maxSize is the maximum size in megabytes of the log file before it gets rotated.
// maxAge is the maximum number of days to retain old log files based on the timestamp encoded in their filename.
// maxBackups is the maximum number of old log files to retain.
// localTime determines if the time used for formatting the timestamps in backup files is the computer's local time.
// compress determines if the rotated log files should be compressed using gzip.
func NewFile(filename string, opts ...FileOption) *File {
	f := defaultFile(filename)
	for _, opt := range opts {
		opt.ApplyToFile(f)
	}
	return f
}

func defaultFile(filename string) *File {
	return &File{lumberjack.Logger{
		Filename:   filename,
		MaxSize:    100,
		MaxAge:     7,
		MaxBackups: 7,
		LocalTime:  true,
		Compress:   false,
	}}
}

type FileOption interface {
	ApplyToFile(*File)
}

type Size uint

func (o Size) ApplyToFile(f *File) {
	f.MaxSize = int(o)
}

type Age uint

func (o Age) ApplyToFile(f *File) {
	f.MaxAge = int(o)
}

type Backups uint

func (o Backups) ApplyToFile(f *File) {
	f.MaxBackups = int(o)
}

type LocalTime bool

func (o LocalTime) ApplyToFile(f *File) {
	f.LocalTime = bool(o)
}

type Compress bool

func (o Compress) ApplyToFile(f *File) {
	f.Compress = bool(o)
}

package log

import (
	"io"
	"log/slog"
	"os"
)

// KisDefaultLog 默认提供的日志对象
type KisDefaultLog struct {
	location   bool
	level      slog.Level
	jsonFormat bool
	writer     io.Writer
}

type KisLogOptions func(k *KisDefaultLog)

func WithLocation(location bool) KisLogOptions {
	return func(k *KisDefaultLog) {
		k.location = location
	}
}

func WithLevel(level slog.Level) KisLogOptions {
	return func(k *KisDefaultLog) {
		k.level = level
	}
}

func WithJSONFormat(jsonFormat bool) KisLogOptions {
	return func(k *KisDefaultLog) {
		k.jsonFormat = jsonFormat
	}
}

func WithWriter(writer io.Writer) KisLogOptions {
	return func(k *KisDefaultLog) {
		k.writer = writer
	}
}

var defaultKisLog = &KisDefaultLog{
	location:   true,
	level:      slog.LevelDebug,
	jsonFormat: false,
	writer:     os.Stdout,
}

func loadKisDefaultLog(opts ...KisLogOptions) *KisDefaultLog {
	kisLog := defaultKisLog
	if opts == nil {
		return kisLog
	}

	for _, opt := range opts {
		opt(kisLog)
	}

	return kisLog
}

// MustNewLog 初始化日志
func MustNewLog(opts ...KisLogOptions) {
	SetLogger(loadKisDefaultLog(opts...))
}

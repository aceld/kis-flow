package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

func init() {
	// 如果没有设置 Logger, 则启动时使用默认的 kisDefaultSlog 对象
	if Logger() == nil {
		MustNewKisDefaultSlog()
	}
}

func MustNewKisDefaultSlog(opts ...KisLogOptions) {
	defaultSlog := getKisDefaultSLog(opts...)
	initDefaultSlog(defaultSlog)
	SetLogger(defaultSlog)
}

// kisDefaultSlog 默认提供的日志对象
type kisDefaultSlog struct {
	location   bool
	level      slog.Level
	jsonFormat bool
	writer     io.Writer

	mu sync.Mutex
}

type KisLogOptions func(k *kisDefaultSlog)

func WithLocation(location bool) KisLogOptions {
	return func(k *kisDefaultSlog) {
		k.location = location
	}
}

func WithLevel(level slog.Level) KisLogOptions {
	return func(k *kisDefaultSlog) {
		k.level = level
	}
}

func WithJSONFormat(jsonFormat bool) KisLogOptions {
	return func(k *kisDefaultSlog) {
		k.jsonFormat = jsonFormat
	}
}

func WithWriter(writer io.Writer) KisLogOptions {
	return func(k *kisDefaultSlog) {
		k.writer = writer
	}
}

var defaultKisLog = &kisDefaultSlog{
	location:   true,
	level:      slog.LevelDebug,
	jsonFormat: false,
	writer:     os.Stdout,
}

func getKisDefaultSLog(opts ...KisLogOptions) *kisDefaultSlog {
	defaultKisSlog := defaultKisLog
	if opts == nil {
		return defaultKisSlog
	}

	for _, opt := range opts {
		opt(defaultKisSlog)
	}

	return defaultKisSlog
}

func (k *kisDefaultSlog) InfoFX(ctx context.Context, str string, v ...interface{}) {
	slog.InfoContext(ctx, fmt.Sprintf(str, v...))
}

func (k *kisDefaultSlog) ErrorFX(ctx context.Context, str string, v ...interface{}) {
	slog.ErrorContext(ctx, fmt.Sprintf(str, v...))
}

func (k *kisDefaultSlog) DebugFX(ctx context.Context, str string, v ...interface{}) {
	slog.DebugContext(ctx, fmt.Sprintf(str, v...))
}

// InfoF 使用格式化格式（xxxF或xxxFX）要使用 fmt.Sprintf() 函数进行格式化包装
func (k *kisDefaultSlog) InfoF(str string, v ...interface{}) {
	slog.Info(fmt.Sprintf(str, v...))
}

func (k *kisDefaultSlog) ErrorF(str string, v ...interface{}) {
	slog.Error(fmt.Sprintf(str, v...))
}

func (k *kisDefaultSlog) DebugF(str string, v ...interface{}) {
	slog.Debug(fmt.Sprintf(str, v...))
}

func (k *kisDefaultSlog) InfoX(ctx context.Context, str string, v ...interface{}) {
	slog.InfoContext(ctx, str, v...)
}

func (k *kisDefaultSlog) ErrorX(ctx context.Context, str string, v ...interface{}) {
	slog.ErrorContext(ctx, str, v...)
}

func (k *kisDefaultSlog) DebugX(ctx context.Context, str string, v ...interface{}) {
	slog.DebugContext(ctx, str, v...)
}

func (k *kisDefaultSlog) Info(str string, v ...interface{}) {
	slog.Info(str, v...)
}

func (k *kisDefaultSlog) Error(str string, v ...interface{}) {
	slog.Error(str, v...)
}

func (k *kisDefaultSlog) Debug(str string, v ...interface{}) {
	slog.Debug(str, v...)
}

func (k *kisDefaultSlog) SetDebugMode() {
	k.mu.Lock()
	defer k.mu.Unlock()
	k.level = slog.LevelDebug
}

func initDefaultSlog(kisLog *kisDefaultSlog) {
	replace := func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.SourceKey {
			source := a.Value.Any().(*slog.Source)
			source.Function = ""
			source.File = filepath.Base(source.File)
		}

		return a
	}

	ho := &slog.HandlerOptions{
		AddSource:   kisLog.location,
		Level:       kisLog.level,
		ReplaceAttr: replace,
	}

	var logger *slog.Logger
	if kisLog.jsonFormat {
		logger = slog.New(slog.NewJSONHandler(kisLog.writer, ho))
	} else {
		logger = slog.New(slog.NewTextHandler(kisLog.writer, ho))
	}
	slog.SetDefault(logger)
}

package log

import (
	"log/slog"
	"path/filepath"
)

// SetLogger 替换 slog
func SetLogger(kisLog *KisDefaultLog) {
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

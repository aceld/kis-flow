package log

import "context"

type KisLogger interface {
	// InfoFX with context Info-level log interface, format string format
	InfoFX(ctx context.Context, str string, v ...interface{})
	// ErrorFX with context Error-level log interface, format string format
	ErrorFX(ctx context.Context, str string, v ...interface{})
	// DebugFX with context Debug-level log interface, format string format
	DebugFX(ctx context.Context, str string, v ...interface{})

	// InfoF without context Info-level log interface, format string format
	InfoF(str string, v ...interface{})
	// ErrorF without context Error-level log interface, format string format
	ErrorF(str string, v ...interface{})
	// DebugF without context Debug-level log interface, format string format
	DebugF(str string, v ...interface{})

	// SetDebugMode set Debug mode
	SetDebugMode(enable bool)
}

// kisLog Default KisLog object, providing default log printing methods, all of which print to standard output.
var kisLog KisLogger

// SetLogger set KisLog object, can be a user-defined Logger object
func SetLogger(newlog KisLogger) {
	kisLog = newlog
}

// Logger get the kisLog object
func Logger() KisLogger {
	return kisLog
}

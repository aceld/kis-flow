package log

import "context"

type KisLogger interface {
	// InfoFX 有上下文的Info级别日志接口, format字符串格式
	InfoFX(ctx context.Context, str string, v ...interface{})
	// ErrorFX 有上下文的Error级别日志接口, format字符串格式
	ErrorFX(ctx context.Context, str string, v ...interface{})
	// DebugFX 有上下文的Debug级别日志接口, format字符串格式
	DebugFX(ctx context.Context, str string, v ...interface{})

	// InfoF 无上下文的Info级别日志接口, format字符串格式
	InfoF(str string, v ...interface{})
	// ErrorF 无上下文的Error级别日志接口, format字符串格式
	ErrorF(str string, v ...interface{})
	// DebugF 无上下文的Debug级别日志接口, format字符串格式
	DebugF(str string, v ...interface{})

	// InfoX 有上下文的Info级别日志接口, format字符串格式
	InfoX(ctx context.Context, str string, v ...interface{})
	// ErrorX 有上下文的Error级别日志接口, format字符串格式
	ErrorX(ctx context.Context, str string, v ...interface{})
	// DebugX 有上下文的Debug级别日志接口, format字符串格式
	DebugX(ctx context.Context, str string, v ...interface{})

	// Info 无上下文的Info级别日志接口, format字符串格式
	Info(str string, v ...interface{})
	// Error 无上下文的Error级别日志接口, format字符串格式
	Error(str string, v ...interface{})
	// Debug 无上下文的Debug级别日志接口, format字符串格式
	Debug(str string, v ...interface{})

	// SetDebugMode 设置Debug模式
	SetDebugMode()
}

// kisLog 默认的KisLog 对象， 提供默认的日志打印方式, 均是打印在标准输出上。
var kisLog KisLogger

// SetLogger 设置KisLog对象, 可以是用户自定义的Logger对象
func SetLogger(newlog KisLogger) {
	kisLog = newlog
}

// Logger 获取到kisLog对象
func Logger() KisLogger {
	return kisLog
}

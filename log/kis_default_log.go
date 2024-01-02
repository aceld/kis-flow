package log

import (
	"context"
	"fmt"
)

// kisDefaultLog 默认提供的日志对象
type kisDefaultLog struct{}

func (log *kisDefaultLog) InfoF(str string, v ...interface{}) {
	fmt.Printf(str, v...)
}

func (log *kisDefaultLog) ErrorF(str string, v ...interface{}) {
	fmt.Printf(str, v...)
}

func (log *kisDefaultLog) DebugF(str string, v ...interface{}) {
	fmt.Printf(str, v...)
}

func (log *kisDefaultLog) InfoFX(ctx context.Context, str string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(str, v...)
}

func (log *kisDefaultLog) ErrorFX(ctx context.Context, str string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(str, v...)
}

func (log *kisDefaultLog) DebugFX(ctx context.Context, str string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(str, v...)
}

func init() {
	// 如果没有设置Logger, 则启动时使用默认的kisDefaultLog对象
	if Logger() == nil {
		SetLogger(&kisDefaultLog{})
	}
}

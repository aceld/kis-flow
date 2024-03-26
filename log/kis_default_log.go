package log

import (
	"context"
	"fmt"
	"sync"
)

// kisDefaultLog 默认提供的日志对象
type kisDefaultLog struct {
	debugMode bool
	mu        sync.Mutex
}

func (log *kisDefaultLog) SetDebugMode(enable bool) {
	log.mu.Lock()
	defer log.mu.Unlock()
	log.debugMode = enable
}

func (log *kisDefaultLog) InfoF(str string, v ...interface{}) {
	fmt.Printf(str, v...)
	fmt.Printf("\n")
}

func (log *kisDefaultLog) ErrorF(str string, v ...interface{}) {
	fmt.Printf(str, v...)
	fmt.Printf("\n")
}

func (log *kisDefaultLog) DebugF(str string, v ...interface{}) {
	log.mu.Lock()
	defer log.mu.Unlock()
	if log.debugMode {
		fmt.Printf(str, v...)
		fmt.Printf("\n")
	}
}

func (log *kisDefaultLog) InfoFX(ctx context.Context, str string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(str, v...)
	fmt.Printf("\n")
}

func (log *kisDefaultLog) ErrorFX(ctx context.Context, str string, v ...interface{}) {
	fmt.Println(ctx)
	fmt.Printf(str, v...)
	fmt.Printf("\n")
}

func (log *kisDefaultLog) DebugFX(ctx context.Context, str string, v ...interface{}) {
	log.mu.Lock()
	defer log.mu.Unlock()
	if log.debugMode {
		fmt.Println(ctx)
		fmt.Printf(str, v...)
		fmt.Printf("\n")
	}
}

func init() {
	// 如果没有设置Logger, 则启动时使用默认的kisDefaultLog对象
	if Logger() == nil {
		SetLogger(&kisDefaultLog{})
	}
}

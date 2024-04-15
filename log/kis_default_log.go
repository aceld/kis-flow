package log

import (
	"context"
	"fmt"
	"sync"
)

// kisDefaultLog Default provided log object
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
	// If no logger is set, use the default kisDefaultLog object at startup
	if Logger() == nil {
		SetLogger(&kisDefaultLog{})
	}
}

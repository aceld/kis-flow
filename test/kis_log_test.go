package test

import (
	"context"
	"github.com/aceld/kis-flow/log"
	"testing"
)

func TestKisLogger(t *testing.T) {
	ctx := context.Background()

	log.Logger().DebugF("TestKisLogger Format DebugF name = %s, age = %d", "kisFlow", 23)
	log.Logger().ErrorF("TestKisLogger Format ErrorF name = %s, age = %d", "kisFlow", 12)
	log.Logger().InfoF("TestKisLogger Format InfoF name = %s, stu =%+v", "kisFlow",
		struct {
			name string
			age  int
		}{
			name: "kisName",
			age:  12,
		})

	log.Logger().InfoX(ctx, "TestKisLogger InfoX")
	log.Logger().ErrorX(ctx, "TestKisLogger ErrorX")
	log.Logger().DebugX(ctx, "TestKisLogger DebugX")

	log.Logger().Info("TestKisLogger Info")
	log.Logger().Error("TestKisLogger Error")
	log.Logger().Debug("TestKisLogger Debug")
}

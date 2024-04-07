package test

import (
	"context"
	"github.com/aceld/kis-flow/log"
	"testing"
)

func TestKisLogger(t *testing.T) {
	ctx := context.Background()

	log.Logger().InfoX(ctx, "TestKisLogger InfoX")
	log.Logger().ErrorX(ctx, "TestKisLogger ErrorX")
	log.Logger().DebugX(ctx, "TestKisLogger DebugX")

	log.Logger().Info("TestKisLogger Info")
	log.Logger().Error("TestKisLogger Error")
	log.Logger().Debug("TestKisLogger Debug")
}

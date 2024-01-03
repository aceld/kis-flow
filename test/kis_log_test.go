package test

import (
	"context"
	"kis-flow/log"
	"testing"
)

func TestKisLogger(t *testing.T) {
	ctx := context.Background()

	log.Logger().InfoFX(ctx, "TestKisLogger InfoFX")
	log.Logger().ErrorFX(ctx, "TestKisLogger ErrorFX")
	log.Logger().DebugFX(ctx, "TestKisLogger DebugFX")

	log.Logger().InfoF("TestKisLogger InfoF")
	log.Logger().ErrorF("TestKisLogger ErrorF")
	log.Logger().DebugF("TestKisLogger DebugF")
}

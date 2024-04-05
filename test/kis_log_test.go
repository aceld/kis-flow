package test

import (
	"context"
	"github.com/aceld/kis-flow/log"
	"log/slog"
	"testing"
)

func TestKisLogger(t *testing.T) {
	log.MustNewLog(log.WithJSONFormat(true))
	ctx := context.Background()

	slog.InfoContext(ctx, "TestKisLogger InfoFX")
	slog.ErrorContext(ctx, "TestKisLogger ErrorFX")
	slog.DebugContext(ctx, "TestKisLogger DebugFX")

	slog.Info("TestKisLogger InfoF")
	slog.Error("TestKisLogger ErrorF")
	slog.Debug("TestKisLogger DebugF")
}

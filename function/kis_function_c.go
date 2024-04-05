package function

import (
	"context"
	"log/slog"

	"github.com/aceld/kis-flow/kis"
)

type KisFunctionC struct {
	BaseFunction
}

func NewKisFunctionC() kis.Function {
	f := new(KisFunctionC)

	// 初始化metaData
	f.metaData = make(map[string]interface{})

	return f
}

func (f *KisFunctionC) Call(ctx context.Context, flow kis.Flow) error {
	slog.Debug("KisFunctionC", "flow", flow)

	// 通过KisPool 路由到具体的执行计算Function中
	if err := kis.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		slog.ErrorContext(ctx, "Function Called Error", "err", err)
		return err
	}

	return nil
}

package function

import (
	"context"
	"log/slog"

	"github.com/aceld/kis-flow/kis"
)

type KisFunctionL struct {
	BaseFunction
}

func NewKisFunctionL() kis.Function {
	f := new(KisFunctionL)

	// 初始化metaData
	f.metaData = make(map[string]interface{})

	return f
}

func (f *KisFunctionL) Call(ctx context.Context, flow kis.Flow) error {
	slog.Debug("KisFunctionL", "flow", flow)

	// 通过KisPool 路由到具体的执行计算Function中
	if err := kis.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		slog.ErrorContext(ctx, "Function Called Error", "err", err)
		return err
	}

	return nil
}

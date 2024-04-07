package function

import (
	"context"
	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/log"
)

type KisFunctionS struct {
	BaseFunction
}

func NewKisFunctionS() kis.Function {
	f := new(KisFunctionS)

	// 初始化metaData
	f.metaData = make(map[string]interface{})

	return f
}

func (f *KisFunctionS) Call(ctx context.Context, flow kis.Flow) error {
	log.Logger().Debug("KisFunctionS", "flow", flow)

	// 通过KisPool 路由到具体的执行计算Function中
	if err := kis.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		log.Logger().ErrorX(ctx, "Function Called Error", "err", err)
		return err
	}

	return nil
}

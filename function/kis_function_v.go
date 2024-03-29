package function

import (
	"context"
	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/log"
)

type KisFunctionV struct {
	BaseFunction
}

func NewKisFunctionV() kis.Function {
	f := new(KisFunctionV)

	// 初始化metaData
	f.metaData = make(map[string]interface{})

	return f
}

func (f *KisFunctionV) Call(ctx context.Context, flow kis.Flow) error {
	log.Logger().DebugF("KisFunctionV, flow = %+v\n", flow)

	// 通过KisPool 路由到具体的执行计算Function中
	if err := kis.Pool().CallFunction(ctx, f.Config.FName, flow); err != nil {
		log.Logger().ErrorFX(ctx, "Function Called Error err = %s\n", err)
		return err
	}

	return nil
}

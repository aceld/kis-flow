package function

import (
	"context"
	"kis-flow/kis"
	"kis-flow/log"
)

type KisFunctionL struct {
	BaseFunction
}

func (f *KisFunctionL) Call(ctx context.Context, flow kis.Flow) error {
	log.Logger().InfoF("KisFunctionL, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法

	return nil
}

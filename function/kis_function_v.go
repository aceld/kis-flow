package function

import (
	"context"
	"kis-flow/kis"
	"kis-flow/log"
)

type KisFunctionV struct {
	BaseFunction
}

func (f *KisFunctionV) Call(ctx context.Context, flow kis.Flow) error {
	log.Logger().InfoF("KisFunctionV, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法

	return nil
}

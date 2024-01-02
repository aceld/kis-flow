package function

import (
	"context"
	"kis-flow/kis"
	"kis-flow/log"
)

type KisFunctionE struct {
	BaseFunction
}

func (f *KisFunctionE) Call(ctx context.Context, flow kis.Flow) error {
	log.Logger().InfoF("KisFunctionE, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法

	return nil
}

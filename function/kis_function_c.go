package function

import (
	"context"
	"kis-flow/kis"
	"kis-flow/log"
)

type KisFunctionC struct {
	BaseFunction
}

func (f *KisFunctionC) Call(ctx context.Context, flow kis.Flow) error {
	log.Logger().InfoF("KisFunctionC, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法

	return nil
}

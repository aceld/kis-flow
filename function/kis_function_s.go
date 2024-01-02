package function

import (
	"context"
	"kis-flow/kis"
	"kis-flow/log"
)

type KisFunctionS struct {
	BaseFunction
}

func (f *KisFunctionS) Call(ctx context.Context, flow kis.Flow) error {
	log.Logger().InfoF("KisFunctionS, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法

	return nil
}

package function

import (
	"context"
	"fmt"
	"kis-flow/kis"
)

type KisFunctionE struct {
	BaseFunction
}

func (f *KisFunctionE) Call(ctx context.Context, flow kis.Flow) error {
	fmt.Printf("KisFunctionE, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法

	return nil
}

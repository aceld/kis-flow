package function

import (
	"context"
	"fmt"
	"kis-flow/kis"
	"kis-flow/log"
)

type KisFunctionE struct {
	BaseFunction
}

func (f *KisFunctionE) Call(ctx context.Context, flow kis.Flow) error {
	log.Logger().InfoF("KisFunctionE, flow = %+v\n", flow)

	// TODO 调用具体的Function执行方法
	//处理业务数据
	for _, row := range flow.Input() {
		fmt.Printf("In KisFunctionE, row = %+v\n", row)
	}

	return nil
}

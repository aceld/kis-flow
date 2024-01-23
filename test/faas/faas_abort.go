package faas

import (
	"context"
	"fmt"
	"kis-flow/kis"
)

// type FaaS func(context.Context, Flow) error

func AbortFuncHandler(ctx context.Context, flow kis.Flow) error {
	fmt.Println("---> Call AbortFuncHandler ----")

	for _, row := range flow.Input() {
		str := fmt.Sprintf("In FuncName = %s, FuncId = %s, row = %s", flow.GetThisFuncConf().FName, flow.GetThisFunction().GetId(), row)
		fmt.Println(str)
	}

	return flow.Next(kis.ActionAbort)
}

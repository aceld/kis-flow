package faas

import (
	"context"
	"fmt"

	"github.com/aceld/kis-flow/kis"
)

// type FaaS func(context.Context, Flow) error

func FuncDemo4Handler(ctx context.Context, flow kis.Flow) error {
	fmt.Println("---> Call FuncDemo4Handler ----")

	for index, row := range flow.Input() {
		str := fmt.Sprintf("In FuncName = %s, FuncId = %s, row = %s", flow.GetThisFuncConf().FName, flow.GetThisFunction().GetID(), row)
		fmt.Println(str)

		resultStr := fmt.Sprintf("data from funcName[%s], index = %d", flow.GetThisFuncConf().FName, index)

		_ = flow.CommitRow(resultStr)
	}

	return nil
}

package faas

import (
	"context"
	"fmt"
	"kis-flow/kis"
	"kis-flow/log"
)

// type FaaS func(context.Context, Flow) error

func FuncDemo2Handler(ctx context.Context, flow kis.Flow) error {
	fmt.Println("---> Call funcName2Handler ----")

	for index, row := range flow.Input() {
		str := fmt.Sprintf("In FuncName = %s, FuncId = %s, row = %s", flow.GetThisFuncConf().FName, flow.GetThisFunction().GetId(), row)
		fmt.Println(str)

		conn, err := flow.GetConnector()
		if err != nil {
			log.Logger().ErrorFX(ctx, "FuncDemo2Handler(): GetConnector err = %s\n", err.Error())
			return err
		}

		if conn.Call(ctx, flow, row) != nil {
			log.Logger().ErrorFX(ctx, "FuncDemo2Handler(): Call err = %s\n", err.Error())
			return err
		}

		// 计算结果数据
		resultStr := fmt.Sprintf("data from funcName[%s], index = %d", flow.GetThisFuncConf().FName, index)

		// 提交结果数据
		_ = flow.CommitRow(resultStr)
	}

	return nil
}

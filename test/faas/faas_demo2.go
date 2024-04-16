package faas

import (
	"context"
	"fmt"

	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/log"
)

// type FaaS func(context.Context, Flow) error

func FuncDemo2Handler(ctx context.Context, flow kis.Flow) error {
	fmt.Println("---> Call funcName2Handler ----")
	fmt.Printf("Params = %+v\n", flow.GetFuncParamAll())

	for index, row := range flow.Input() {
		str := fmt.Sprintf("In FuncName = %s, FuncId = %s, row = %s", flow.GetThisFuncConf().FName, flow.GetThisFunction().GetID(), row)
		fmt.Println(str)

		conn, err := flow.GetConnector()
		if err != nil {
			log.Logger().ErrorFX(ctx, "FuncDemo2Handler(): GetConnector err = %s\n", err.Error())
			return err
		}

		if _, err := conn.Call(ctx, flow, row); err != nil {
			log.Logger().ErrorFX(ctx, "FuncDemo2Handler(): Call err = %s\n", err.Error())
			return err
		}

		resultStr := fmt.Sprintf("data from funcName[%s], index = %d", flow.GetThisFuncConf().FName, index)

		_ = flow.CommitRow(resultStr)
	}

	return nil
}

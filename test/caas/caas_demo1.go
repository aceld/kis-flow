package caas

import (
	"context"
	"fmt"
	"kis-flow/kis"
)

// type CaaS func(context.Context, Connector, Function, Flow, interface{}) error

func CaasDemoHanler1(ctx context.Context, conn kis.Connector, fn kis.Function, flow kis.Flow, args interface{}) error {
	fmt.Printf("===> In CaasDemoHanler1: flowName: %s, cName:%s, fnName:%s, mode:%s\n",
		flow.GetName(), conn.GetName(), fn.GetConfig().FName, fn.GetConfig().FMode)

	fmt.Printf("===> Call Connector CaasDemoHanler1, args from funciton: %s\n", args)

	return nil
}

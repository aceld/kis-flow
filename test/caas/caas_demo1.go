package caas

import (
	"context"
	"fmt"
	"github.com/aceld/kis-flow/kis"
)

// type CaaS func(context.Context, Connector, Function, Flow, interface{}) error

func CaasDemoHanler1(ctx context.Context, conn kis.Connector, fn kis.Function, flow kis.Flow, args interface{}) (interface{}, error) {
	fmt.Printf("===> In CaasDemoHanler1: flowName: %s, cName:%s, fnName:%s, mode:%s\n",
		flow.GetName(), conn.GetName(), fn.GetConfig().FName, fn.GetConfig().FMode)

	fmt.Printf("Params = %+v\n", conn.GetConfig().Params)

	fmt.Printf("===> Call Connector CaasDemoHanler1, args from funciton: %s\n", args)

	return nil, nil
}

package function

import (
	"context"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/flow"
	"testing"
)

func TestNewKisFunction(t *testing.T) {
	ctx := context.Background()

	// 1. 创建一个KisFunction配置实例
	source := config.KisSource{
		Name: "公众号抖音商城户订单数据",
		Must: []string{"order_id", "user_id"},
	}

	myFuncConfig1 := config.NewFuncConfig("funcId1", "funcName", common.C, &source, nil)
	if myFuncConfig1 == nil {
		panic("myFuncConfig1 is nil")
	}

	// 2. 创建一个 KisFlow 配置实例
	flowFuncParams1 := config.KisFlowFunctionParam{
		Fid: "funcId1",
	}

	myFlowConfig1 := config.NewFlowConfig("flowId", "flowName", common.FlowEnable)
	myFlowConfig1.AppendFunctionConfig(flowFuncParams1)

	// 3. 创建一个KisFlow对象
	flow1 := flow.NewKisFlow(myFlowConfig1)

	// 4. 创建一个KisFunction对象
	func1 := NewKisFunction(flow1, myFuncConfig1)

	if err := func1.Call(ctx, flow1); err != nil {
		t.Errorf("func1.Call() error = %v", err)
	}
}

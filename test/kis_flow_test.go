package test

import (
	"context"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/flow"
	"testing"
)

func TestNewKisFlow(t *testing.T) {
	ctx := context.Background()

	// 1. 创建2个KisFunction配置实例
	source1 := config.KisSource{
		Name: "公众号抖音商城户订单数据",
		Must: []string{"order_id", "user_id"},
	}

	source2 := config.KisSource{
		Name: "用户订单错误率",
		Must: []string{"order_id", "user_id"},
	}

	myFuncConfig1 := config.NewFuncConfig("funcName1", common.C, &source1, nil)
	if myFuncConfig1 == nil {
		panic("myFuncConfig1 is nil")
	}

	myFuncConfig2 := config.NewFuncConfig("funcName2", common.V, &source2, nil)
	if myFuncConfig2 == nil {
		panic("myFuncConfig2 is nil")
	}

	// 2. 创建一个 KisFlow 配置实例
	myFlowConfig1 := config.NewFlowConfig("flowName1", common.FlowEnable)

	// 3. 创建一个KisFlow对象
	flow1 := flow.NewKisFlow(myFlowConfig1)

	// 4. 拼接Functioin 到 Flow 上
	if err := flow1.Link(myFuncConfig1, nil); err != nil {
		panic(err)
	}
	if err := flow1.Link(myFuncConfig2, nil); err != nil {
		panic(err)
	}

	// 5. 执行flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

func TestNewKisFlowData(t *testing.T) {
	ctx := context.Background()

	// 1. 创建2个KisFunction配置实例
	source1 := config.KisSource{
		Name: "公众号抖音商城户订单数据",
		Must: []string{"order_id", "user_id"},
	}

	source2 := config.KisSource{
		Name: "用户订单错误率",
		Must: []string{"order_id", "user_id"},
	}

	myFuncConfig1 := config.NewFuncConfig("funcName1", common.C, &source1, nil)
	if myFuncConfig1 == nil {
		panic("myFuncConfig1 is nil")
	}

	myFuncConfig2 := config.NewFuncConfig("funcName2", common.E, &source2, nil)
	if myFuncConfig2 == nil {
		panic("myFuncConfig2 is nil")
	}

	// 2. 创建一个 KisFlow 配置实例
	myFlowConfig1 := config.NewFlowConfig("flowName1", common.FlowEnable)

	// 3. 创建一个KisFlow对象
	flow1 := flow.NewKisFlow(myFlowConfig1)

	// 4. 拼接Functioin 到 Flow 上
	if err := flow1.Link(myFuncConfig1, nil); err != nil {
		panic(err)
	}
	if err := flow1.Link(myFuncConfig2, nil); err != nil {
		panic(err)
	}

	// 5. 提交原始数据
	_ = flow1.CommitRow("This is Data1 from Test")
	_ = flow1.CommitRow("This is Data2 from Test")
	_ = flow1.CommitRow("This is Data3 from Test")

	// 6. 执行flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

package test

import (
	"context"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/flow"
	"kis-flow/kis"
	"kis-flow/test/caas"
	"kis-flow/test/faas"
	"testing"
)

func TestNewKisConnector(t *testing.T) {

	ctx := context.Background()

	// 0. 注册Function 回调业务
	kis.Pool().FaaS("funcName1", faas.FuncDemo1Handler)
	kis.Pool().FaaS("funcName2", faas.FuncDemo2Handler)
	kis.Pool().FaaS("funcName3", faas.FuncDemo3Handler)

	// 0. 注册ConnectorInit 和 Connector 回调业务
	kis.Pool().CaaSInit("ConnName1", caas.InitConnDemo1)
	kis.Pool().CaaS("ConnName1", "funcName2", common.S, caas.CaasDemoHanler1)

	// 1. 创建3个KisFunction配置实例, 其中myFuncConfig2 有Connector配置
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

	option := config.KisFuncOption{
		CName: "ConnName1",
	}

	myFuncConfig2 := config.NewFuncConfig("funcName2", common.S, &source2, &option)
	if myFuncConfig2 == nil {
		panic("myFuncConfig2 is nil")
	}

	myFuncConfig3 := config.NewFuncConfig("funcName3", common.E, &source2, nil)
	if myFuncConfig3 == nil {
		panic("myFuncConfig3 is nil")
	}

	// 2. 创建一个KisConnector配置实例
	myConnConfig1 := config.NewConnConfig("ConnName1", "0.0.0.0:9998", common.REDIS, "redis-key", nil)
	if myConnConfig1 == nil {
		panic("myConnConfig1 is nil")
	}

	// 3. 将KisConnector配置实例绑定到KisFunction配置实例上
	_ = myFuncConfig2.AddConnConfig(myConnConfig1)

	// 4. 创建一个 KisFlow 配置实例
	myFlowConfig1 := config.NewFlowConfig("flowName1", common.FlowEnable)

	// 5. 创建一个KisFlow对象
	flow1 := flow.NewKisFlow(myFlowConfig1)

	// 6. 拼接Functioin 到 Flow 上
	if err := flow1.Link(myFuncConfig1, nil); err != nil {
		panic(err)
	}
	if err := flow1.Link(myFuncConfig2, nil); err != nil {
		panic(err)
	}
	if err := flow1.Link(myFuncConfig3, nil); err != nil {
		panic(err)
	}

	// 7. 提交原始数据
	_ = flow1.CommitRow("This is Data1 from Test")
	_ = flow1.CommitRow("This is Data2 from Test")
	_ = flow1.CommitRow("This is Data3 from Test")

	// 8. 执行flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

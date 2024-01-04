package test

import (
	"context"
	"fmt"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/flow"
	"kis-flow/kis"
	"testing"
)

func funcName1Handler(ctx context.Context, flow kis.Flow) error {
	fmt.Println("---> Call funcName1Handler ----")

	for index, row := range flow.Input() {
		// 打印数据
		str := fmt.Sprintf("In FuncName = %s, FuncId = %s, row = %s", flow.GetThisFuncConf().FName, flow.GetThisFunction().GetId(), row)
		fmt.Println(str)

		// 计算结果数据
		resultStr := fmt.Sprintf("data from funcName[%s], index = %d", flow.GetThisFuncConf().FName, index)

		// 提交结果数据
		_ = flow.CommitRow(resultStr)
	}

	return nil
}

func funcName2Handler(ctx context.Context, flow kis.Flow) error {

	for _, row := range flow.Input() {
		str := fmt.Sprintf("In FuncName = %s, FuncId = %s, row = %s", flow.GetThisFuncConf().FName, flow.GetThisFunction().GetId(), row)
		fmt.Println(str)
	}

	return nil
}

func TestNewKisPool(t *testing.T) {

	ctx := context.Background()

	// 0. 注册Function
	kis.Pool().FaaS("funcName1", funcName1Handler)
	kis.Pool().FaaS("funcName2", funcName2Handler)

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

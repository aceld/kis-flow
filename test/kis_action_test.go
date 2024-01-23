package test

import (
	"context"
	"kis-flow/common"
	"kis-flow/file"
	"kis-flow/kis"
	"kis-flow/test/caas"
	"kis-flow/test/faas"
	"testing"
)

func TestActionAbort(t *testing.T) {
	ctx := context.Background()

	// 0. 注册Function 回调业务
	kis.Pool().FaaS("funcName1", faas.FuncDemo1Handler)
	kis.Pool().FaaS("abortFunc", faas.AbortFuncHandler) // 添加abortFunc 业务
	kis.Pool().FaaS("funcName3", faas.FuncDemo3Handler)

	// 0. 注册ConnectorInit 和 Connector 回调业务
	kis.Pool().CaaSInit("ConnName1", caas.InitConnDemo1)
	kis.Pool().CaaS("ConnName1", "funcName2", common.S, caas.CaasDemoHanler1)

	// 1. 加载配置文件并构建Flow
	if err := file.ConfigImportYaml("/Users/tal/gopath/src/kis-flow/test/load_conf/"); err != nil {
		panic(err)
	}

	// 2. 获取Flow
	flow1 := kis.Pool().GetFlow("flowName2")

	// 3. 提交原始数据
	_ = flow1.CommitRow("This is Data1 from Test")
	_ = flow1.CommitRow("This is Data2 from Test")
	_ = flow1.CommitRow("This is Data3 from Test")

	// 4. 执行flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

func TestActionDataReuse(t *testing.T) {
	ctx := context.Background()

	// 0. 注册Function 回调业务
	kis.Pool().FaaS("funcName1", faas.FuncDemo1Handler)
	kis.Pool().FaaS("dataReuseFunc", faas.DataReuseFuncHandler) // 添加dataReuesFunc 业务
	kis.Pool().FaaS("funcName3", faas.FuncDemo3Handler)

	// 0. 注册ConnectorInit 和 Connector 回调业务
	kis.Pool().CaaSInit("ConnName1", caas.InitConnDemo1)
	kis.Pool().CaaS("ConnName1", "funcName2", common.S, caas.CaasDemoHanler1)

	// 1. 加载配置文件并构建Flow
	if err := file.ConfigImportYaml("/Users/tal/gopath/src/kis-flow/test/load_conf/"); err != nil {
		panic(err)
	}

	// 2. 获取Flow
	flow1 := kis.Pool().GetFlow("flowName3")

	// 3. 提交原始数据
	_ = flow1.CommitRow("This is Data1 from Test")
	_ = flow1.CommitRow("This is Data2 from Test")
	_ = flow1.CommitRow("This is Data3 from Test")

	// 4. 执行flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

func TestActionForceEntry(t *testing.T) {
	ctx := context.Background()

	// 0. 注册Function 回调业务
	kis.Pool().FaaS("funcName1", faas.FuncDemo1Handler)
	kis.Pool().FaaS("noResultFunc", faas.NoResultFuncHandler) // 添加noResultFunc 业务
	kis.Pool().FaaS("funcName3", faas.FuncDemo3Handler)

	// 0. 注册ConnectorInit 和 Connector 回调业务
	kis.Pool().CaaSInit("ConnName1", caas.InitConnDemo1)
	kis.Pool().CaaS("ConnName1", "funcName2", common.S, caas.CaasDemoHanler1)

	// 1. 加载配置文件并构建Flow
	if err := file.ConfigImportYaml("/Users/tal/gopath/src/kis-flow/test/load_conf/"); err != nil {
		panic(err)
	}

	// 2. 获取Flow
	flow1 := kis.Pool().GetFlow("flowName4")

	// 3. 提交原始数据
	_ = flow1.CommitRow("This is Data1 from Test")
	_ = flow1.CommitRow("This is Data2 from Test")
	_ = flow1.CommitRow("This is Data3 from Test")

	// 4. 执行flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

func TestActionJumpFunc(t *testing.T) {
	ctx := context.Background()

	// 0. 注册Function 回调业务
	kis.Pool().FaaS("funcName1", faas.FuncDemo1Handler)
	kis.Pool().FaaS("funcName2", faas.FuncDemo2Handler)
	kis.Pool().FaaS("jumpFunc", faas.JumpFuncHandler) // 添加jumpFunc 业务

	// 0. 注册ConnectorInit 和 Connector 回调业务
	kis.Pool().CaaSInit("ConnName1", caas.InitConnDemo1)
	kis.Pool().CaaS("ConnName1", "funcName2", common.S, caas.CaasDemoHanler1)

	// 1. 加载配置文件并构建Flow
	if err := file.ConfigImportYaml("/Users/tal/gopath/src/kis-flow/test/load_conf/"); err != nil {
		panic(err)
	}

	// 2. 获取Flow
	flow1 := kis.Pool().GetFlow("flowName5")

	// 3. 提交原始数据
	_ = flow1.CommitRow("This is Data1 from Test")
	_ = flow1.CommitRow("This is Data2 from Test")
	_ = flow1.CommitRow("This is Data3 from Test")

	// 4. 执行flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

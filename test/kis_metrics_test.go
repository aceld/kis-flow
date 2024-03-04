package test

import (
	"context"
	"kis-flow/common"
	"kis-flow/file"
	"kis-flow/kis"
	"kis-flow/test/caas"
	"kis-flow/test/faas"
	"testing"
	"time"
)

func TestMetricsDataTotal(t *testing.T) {
	ctx := context.Background()

	// 0. 注册Function 回调业务
	kis.Pool().FaaS("funcName1", faas.FuncDemo1Handler)
	kis.Pool().FaaS("funcName2", faas.FuncDemo2Handler)
	kis.Pool().FaaS("funcName3", faas.FuncDemo3Handler)

	// 0. 注册ConnectorInit 和 Connector 回调业务
	kis.Pool().CaaSInit("ConnName1", caas.InitConnDemo1)
	kis.Pool().CaaS("ConnName1", "funcName2", common.S, caas.CaasDemoHanler1)

	// 1. 加载配置文件并构建Flow
	if err := file.ConfigImportYaml("/Users/tal/gopath/src/kis-flow/test/load_conf/"); err != nil {
		panic(err)
	}

	// 2. 获取Flow
	flow1 := kis.Pool().GetFlow("flowName1")

	n := 0

	for n < 10 {
		// 3. 提交原始数据
		_ = flow1.CommitRow("This is Data1 from Test")

		// 4. 执行flow1
		if err := flow1.Run(ctx); err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)
		n++
	}

	select {}
}

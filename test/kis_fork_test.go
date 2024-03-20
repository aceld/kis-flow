package test

import (
	"context"
	"kis-flow/file"
	"kis-flow/kis"
	"testing"
)

func TestForkFlow(t *testing.T) {
	ctx := context.Background()

	// 1. 加载配置文件并构建Flow
	if err := file.ConfigImportYaml("load_conf/"); err != nil {
		panic(err)
	}

	// 2. 获取Flow
	flow1 := kis.Pool().GetFlow("flowName1")

	flow1Clone1 := flow1.Fork(ctx)

	// 3. 提交原始数据
	_ = flow1Clone1.CommitRow("This is Data1 from Test")
	_ = flow1Clone1.CommitRow("This is Data2 from Test")
	_ = flow1Clone1.CommitRow("This is Data3 from Test")

	// 4. 执行flow1
	if err := flow1Clone1.Run(ctx); err != nil {
		panic(err)
	}
}

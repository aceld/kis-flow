package test

import (
	"context"
	"github.com/aceld/kis-flow/file"
	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/log"
	"testing"
)

func TestForkFlowCommitBatch(t *testing.T) {
	ctx := context.Background()

	// 1. 加载配置文件并构建Flow
	if err := file.ConfigImportYaml("load_conf/"); err != nil {
		panic(err)
	}

	// 2. 获取Flow
	flow1 := kis.Pool().GetFlow("flowName1")

	stringRows := []string{
		"This is Data1 from Test",
		"This is Data2 from Test",
		"This is Data3 from Test",
	}

	// 3. 提交原始数据
	if err := flow1.CommitRowBatch(stringRows); err != nil {
		log.Logger().Error("CommitRowBatch Error", "err", err)
		panic(err)
	}

	// 4. 执行flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

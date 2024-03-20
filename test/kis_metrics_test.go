package test

import (
	"context"
	"fmt"
	"kis-flow/file"
	"kis-flow/kis"
	"testing"
	"time"
)

func TestMetricsDataTotal(t *testing.T) {
	ctx := context.Background()

	// 1. 加载配置文件并构建Flow
	if err := file.ConfigImportYaml("/Users/Aceld/go/src/kis-flow/test/load_conf/"); err != nil {
		fmt.Println("Wrong Config Yaml Path!")
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

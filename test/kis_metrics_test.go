package test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aceld/kis-flow/file"
	"github.com/aceld/kis-flow/kis"
)

func TestMetricsDataTotal(t *testing.T) {
	ctx := context.Background()

	if err := file.ConfigImportYaml("load_conf/"); err != nil {
		fmt.Println("Wrong Config Yaml Path!")
		panic(err)
	}

	flow1 := kis.Pool().GetFlow("flowName1")

	n := 0

	for n < 10 {
		_ = flow1.CommitRow("This is Data1 from Test")

		if err := flow1.Run(ctx); err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)
		n++
	}

	select {}
}

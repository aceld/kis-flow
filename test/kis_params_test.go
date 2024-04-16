package test

import (
	"context"
	"testing"

	"github.com/aceld/kis-flow/file"
	"github.com/aceld/kis-flow/kis"
)

func TestParams(t *testing.T) {
	ctx := context.Background()

	if err := file.ConfigImportYaml("load_conf/"); err != nil {
		panic(err)
	}

	flow1 := kis.Pool().GetFlow("flowName1")

	_ = flow1.CommitRow("This is Data1 from Test")
	_ = flow1.CommitRow("This is Data2 from Test")
	_ = flow1.CommitRow("This is Data3 from Test")

	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

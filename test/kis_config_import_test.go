package test

import (
	"context"
	"testing"

	"github.com/aceld/kis-flow/file"
	"github.com/aceld/kis-flow/kis"
)

func TestConfigImportYaml(t *testing.T) {
	ctx := context.Background()

	// 1. Load the configuration file and build the Flow
	if err := file.ConfigImportYaml("load_conf/"); err != nil {
		panic(err)
	}

	// 2. Get the Flow
	flow1 := kis.Pool().GetFlow("flowName1")

	// 3. Commit the raw data
	_ = flow1.CommitRow("This is Data1 from Test")
	_ = flow1.CommitRow("This is Data2 from Test")
	_ = flow1.CommitRow("This is Data3 from Test")

	// 4. Execute flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

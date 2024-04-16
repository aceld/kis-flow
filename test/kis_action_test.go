package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aceld/kis-flow/file"
	"github.com/aceld/kis-flow/kis"
)

func TestActionAbort(t *testing.T) {
	ctx := context.Background()

	// 1. Load the configuration file and build the Flow
	if err := file.ConfigImportYaml("load_conf/"); err != nil {
		fmt.Println("Wrong Config Yaml Path!")
		panic(err)
	}

	// 2. Get the Flow
	flow1 := kis.Pool().GetFlow("flowName2")

	// 3. Commit original data
	_ = flow1.CommitRow("This is Data1 from Test")
	_ = flow1.CommitRow("This is Data2 from Test")
	_ = flow1.CommitRow("This is Data3 from Test")

	// 4. Execute flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

func TestActionDataReuse(t *testing.T) {
	ctx := context.Background()

	// 1. Load the configuration file and build the Flow
	if err := file.ConfigImportYaml("load_conf/"); err != nil {
		fmt.Println("Wrong Config Yaml Path!")
		panic(err)
	}

	// 2. Get the Flow
	flow1 := kis.Pool().GetFlow("flowName3")

	// 3. Commit original data
	_ = flow1.CommitRow("This is Data1 from Test")
	_ = flow1.CommitRow("This is Data2 from Test")
	_ = flow1.CommitRow("This is Data3 from Test")

	// 4. Execute flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

func TestActionForceEntry(t *testing.T) {
	ctx := context.Background()

	// 1. Load the configuration file and build the Flow
	if err := file.ConfigImportYaml("load_conf/"); err != nil {
		fmt.Println("Wrong Config Yaml Path!")
		panic(err)
	}

	// 2. Get the Flow
	flow1 := kis.Pool().GetFlow("flowName4")

	// 3. Commit original data
	_ = flow1.CommitRow("This is Data1 from Test")
	_ = flow1.CommitRow("This is Data2 from Test")
	_ = flow1.CommitRow("This is Data3 from Test")

	// 4. Execute flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

func TestActionJumpFunc(t *testing.T) {
	ctx := context.Background()

	// 1. Load the configuration file and build the Flow
	if err := file.ConfigImportYaml("load_conf/"); err != nil {
		fmt.Println("Wrong Config Yaml Path!")
		panic(err)
	}

	// 2. Get the Flow
	flow1 := kis.Pool().GetFlow("flowName5")

	// 3. Commit original data
	_ = flow1.CommitRow("This is Data1 from Test")
	_ = flow1.CommitRow("This is Data2 from Test")
	_ = flow1.CommitRow("This is Data3 from Test")

	// 4. Execute flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

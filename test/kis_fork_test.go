package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/config"
	"github.com/aceld/kis-flow/file"
	"github.com/aceld/kis-flow/flow"
	"github.com/aceld/kis-flow/kis"
)

func TestForkFlow(t *testing.T) {
	ctx := context.Background()

	// 1. Load configuration file and build Flow
	if err := file.ConfigImportYaml("load_conf/"); err != nil {
		panic(err)
	}

	// 2. Get Flow
	flow1 := kis.Pool().GetFlow("flowFork1")

	fmt.Println("----> flow1: ", flow1.GetFuncParamsAllFuncs())

	flow1Clone1 := flow1.Fork(ctx)
	fmt.Println("----> flow1Clone1: ", flow1Clone1.GetFuncParamsAllFuncs())

	// 3. Commit raw data
	_ = flow1Clone1.CommitRow("This is Data1 from Test")

	// 4. Execute flow1
	if err := flow1Clone1.Run(ctx); err != nil {
		panic(err)
	}
}

func TestForkFlowWithLink(t *testing.T) {
	ctx := context.Background()

	// Create a new flow configuration
	myFlowConfig1 := config.NewFlowConfig("flowFork1", common.FlowEnable)

	// Create new function configuration
	func1Config := config.NewFuncConfig("funcName1", common.V, nil, nil)
	func3Config := config.NewFuncConfig("funcName3", common.E, nil, nil)

	// Create a new flow
	flow1 := flow.NewKisFlow(myFlowConfig1)

	_ = flow1.Link(func1Config, config.FParam{"school": "TsingHua1", "country": "China1"})
	_ = flow1.Link(func3Config, config.FParam{"school": "TsingHua3", "country": "China3"})

	fmt.Println("----> flow1: ", flow1.GetFuncParamsAllFuncs())

	flow1Clone1 := flow1.Fork(ctx)

	fmt.Println("----> flow1Clone1: ", flow1Clone1.GetFuncParamsAllFuncs())

	// 3. Commit raw data
	_ = flow1Clone1.CommitRow("This is Data1 from Test")

	// 4. Execute flow1
	if err := flow1Clone1.Run(ctx); err != nil {
		panic(err)
	}
}

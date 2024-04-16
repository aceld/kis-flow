package test

import (
	"context"
	"testing"

	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/config"
	"github.com/aceld/kis-flow/flow"
	"github.com/aceld/kis-flow/function"
)

func TestNewKisFunction(t *testing.T) {
	ctx := context.Background()

	// 1. Create a KisFunction configuration instance
	source := config.KisSource{
		Name: "TikTokOrder",
		Must: []string{"order_id", "user_id"},
	}

	myFuncConfig1 := config.NewFuncConfig("funcName1", common.C, &source, nil)
	if myFuncConfig1 == nil {
		panic("myFuncConfig1 is nil")
	}

	// 2. Create a KisFlow configuration instance
	myFlowConfig1 := config.NewFlowConfig("flowName1", common.FlowEnable)

	// 3. Create a KisFlow object
	flow1 := flow.NewKisFlow(myFlowConfig1)

	// 4. Create a KisFunction object
	func1 := function.NewKisFunction(flow1, myFuncConfig1)

	if err := func1.Call(ctx, flow1); err != nil {
		t.Errorf("func1.Call() error = %v", err)
	}
}

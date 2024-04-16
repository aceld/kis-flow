package test

import (
	"context"
	"testing"

	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/config"
	"github.com/aceld/kis-flow/flow"
)

func TestNewKisConnector(t *testing.T) {

	ctx := context.Background()

	// 1. Create three KisFunction configuration instances, with myFuncConfig2 having a Connector configuration
	source1 := config.KisSource{
		Name: "TikTokOrder",
		Must: []string{"order_id", "user_id"},
	}

	source2 := config.KisSource{
		Name: "UserOrderErrorRate",
		Must: []string{"order_id", "user_id"},
	}

	myFuncConfig1 := config.NewFuncConfig("funcName1", common.C, &source1, nil)
	if myFuncConfig1 == nil {
		panic("myFuncConfig1 is nil")
	}

	option := config.KisFuncOption{
		CName: "ConnName1",
	}

	myFuncConfig2 := config.NewFuncConfig("funcName2", common.S, &source2, &option)
	if myFuncConfig2 == nil {
		panic("myFuncConfig2 is nil")
	}

	myFuncConfig3 := config.NewFuncConfig("funcName3", common.E, &source2, nil)
	if myFuncConfig3 == nil {
		panic("myFuncConfig3 is nil")
	}

	// 2. Create a KisConnector configuration instance
	myConnConfig1 := config.NewConnConfig("ConnName1", "0.0.0.0:9998", common.REDIS, "redis-key", nil)
	if myConnConfig1 == nil {
		panic("myConnConfig1 is nil")
	}

	// 3. Bind the KisConnector configuration instance to the KisFunction configuration instance
	_ = myFuncConfig2.AddConnConfig(myConnConfig1)

	// 4. Create a KisFlow configuration instance
	myFlowConfig1 := config.NewFlowConfig("flowName1", common.FlowEnable)

	// 5. Create a KisFlow object
	flow1 := flow.NewKisFlow(myFlowConfig1)

	// 6. Link Functions to the Flow
	if err := flow1.Link(myFuncConfig1, nil); err != nil {
		panic(err)
	}
	if err := flow1.Link(myFuncConfig2, nil); err != nil {
		panic(err)
	}
	if err := flow1.Link(myFuncConfig3, nil); err != nil {
		panic(err)
	}

	// 7. Commit raw data
	_ = flow1.CommitRow("This is Data1 from Test")
	_ = flow1.CommitRow("This is Data2 from Test")
	_ = flow1.CommitRow("This is Data3 from Test")

	// 8. Execute flow1
	if err := flow1.Run(ctx); err != nil {
		panic(err)
	}
}

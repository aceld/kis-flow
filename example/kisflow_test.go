package test

import (
	"fmt"
	"kis-flow/flow/kis_config"
	"testing"
)

func TestNewFlowConfig(t *testing.T) {

	flowFuncParams1 := flow.KisFlowFunctionParam{
		Fid: "funcId1",
		Params: flow.FParam{
			"flowSetFunParam1": "value1",
			"flowSetFunParam2": "value2",
		},
	}

	flowFuncParams2 := flow.KisFlowFunctionParam{
		Fid: "funcId2",
		Params: flow.FParam{
			"default": "value1",
		},
	}

	myFlow1 := flow.NewFlowConfig("flowId", "flowName", 1)
	myFlow1.AppendFunctionConfig(flowFuncParams1)
	myFlow1.AppendFunctionConfig(flowFuncParams2)

	fmt.Printf("myFlow1: %+v\n", myFlow1)
}

func TestNewFuncConfig(t *testing.T) {
	source := flow.KisSource{
		Name: "公众号抖音商城户订单数据",
		Must: []string{"order_id", "user_id"},
	}

	option := flow.KisFuncOption{
		Cid:          "connector_id",
		RetryTimes:   3,
		RetryDuriton: 300,

		Params: flow.FParam{
			"param1": "value1",
			"param2": "value2",
		},
	}

	myFunc1 := flow.NewFuncConfig("funcId", "funcName", "Save", &source, &option)

	fmt.Printf("myFunc1: %+v\n", myFunc1)
}

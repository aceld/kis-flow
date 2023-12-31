package config

import (
	"fmt"
	"kis-flow/common"
	"testing"
)

func TestNewFlowConfig(t *testing.T) {

	flowFuncParams1 := KisFlowFunctionParam{
		Fid: "funcId1",
		Params: FParam{
			"flowSetFunParam1": "value1",
			"flowSetFunParam2": "value2",
		},
	}

	flowFuncParams2 := KisFlowFunctionParam{
		Fid: "funcId2",
		Params: FParam{
			"default": "value1",
		},
	}

	myFlow1 := NewFlowConfig("flowId", "flowName", 1)
	myFlow1.AppendFunctionConfig(flowFuncParams1)
	myFlow1.AppendFunctionConfig(flowFuncParams2)

	fmt.Printf("myFlow1: %+v\n", myFlow1)
}

func TestNewFuncConfig(t *testing.T) {
	source := KisSource{
		Name: "公众号抖音商城户订单数据",
		Must: []string{"order_id", "user_id"},
	}

	option := KisFuncOption{
		Cid:          "connector_id",
		RetryTimes:   3,
		RetryDuriton: 300,

		Params: FParam{
			"param1": "value1",
			"param2": "value2",
		},
	}

	myFunc1 := NewFuncConfig("funcId", "funcName", "Save", &source, &option)

	fmt.Printf("myFunc1: %+v\n", myFunc1)
}

func TestNewConnConfig(t *testing.T) {

	source := KisSource{
		Name: "公众号抖音商城户订单数据",
		Must: []string{"order_id", "user_id"},
	}

	option := KisFuncOption{
		Cid:          "connector_id",
		RetryTimes:   3,
		RetryDuriton: 300,

		Params: FParam{
			"param1": "value1",
			"param2": "value2",
		},
	}

	myFunc1 := NewFuncConfig("funcId", "funcName", "Save", &source, &option)

	connParams := FParam{
		"param1": "value1",
		"param2": "value2",
	}

	myConnector1 := NewConnConfig("connectorId", "connectorName", "0.0.0.0:9987,0.0.0.0:9997", common.REDIS, "key", connParams)

	if err := myConnector1.WithFunc(myFunc1); err != nil {
		fmt.Printf("WithFunc err: %s\n", err.Error())
	}

	fmt.Printf("myConnector1: %+v\n", myConnector1)
}

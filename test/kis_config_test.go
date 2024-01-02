package test

import (
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/log"
	"testing"
)

func TestNewFlowConfig(t *testing.T) {

	flowFuncParams1 := config.KisFlowFunctionParam{
		Fid: "funcId1",
		Params: config.FParam{
			"flowSetFunParam1": "value1",
			"flowSetFunParam2": "value2",
		},
	}

	flowFuncParams2 := config.KisFlowFunctionParam{
		Fid: "funcId2",
		Params: config.FParam{
			"default": "value1",
		},
	}

	myFlow1 := config.NewFlowConfig("flowId", "flowName", 1)
	myFlow1.AppendFunctionConfig(flowFuncParams1)
	myFlow1.AppendFunctionConfig(flowFuncParams2)

	log.Logger().InfoF("myFlow1: %+v\n", myFlow1)
}

func TestNewFuncConfig(t *testing.T) {
	source := config.KisSource{
		Name: "公众号抖音商城户订单数据",
		Must: []string{"order_id", "user_id"},
	}

	option := config.KisFuncOption{
		Cid:          "connector_id",
		RetryTimes:   3,
		RetryDuriton: 300,

		Params: config.FParam{
			"param1": "value1",
			"param2": "value2",
		},
	}

	myFunc1 := config.NewFuncConfig("funcId", "funcName", "Save", &source, &option)

	log.Logger().InfoF("myFunc1: %+v\n", myFunc1)
}

func TestNewConnConfig(t *testing.T) {

	source := config.KisSource{
		Name: "公众号抖音商城户订单数据",
		Must: []string{"order_id", "user_id"},
	}

	option := config.KisFuncOption{
		Cid:          "connector_id",
		RetryTimes:   3,
		RetryDuriton: 300,

		Params: config.FParam{
			"param1": "value1",
			"param2": "value2",
		},
	}

	myFunc1 := config.NewFuncConfig("funcId", "funcName", "Save", &source, &option)

	connParams := config.FParam{
		"param1": "value1",
		"param2": "value2",
	}

	myConnector1 := config.NewConnConfig("connectorId", "connectorName", "0.0.0.0:9987,0.0.0.0:9997", common.REDIS, "key", connParams)

	if err := myConnector1.WithFunc(myFunc1); err != nil {
		log.Logger().ErrorF("WithFunc err: %s\n", err.Error())
	}

	log.Logger().InfoF("myConnector1: %+v\n", myConnector1)
}

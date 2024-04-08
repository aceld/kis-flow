package test

import (
	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/config"
	"github.com/aceld/kis-flow/log"
	"testing"
)

func TestNewFuncConfig(t *testing.T) {
	source := config.KisSource{
		Name: "公众号抖音商城户订单数据",
		Must: []string{"order_id", "user_id"},
	}

	option := config.KisFuncOption{
		CName:         "connectorName1",
		RetryTimes:    3,
		RetryDuration: 300,

		Params: config.FParam{
			"param1": "value1",
			"param2": "value2",
		},
	}

	myFunc1 := config.NewFuncConfig("funcName1", common.S, &source, &option)

	log.Logger().Info("funcName1", myFunc1)
}

func TestNewFlowConfig(t *testing.T) {

	flowFuncParams1 := config.KisFlowFunctionParam{
		FuncName: "funcName1",
		Params: config.FParam{
			"flowSetFunParam1": "value1",
			"flowSetFunParam2": "value2",
		},
	}

	flowFuncParams2 := config.KisFlowFunctionParam{
		FuncName: "funcName2",
		Params: config.FParam{
			"default": "value1",
		},
	}

	myFlow1 := config.NewFlowConfig("flowName1", common.FlowEnable)
	myFlow1.AppendFunctionConfig(flowFuncParams1)
	myFlow1.AppendFunctionConfig(flowFuncParams2)

	log.Logger().Info("myFlow1", myFlow1)
}

func TestNewConnConfig(t *testing.T) {

	source := config.KisSource{
		Name: "公众号抖音商城户订单数据",
		Must: []string{"order_id", "user_id"},
	}

	option := config.KisFuncOption{
		CName:         "connectorName1",
		RetryTimes:    3,
		RetryDuration: 300,

		Params: config.FParam{
			"param1": "value1",
			"param2": "value2",
		},
	}

	myFunc1 := config.NewFuncConfig("funcName1", common.S, &source, &option)

	connParams := config.FParam{
		"param1": "value1",
		"param2": "value2",
	}

	myConnector1 := config.NewConnConfig("connectorName1", "0.0.0.0:9987,0.0.0.0:9997", common.REDIS, "key", connParams)

	if err := myConnector1.WithFunc(myFunc1); err != nil {
		log.Logger().Error("WithFunc", "err", err.Error())
	}

	log.Logger().Info("myConnector1", myConnector1)
}

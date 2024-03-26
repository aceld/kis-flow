package test

import (
	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/test/caas"
	"github.com/aceld/kis-flow/test/faas"
)

func init() {
	// 0. 注册Function 回调业务
	kis.Pool().FaaS("funcName1", faas.FuncDemo1Handler)
	kis.Pool().FaaS("funcName2", faas.FuncDemo2Handler)
	kis.Pool().FaaS("funcName3", faas.FuncDemo3Handler)
	kis.Pool().FaaS("funcName4", faas.FuncDemo4Handler)
	kis.Pool().FaaS("abortFunc", faas.AbortFuncHandler)         // abortFunc 业务
	kis.Pool().FaaS("dataReuseFunc", faas.DataReuseFuncHandler) // dataReuseFunc 业务
	kis.Pool().FaaS("noResultFunc", faas.NoResultFuncHandler)   // noResultFunc 业务
	kis.Pool().FaaS("jumpFunc", faas.JumpFuncHandler)           // jumpFunc 业务

	// 0. 注册ConnectorInit 和 Connector 回调业务
	kis.Pool().CaaSInit("ConnName1", caas.InitConnDemo1)
	kis.Pool().CaaS("ConnName1", "funcName2", common.S, caas.CaasDemoHanler1)

}

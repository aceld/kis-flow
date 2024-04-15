package test

import (
	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/test/caas"
	"github.com/aceld/kis-flow/test/faas"
)

func init() {
	// Register Function callback business
	kis.Pool().FaaS("funcName1", faas.FuncDemo1Handler)
	kis.Pool().FaaS("funcName2", faas.FuncDemo2Handler)
	kis.Pool().FaaS("funcName3", faas.FuncDemo3Handler)
	kis.Pool().FaaS("funcName4", faas.FuncDemo4Handler)
	kis.Pool().FaaS("abortFunc", faas.AbortFuncHandler)         // abortFunc business
	kis.Pool().FaaS("dataReuseFunc", faas.DataReuseFuncHandler) // dataReuseFunc business
	kis.Pool().FaaS("noResultFunc", faas.NoResultFuncHandler)   // noResultFunc business
	kis.Pool().FaaS("jumpFunc", faas.JumpFuncHandler)           // jumpFunc business

	// Register ConnectorInit and Connector callback business
	kis.Pool().CaaSInit("ConnName1", caas.InitConnDemo1)
	kis.Pool().CaaS("ConnName1", "funcName2", common.S, caas.CaasDemoHanler1)
}

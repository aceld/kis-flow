package kis

import (
	"context"

	"github.com/aceld/kis-flow/common"
)

/*
	Function Call
*/
// funcRouter
// key: Function Name
// value: FaaSDesc callback description for custom business
type funcRouter map[string]*FaaSDesc

// flowRouter
// key: Flow Name
// value: Flow
type flowRouter map[string]Flow

/*
	Connector Init
*/
// ConnInit Connector third-party storage initialization
type ConnInit func(conn Connector) error

// connInitRouter
// key: Connector Name
// value: ConnInit
type connInitRouter map[string]ConnInit

/*
	Connector Call
*/
// CaaS Connector storage read/write business implementation
type CaaS func(context.Context, Connector, Function, Flow, interface{}) (interface{}, error)

// connFuncRouter Maps CaaS callback storage business to FunctionName
// key: Function Name
// value: Connector storage read/write business implementation
type connFuncRouter map[string]CaaS

// connSL Splits connFuncRouter into two subtrees based on KisMode
// key: Function KisMode S/L
// value: connFuncRouter
type connSL map[common.KisMode]connFuncRouter

// connTree
// key: Connector Name
// value: connSL second-level tree
type connTree map[string]connSL

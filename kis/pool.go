package kis

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/log"
)

var _poolOnce sync.Once

// KisPool manages all Function and Flow configurations
type KisPool struct {
	fnRouter funcRouter   // All Function management routes
	fnLock   sync.RWMutex // fnRouter lock

	flowRouter flowRouter   // All flow objects
	flowLock   sync.RWMutex // flowRouter lock

	cInitRouter connInitRouter // All Connector initialization routes
	ciLock      sync.RWMutex   // cInitRouter lock

	cTree connTree     // All Connector management routes
	cLock sync.RWMutex // cTree lock
}

// Singleton
var _pool *KisPool

// Pool Singleton constructor
func Pool() *KisPool {
	_poolOnce.Do(func() {
		// Create KisPool object
		_pool = &KisPool{}

		// Initialize fnRouter
		_pool.fnRouter = make(funcRouter)

		// Initialize flowRouter
		_pool.flowRouter = make(flowRouter)

		// Initialize connTree
		_pool.cTree = make(connTree)
		_pool.cInitRouter = make(connInitRouter)
	})

	return _pool
}

func (pool *KisPool) AddFlow(name string, flow Flow) {
	pool.flowLock.Lock() // Write lock
	defer pool.flowLock.Unlock()

	if _, ok := pool.flowRouter[name]; !ok {
		pool.flowRouter[name] = flow
	} else {
		errString := fmt.Sprintf("Pool AddFlow Repeat FlowName=%s\n", name)
		panic(errString)
	}

	log.Logger().InfoF("Add FlowRouter FlowName=%s", name)
}

func (pool *KisPool) GetFlow(name string) Flow {
	pool.flowLock.RLock() // Read lock
	defer pool.flowLock.RUnlock()

	if flow, ok := pool.flowRouter[name]; ok {
		return flow
	} else {
		return nil
	}
}

// FaaS registers Function computation business logic, indexed and registered by Function Name
func (pool *KisPool) FaaS(fnName string, f FaaS) {

	// When registering the FaaS computation logic callback, create a FaaSDesc description object
	faaSDesc, err := NewFaaSDesc(fnName, f)
	if err != nil {
		panic(err)
	}

	pool.fnLock.Lock() // Write lock
	defer pool.fnLock.Unlock()

	if _, ok := pool.fnRouter[fnName]; !ok {
		// Register the FaaSDesc description object to fnRouter
		pool.fnRouter[fnName] = faaSDesc
	} else {
		errString := fmt.Sprintf("KisPoll FaaS Repeat FuncName=%s", fnName)
		panic(errString)
	}

	log.Logger().InfoF("Add KisPool FuncName=%s", fnName)
}

// CallFunction schedules Function
func (pool *KisPool) CallFunction(ctx context.Context, fnName string, flow Flow) error {
	pool.fnLock.RLock() // Read lock
	defer pool.fnLock.RUnlock()
	if funcDesc, ok := pool.fnRouter[fnName]; ok {

		// Parameters list for the scheduled Function
		params := make([]reflect.Value, 0, funcDesc.ArgNum)

		for _, argType := range funcDesc.ArgsType {

			// If it is a Flow type parameter, pass in the value of flow
			if isFlowType(argType) {
				params = append(params, reflect.ValueOf(flow))
				continue
			}

			// If it is a Context type parameter, pass in the value of ctx
			if isContextType(argType) {
				params = append(params, reflect.ValueOf(ctx))
				continue
			}

			// If it is a Slice type parameter, pass in the value of flow.Input()
			if isSliceType(argType) {

				// Deserialize the raw data in flow.Input() to data of type argType
				value, err := funcDesc.Serialize.UnMarshal(flow.Input(), argType)
				if err != nil {
					log.Logger().ErrorFX(ctx, "funcDesc.Serialize.DecodeParam err=%v", err)
				} else {
					params = append(params, value)
					continue
				}

			}

			// If the passed parameter is neither a Flow type, nor a Context type, nor a Slice type, it defaults to zero value
			params = append(params, reflect.Zero(argType))
		}

		// Call the computation logic of the current Function
		retValues := funcDesc.FuncValue.Call(params)

		// Extract the first return value, if it is nil, return nil
		ret := retValues[0].Interface()
		if ret == nil {
			return nil
		}

		// If the return value is of type error, return error
		return retValues[0].Interface().(error)

	}

	log.Logger().ErrorFX(ctx, "FuncName: %s Can not find in KisPool, Not Added.\n", fnName)

	return errors.New("FuncName: " + fnName + " Can not find in NsPool, Not Added.")
}

// CaaSInit registers Connector initialization business
func (pool *KisPool) CaaSInit(cname string, c ConnInit) {
	pool.ciLock.Lock() // Write lock
	defer pool.ciLock.Unlock()

	if _, ok := pool.cInitRouter[cname]; !ok {
		pool.cInitRouter[cname] = c
	} else {
		errString := fmt.Sprintf("KisPool Reg CaaSInit Repeat CName=%s\n", cname)
		panic(errString)
	}

	log.Logger().InfoF("Add KisPool CaaSInit CName=%s", cname)
}

// CallConnInit schedules ConnInit
func (pool *KisPool) CallConnInit(conn Connector) error {
	pool.ciLock.RLock() // Read lock
	defer pool.ciLock.RUnlock()

	init, ok := pool.cInitRouter[conn.GetName()]

	if !ok {
		panic(errors.New(fmt.Sprintf("init connector cname = %s not reg..", conn.GetName())))
	}

	return init(conn)
}

// CaaS registers Connector Call business
func (pool *KisPool) CaaS(cname string, fname string, mode common.KisMode, c CaaS) {
	pool.cLock.Lock() // Write lock
	defer pool.cLock.Unlock()

	if _, ok := pool.cTree[cname]; !ok {
		//cid First registration, does not exist, create a second-level tree NsConnSL
		pool.cTree[cname] = make(connSL)

		// Initialize various FunctionMode
		pool.cTree[cname][common.S] = make(connFuncRouter)
		pool.cTree[cname][common.L] = make(connFuncRouter)
	}

	if _, ok := pool.cTree[cname][mode][fname]; !ok {
		pool.cTree[cname][mode][fname] = c
	} else {
		errString := fmt.Sprintf("CaaS Repeat CName=%s, FName=%s, Mode =%s\n", cname, fname, mode)
		panic(errString)
	}

	log.Logger().InfoF("Add KisPool CaaS CName=%s, FName=%s, Mode =%s", cname, fname, mode)
}

// CallConnector schedules Connector
func (pool *KisPool) CallConnector(ctx context.Context, flow Flow, conn Connector, args interface{}) (interface{}, error) {
	pool.cLock.RLock() // Read lock
	defer pool.cLock.RUnlock()
	fn := flow.GetThisFunction()
	fnConf := fn.GetConfig()
	mode := common.KisMode(fnConf.FMode)

	if callback, ok := pool.cTree[conn.GetName()][mode][fnConf.FName]; ok {
		return callback(ctx, conn, fn, flow, args)
	}

	log.Logger().ErrorFX(ctx, "CName:%s FName:%s mode:%s Can not find in KisPool, Not Added.\n", conn.GetName(), fnConf.FName, mode)

	return nil, errors.New(fmt.Sprintf("CName:%s FName:%s mode:%s Can not find in KisPool, Not Added.", conn.GetName(), fnConf.FName, mode))
}

// GetFlows retrieves all Flows
func (pool *KisPool) GetFlows() []Flow {
	pool.flowLock.RLock() // Read lock
	defer pool.flowLock.RUnlock()

	var flows []Flow

	for _, flow := range pool.flowRouter {
		flows = append(flows, flow)
	}

	return flows
}

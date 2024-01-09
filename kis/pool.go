package kis

import (
	"context"
	"errors"
	"fmt"
	"kis-flow/common"
	"kis-flow/log"
	"sync"
)

var _poolOnce sync.Once

//  kisPool 用于管理全部的Function和Flow配置的池子
type kisPool struct {
	fnRouter funcRouter   // 全部的Function管理路由
	fnLock   sync.RWMutex // fnRouter 锁

	flowRouter flowRouter   // 全部的flow对象
	flowLock   sync.RWMutex // flowRouter 锁

	cInitRouter connInitRouter // 全部的Connector初始化路由
	ciLock      sync.RWMutex   // cInitRouter 锁

	cTree      connTree             //全部Connector管理路由
	connectors map[string]Connector // 全部的Connector对象
	cLock      sync.RWMutex         // cTree 锁
}

// 单例
var _pool *kisPool

// Pool 单例构造
func Pool() *kisPool {
	_poolOnce.Do(func() {
		//创建kisPool对象
		_pool = new(kisPool)

		// fnRouter初始化
		_pool.fnRouter = make(funcRouter)

		// flowRouter初始化
		_pool.flowRouter = make(flowRouter)

		// connTree初始化
		_pool.cTree = make(connTree)
		_pool.cInitRouter = make(connInitRouter)
		_pool.connectors = make(map[string]Connector)
	})

	return _pool
}

func (pool *kisPool) AddFlow(name string, flow Flow) {
	pool.flowLock.Lock() // 写锁
	defer pool.flowLock.Unlock()

	if _, ok := pool.flowRouter[name]; !ok {
		pool.flowRouter[name] = flow
	} else {
		errString := fmt.Sprintf("Pool AddFlow Repeat FlowName=%s\n", name)
		panic(errString)
	}

	log.Logger().InfoF("Add FlowRouter FlowName=%s\n", name)
}

func (pool *kisPool) GetFlow(name string) Flow {
	pool.flowLock.RLock() // 读锁
	defer pool.flowLock.RUnlock()

	if flow, ok := pool.flowRouter[name]; ok {
		return flow
	} else {
		return nil
	}
}

// FaaS 注册 Function 计算业务逻辑, 通过Function Name 索引及注册
func (pool *kisPool) FaaS(fnName string, f FaaS) {
	pool.fnLock.Lock() // 写锁
	defer pool.fnLock.Unlock()

	if _, ok := pool.fnRouter[fnName]; !ok {
		pool.fnRouter[fnName] = f
	} else {
		errString := fmt.Sprintf("KisPoll FaaS Repeat FuncName=%s", fnName)
		panic(errString)
	}

	log.Logger().InfoF("Add KisPool FuncName=%s", fnName)
}

// CallFunction 调度 Function
func (pool *kisPool) CallFunction(ctx context.Context, fnName string, flow Flow) error {

	if f, ok := pool.fnRouter[fnName]; ok {
		return f(ctx, flow)
	}

	log.Logger().ErrorFX(ctx, "FuncName: %s Can not find in KisPool, Not Added.\n", fnName)

	return errors.New("FuncName: " + fnName + " Can not find in NsPool, Not Added.")
}

// CaaSInit 注册Connector初始化业务
func (pool *kisPool) CaaSInit(cname string, c ConnInit) {
	pool.ciLock.Lock() // 写锁
	defer pool.ciLock.Unlock()

	if _, ok := pool.cInitRouter[cname]; !ok {
		pool.cInitRouter[cname] = c
	} else {
		errString := fmt.Sprintf("KisPool Reg CaaSInit Repeat CName=%s\n", cname)
		panic(errString)
	}

	log.Logger().InfoF("Add KisPool CaaSInit CName=%s", cname)
}

// CallConnInit 调度 ConnInit
func (pool *kisPool) CallConnInit(conn Connector) error {
	pool.ciLock.RLock() // 读锁
	defer pool.ciLock.RUnlock()

	init, ok := pool.cInitRouter[conn.GetName()]

	if !ok {
		panic(errors.New(fmt.Sprintf("init connector cname = %s not reg..", conn.GetName())))
	}

	return init(conn)
}

// CaaS 注册Connector Call业务
func (pool *kisPool) CaaS(cname string, fname string, mode common.KisMode, c CaaS) {
	pool.cLock.Lock() // 写锁
	defer pool.cLock.Unlock()

	if _, ok := pool.cTree[cname]; !ok {
		//cid 首次注册，不存在，创建二级树NsConnSL
		pool.cTree[cname] = make(connSL)

		//初始化各类型FunctionMode
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

// CallConnector 调度 Connector
func (pool *kisPool) CallConnector(ctx context.Context, flow Flow, conn Connector, args interface{}) error {
	fn := flow.GetThisFunction()
	fnConf := fn.GetConfig()
	mode := common.KisMode(fnConf.FMode)

	if callback, ok := pool.cTree[conn.GetName()][mode][fnConf.FName]; ok {
		return callback(ctx, conn, fn, flow, args)
	}

	log.Logger().ErrorFX(ctx, "CName:%s FName:%s mode:%s Can not find in KisPool, Not Added.\n", conn.GetName(), fnConf.FName, mode)

	return errors.New(fmt.Sprintf("CName:%s FName:%s mode:%s Can not find in KisPool, Not Added.", conn.GetName(), fnConf.FName, mode))
}

package function

import (
	"context"
	"errors"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/id"
	"kis-flow/kis"
)

type BaseFunction struct {
	// Id , KisFunction的实例ID，用于KisFlow内部区分不同的实例对象
	Id     string
	Config *config.KisFuncConfig

	// flow
	flow kis.Flow //上下文环境KisFlow

	// connector
	connector kis.Connector

	// link
	N kis.Function //下一个流计算Function
	P kis.Function //上一个流计算Function
}

// Call
// BaseFunction 为空实现，目的为了让其他具体类型的KisFunction，如KisFunction_V 来继承BaseFuncion来重写此方法
func (base *BaseFunction) Call(ctx context.Context, flow kis.Flow) error { return nil }

func (base *BaseFunction) Next() kis.Function {
	return base.N
}

func (base *BaseFunction) Prev() kis.Function {
	return base.P
}

func (base *BaseFunction) SetN(f kis.Function) {
	base.N = f
}

func (base *BaseFunction) SetP(f kis.Function) {
	base.P = f
}

func (base *BaseFunction) SetConfig(s *config.KisFuncConfig) error {
	if s == nil {
		return errors.New("KisFuncConfig is nil")
	}

	base.Config = s

	return nil
}

func (base *BaseFunction) GetId() string {
	return base.Id
}

func (base *BaseFunction) GetPrevId() string {
	if base.P == nil {
		//Function为首结点
		return common.FunctionIdFirstVirtual
	}
	return base.P.GetId()
}

func (base *BaseFunction) GetNextId() string {
	if base.N == nil {
		//Function为尾结点
		return common.FunctionIdLastVirtual
	}
	return base.N.GetId()
}

func (base *BaseFunction) GetConfig() *config.KisFuncConfig {
	return base.Config
}

func (base *BaseFunction) SetFlow(f kis.Flow) error {
	if f == nil {
		return errors.New("KisFlow is nil")
	}
	base.flow = f
	return nil
}

func (base *BaseFunction) GetFlow() kis.Flow {
	return base.flow
}

// AddConnector 给当前Function实例添加一个Connector
func (base *BaseFunction) AddConnector(conn kis.Connector) error {
	if conn == nil {
		return errors.New("conn is nil")
	}

	base.connector = conn

	return nil
}

// GetConnector 获取当前Function实例所关联的Connector
func (base *BaseFunction) GetConnector() kis.Connector {
	return base.connector
}

func (base *BaseFunction) CreateId() {
	base.Id = id.KisID(common.KisIdTypeFunction)
}

// NewKisFunction 创建一个NsFunction
// flow: 当前所属的flow实例
// s : 当前function的配置策略
func NewKisFunction(flow kis.Flow, config *config.KisFuncConfig) kis.Function {
	var f kis.Function

	//工厂生产泛化对象
	switch common.KisMode(config.FMode) {
	case common.V:
		f = new(KisFunctionV)
		break
	case common.S:
		f = new(KisFunctionS)
	case common.L:
		f = new(KisFunctionL)
	case common.C:
		f = new(KisFunctionC)
	case common.E:
		f = new(KisFunctionE)
	default:
		//LOG ERROR
		return nil
	}

	// 生成随机实例唯一ID
	f.CreateId()

	// 设置基础信息属性
	if err := f.SetConfig(config); err != nil {
		panic(err)
	}

	// 设置Flow
	if err := f.SetFlow(flow); err != nil {
		panic(err)
	}

	return f
}

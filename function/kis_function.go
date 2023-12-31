package function

import (
	"context"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/flow"
)

// KisFunction 流式计算基础计算模块，KisFunction是一条流式计算的基本计算逻辑单元，
// 			   任意个KisFunction可以组合成一个KisFlow
type KisFunction interface {
	// Call 执行流式计算逻辑
	Call(ctx context.Context, flow *flow.KisFlow) error

	// SetConfig 给当前Function实例配置策略
	SetConfig(s *config.KisFuncConfig) error
	// GetConfig 获取当前Function实例配置策略
	GetConfig() *config.KisFuncConfig

	// SetFlow 给当前Function实例设置所依赖的Flow实例
	SetFlow(f *flow.KisFlow) error
	// GetFlow 获取当前Functioin实力所依赖的Flow
	GetFlow() *flow.KisFlow

	// SetConnId 如果当前Function为S或者L 那么建议设置当前Funciton所关联的Connector
	SetConnId(string)
	// GetConnId 获取所关联的Connector CID
	GetConnId() string

	// GetPrevId 获取当前Function上一个Function节点FID
	GetPrevId() string
	// GetNextId 获取当前Function下一个Function节点FID
	GetNextId() string
	// GetId 获取当前Function的FID
	GetId() string

	// CreateKisId 给当前Funciton实力生成一个随机的实例KisID
	CreateKisId()
	// GetKisId 获取当前Function的唯一实例KisID
	GetKisId() string

	// Next 返回下一层计算流Function，如果当前层为最后一层，则返回nil
	Next() KisFunction
	// Prev 返回上一层计算流Function，如果当前层为最后一层，则返回nil
	Prev() KisFunction
	// SetN 设置下一层Function实例
	SetN(f KisFunction)
	// SetP 设置上一层Function实例
	SetP(f KisFunction)
}

// NewKisFunction 创建一个NsFunction
// flow: 当前所属的flow实例
// s : 当前function的配置策略
func NewKisFunction(flow *flow.KisFlow, config *config.KisFuncConfig) KisFunction {
	var f KisFunction

	//工厂生产泛化对象
	switch common.KisMode(config.Fmode) {
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

	//设置基础信息属性
	if err := f.SetConfig(config); err != nil {
		panic(err)
	}

	if err := f.SetFlow(flow); err != nil {
		panic(err)
	}

	if config.Option.Cid != "" {
		f.SetConnId(config.Option.Cid)
	}

	// 生成随机实力唯一ID
	f.CreateKisId()

	return f
}

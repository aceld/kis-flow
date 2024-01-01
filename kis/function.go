package kis

import (
	"context"
	"kis-flow/config"
)

// Function 流式计算基础计算模块，KisFunction是一条流式计算的基本计算逻辑单元，
// 			   任意个KisFunction可以组合成一个KisFlow
type Function interface {
	// Call 执行流式计算逻辑
	Call(ctx context.Context, flow Flow) error

	// SetConfig 给当前Function实例配置策略
	SetConfig(s *config.KisFuncConfig) error
	// GetConfig 获取当前Function实例配置策略
	GetConfig() *config.KisFuncConfig

	// SetFlow 给当前Function实例设置所依赖的Flow实例
	SetFlow(f Flow) error
	// GetFlow 获取当前Functioin实力所依赖的Flow
	GetFlow() Flow

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
	Next() Function
	// Prev 返回上一层计算流Function，如果当前层为最后一层，则返回nil
	Prev() Function
	// SetN 设置下一层Function实例
	SetN(f Function)
	// SetP 设置上一层Function实例
	SetP(f Function)
}

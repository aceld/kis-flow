package kis

import (
	"context"
	"kis-flow/config"
)

type Connector interface {
	// Init 初始化Connector所关联的存储引擎链接等
	Init() error
	// Call 调用Connector 外挂存储逻辑的读写操作
	Call(ctx context.Context, flow Flow, args interface{}) error
	// GetId 获取Connector的ID
	GetId() string
	// GetName 获取Connector的名称
	GetName() string
	// GetConfig 获取Connector的配置信息
	GetConfig() *config.KisConnConfig
}

package conn

import (
	"context"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/id"
	"kis-flow/kis"
	"sync"
)

type KisConnector struct {
	// Connector ID
	CId string
	// Connector Name
	CName string
	// Connector Config
	Conf *config.KisConnConfig

	// Connector Init
	onceInit sync.Once
}

// NewKisConnector 根据配置策略创建一个KisConnector
func NewKisConnector(config *config.KisConnConfig) *KisConnector {
	conn := new(KisConnector)
	conn.CId = id.KisID(common.KisIdTypeConnnector)
	conn.CName = config.CName
	conn.Conf = config

	return conn
}

// Init 初始化Connector所关联的存储引擎链接等
func (conn *KisConnector) Init() error {
	var err error

	//一个Connector只能执行初始化业务一次
	conn.onceInit.Do(func() {
		err = kis.Pool().CallConnInit(conn)
	})

	return err
}

// Call 调用Connector 外挂存储逻辑的读写操作
func (conn *KisConnector) Call(ctx context.Context, flow kis.Flow, args interface{}) error {
	if err := kis.Pool().CallConnector(ctx, flow, conn, args); err != nil {
		return err
	}

	return nil
}

func (conn *KisConnector) GetName() string {
	return conn.CName
}

func (conn *KisConnector) GetConfig() *config.KisConnConfig {
	return conn.Conf
}

func (conn *KisConnector) GetId() string {
	return conn.CId
}

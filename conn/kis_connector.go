package conn

import (
	"context"
	"sync"

	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/config"
	"github.com/aceld/kis-flow/id"
	"github.com/aceld/kis-flow/kis"
)

// KisConnector represents a KisConnector instance
type KisConnector struct {
	// Connector ID
	CId string
	// Connector Name
	CName string
	// Connector Config
	Conf *config.KisConnConfig

	// Connector Init
	onceInit sync.Once

	// KisConnector's custom temporary data
	metaData map[string]interface{}
	// Lock for reading and writing metaData
	mLock sync.RWMutex
}

// NewKisConnector creates a KisConnector based on the given configuration
func NewKisConnector(config *config.KisConnConfig) *KisConnector {
	conn := new(KisConnector)
	conn.CId = id.KisID(common.KisIDTypeConnector)
	conn.CName = config.CName
	conn.Conf = config
	conn.metaData = make(map[string]interface{})

	return conn
}

// Init initializes the connection to the associated storage engine of the Connector
func (conn *KisConnector) Init() error {
	var err error

	// The initialization business of a Connector can only be executed once
	conn.onceInit.Do(func() {
		err = kis.Pool().CallConnInit(conn)
	})

	return err
}

// Call invokes the read-write operations of the external storage logic through the Connector
func (conn *KisConnector) Call(ctx context.Context, flow kis.Flow, args interface{}) (interface{}, error) {
	var result interface{}
	var err error

	result, err = kis.Pool().CallConnector(ctx, flow, conn, args)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetName returns the name of the Connector
func (conn *KisConnector) GetName() string {
	return conn.CName
}

// GetConfig returns the configuration of the Connector
func (conn *KisConnector) GetConfig() *config.KisConnConfig {
	return conn.Conf
}

// GetID returns the ID of the Connector
func (conn *KisConnector) GetID() string {
	return conn.CId
}

// GetMetaData gets the temporary data of the current Connector
func (conn *KisConnector) GetMetaData(key string) interface{} {
	conn.mLock.RLock()
	defer conn.mLock.RUnlock()

	data, ok := conn.metaData[key]
	if !ok {
		return nil
	}

	return data
}

// SetMetaData sets the temporary data of the current Connector
func (conn *KisConnector) SetMetaData(key string, value interface{}) {
	conn.mLock.Lock()
	defer conn.mLock.Unlock()

	conn.metaData[key] = value
}

package kis

import (
	"context"
	"time"

	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/config"
)

type Flow interface {
	// Run schedules the Flow, sequentially dispatching and executing Functions in the Flow
	Run(ctx context.Context) error
	// Link connects the Functions in the Flow according to the configuration in the config file, and the Flow's configuration will also be updated
	Link(fConf *config.KisFuncConfig, fParams config.FParam) error
	// AppendNewFunction appends a new Function to the Flow
	AppendNewFunction(fConf *config.KisFuncConfig, fParams config.FParam) error
	// CommitRow submits Flow data to the upcoming Function layer
	CommitRow(row interface{}) error
	// CommitRowBatch submits Flow data to the upcoming Function layer (batch submission)
	// row: Must be a slice
	CommitRowBatch(row interface{}) error
	// Input gets the input source data of the currently executing Function in the Flow
	Input() common.KisRowArr
	// GetName gets the name of the Flow
	GetName() string
	// GetThisFunction gets the currently executing Function
	GetThisFunction() Function
	// GetThisFuncConf gets the configuration of the currently executing Function
	GetThisFuncConf() *config.KisFuncConfig
	// GetConnector gets the Connector of the currently executing Function
	GetConnector() (Connector, error)
	// GetConnConf gets the configuration of the Connector of the currently executing Function
	GetConnConf() (*config.KisConnConfig, error)
	// GetConfig gets the configuration of the current Flow
	GetConfig() *config.KisFlowConfig
	// GetFuncConfigByName gets the configuration of the current Flow by Function name
	GetFuncConfigByName(funcName string) *config.KisFuncConfig
	// Next carries the Action actions of the next layer Function that the current Flow is executing
	Next(acts ...ActionFunc) error
	// GetCacheData gets the cached data of the current Flow
	GetCacheData(key string) interface{}
	// SetCacheData sets the cached data of the current Flow
	SetCacheData(key string, value interface{}, Exp time.Duration)
	// GetMetaData gets the temporary data of the current Flow
	GetMetaData(key string) interface{}
	// SetMetaData sets the temporary data of the current Flow
	SetMetaData(key string, value interface{})
	// GetFuncParam gets the default parameters of the current Flow's currently executing Function, retrieving a key-value pair
	GetFuncParam(key string) string
	// GetFuncParamAll gets the default parameters of the current Flow's currently executing Function, retrieving all Key-Value pairs
	GetFuncParamAll() config.FParam
	// GetFuncParamsAllFuncs gets the FuncParams of all Functions in the Flow, retrieving all Key-Value pairs
	GetFuncParamsAllFuncs() map[string]config.FParam
	// Fork gets a copy of the Flow (deep copy)
	Fork(ctx context.Context) Flow
	// GetID gets the Id of the Flow
	GetID() string
}

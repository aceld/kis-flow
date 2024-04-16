package function

import (
	"context"
	"errors"
	"sync"

	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/config"
	"github.com/aceld/kis-flow/id"
	"github.com/aceld/kis-flow/kis"
)

type BaseFunction struct {
	// Id, the instance ID of KisFunction, used to differentiate different instance objects within KisFlow
	Id     string
	Config *config.KisFuncConfig

	// flow
	flow kis.Flow // Context environment KisFlow

	// connector
	connector kis.Connector

	// Custom temporary data of Function
	metaData map[string]interface{}
	// Manage the read-write lock of metaData
	mLock sync.RWMutex

	// link
	N kis.Function // Next flow computing Function
	P kis.Function // Previous flow computing Function
}

// Call
// BaseFunction is an empty implementation, designed to allow other specific types of KisFunction,
// such as KisFunction_V, to inherit BaseFuncion and override this method
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

func (base *BaseFunction) GetID() string {
	return base.Id
}

func (base *BaseFunction) GetPrevId() string {
	if base.P == nil {
		// Function is the first node
		return common.FunctionIDFirstVirtual
	}
	return base.P.GetID()
}

func (base *BaseFunction) GetNextId() string {
	if base.N == nil {
		// Function is the last node
		return common.FunctionIDLastVirtual
	}
	return base.N.GetID()
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

// AddConnector adds a Connector to the current Function instance
func (base *BaseFunction) AddConnector(conn kis.Connector) error {
	if conn == nil {
		return errors.New("conn is nil")
	}

	base.connector = conn

	return nil
}

// GetConnector gets the Connector associated with the current Function instance
func (base *BaseFunction) GetConnector() kis.Connector {
	return base.connector
}

func (base *BaseFunction) CreateId() {
	base.Id = id.KisID(common.KisIDTypeFunction)
}

// NewKisFunction creates a new NsFunction
// flow: the current belonging flow instance
// s: the configuration strategy of the current function
func NewKisFunction(flow kis.Flow, config *config.KisFuncConfig) kis.Function {
	var f kis.Function

	// Factory produces generic objects
	switch common.KisMode(config.FMode) {
	case common.V:
		f = NewKisFunctionV()
	case common.S:
		f = NewKisFunctionS()
	case common.L:
		f = NewKisFunctionL()
	case common.C:
		f = NewKisFunctionC()
	case common.E:
		f = NewKisFunctionE()
	default:
		// LOG ERROR
		return nil
	}

	// Generate a random unique instance ID
	f.CreateId()

	// Set basic information attributes
	if err := f.SetConfig(config); err != nil {
		panic(err)
	}

	// Set Flow
	if err := f.SetFlow(flow); err != nil {
		panic(err)
	}

	return f
}

// GetMetaData gets the temporary data of the current Function
func (base *BaseFunction) GetMetaData(key string) interface{} {
	base.mLock.RLock()
	defer base.mLock.RUnlock()

	data, ok := base.metaData[key]
	if !ok {
		return nil
	}

	return data
}

// SetMetaData sets the temporary data of the current Function
func (base *BaseFunction) SetMetaData(key string, value interface{}) {
	base.mLock.Lock()
	defer base.mLock.Unlock()

	base.metaData[key] = value
}

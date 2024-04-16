package kis

import (
	"context"

	"github.com/aceld/kis-flow/config"
)

// Function is the basic computation unit of streaming computation. KisFunction is a basic logical unit of streaming computation, any number of KisFunctions can be combined into a KisFlow
type Function interface {
	// Call executes the streaming computation logic
	Call(ctx context.Context, flow Flow) error

	// SetConfig configures the current Function instance
	SetConfig(s *config.KisFuncConfig) error
	// GetConfig retrieves the configuration of the current Function instance
	GetConfig() *config.KisFuncConfig

	// SetFlow sets the Flow instance that the current Function instance depends on
	SetFlow(f Flow) error
	// GetFlow retrieves the Flow instance that the current Function instance depends on
	GetFlow() Flow

	// AddConnector adds a Connector to the current Function instance
	AddConnector(conn Connector) error
	// GetConnector retrieves the Connector associated with the current Function instance
	GetConnector() Connector

	// CreateId generates a random KisID for the current Function instance
	CreateId()
	// GetID retrieves the FID of the current Function
	GetID() string
	// GetPrevId retrieves the FID of the previous Function node of the current Function
	GetPrevId() string
	// GetNextId retrieves the FID of the next Function node of the current Function
	GetNextId() string

	// Next returns the next layer of the computation flow Function. If the current layer is the last layer, it returns nil
	Next() Function
	// Prev returns the previous layer of the computation flow Function. If the current layer is the last layer, it returns nil
	Prev() Function
	// SetN sets the next Function instance
	SetN(f Function)
	// SetP sets the previous Function instance
	SetP(f Function)
	// GetMetaData retrieves the temporary data of the current Function
	GetMetaData(key string) interface{}
	// SetMetaData sets the temporary data of the current Function
	SetMetaData(key string, value interface{})
}

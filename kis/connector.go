package kis

import (
	"context"

	"github.com/aceld/kis-flow/config"
)

// Connector defines the interface for connectors associated with external storage.
type Connector interface {
	// Init initializes the connection to the storage engine associated with the Connector.
	Init() error
	// Call invokes the read-write operations of the external storage logic.
	Call(ctx context.Context, flow Flow, args interface{}) (interface{}, error)
	// GetID returns the ID of the Connector.
	GetID() string
	// GetName returns the name of the Connector.
	GetName() string
	// GetConfig returns the configuration information of the Connector.
	GetConfig() *config.KisConnConfig
	// GetMetaData gets the temporary data of the current Connector.
	GetMetaData(key string) interface{}
	// SetMetaData sets the temporary data of the current Connector.
	SetMetaData(key string, value interface{})
}

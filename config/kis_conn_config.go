package config

import (
	"fmt"

	"github.com/aceld/kis-flow/common"
)

// KisConnConfig describes the KisConnector strategy configuration
type KisConnConfig struct {
	KisType    string             `yaml:"kistype"` // Configuration type
	CName      string             `yaml:"cname"`   // Unique descriptive identifier
	AddrString string             `yaml:"addrs"`   // Base storage medium address
	Type       common.KisConnType `yaml:"type"`    // Storage medium engine type: "Mysql", "Redis", "Kafka", etc.
	Key        string             `yaml:"key"`     // Identifier for a single storage: Key name for Redis, Table name for Mysql, Topic name for Kafka, etc.
	Params     map[string]string  `yaml:"params"`  // Custom parameters in the configuration information

	// NsFuncionID bound to storage reading
	Load []string `yaml:"load"`
	Save []string `yaml:"save"`
}

// NewConnConfig creates a KisConnector strategy configuration object, used to describe a KisConnector information
func NewConnConfig(cName string, addr string, t common.KisConnType, key string, param map[string]string) *KisConnConfig {
	strategy := new(KisConnConfig)
	strategy.CName = cName
	strategy.AddrString = addr
	strategy.Type = t
	strategy.Key = key
	strategy.Params = param

	return strategy
}

// WithFunc binds Connector to Function
func (cConfig *KisConnConfig) WithFunc(fConfig *KisFuncConfig) error {

	switch common.KisMode(fConfig.FMode) {
	case common.S:
		cConfig.Save = append(cConfig.Save, fConfig.FName)
	case common.L:
		cConfig.Load = append(cConfig.Load, fConfig.FName)
	default:
		return fmt.Errorf("Wrong KisMode %s", fConfig.FMode)
	}

	return nil
}

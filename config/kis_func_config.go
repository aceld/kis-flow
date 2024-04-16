package config

import (
	"fmt"

	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/log"
)

// FParam represents the type for custom fixed configuration parameters for the Function in the current Flow
type FParam map[string]string

// KisSource represents the business source of the current Function
type KisSource struct {
	Name string   `yaml:"name"` // Description of the data source for this layer Function
	Must []string `yaml:"must"` // Required fields for the source
}

// KisFuncOption represents optional configurations
type KisFuncOption struct {
	CName         string `yaml:"cname"`           // Connector name
	RetryTimes    int    `yaml:"retry_times"`     // Optional, maximum retry times for Function scheduling (excluding normal scheduling)
	RetryDuration int    `yaml:"return_duration"` // Optional, maximum time interval for each retry in Function scheduling (unit: ms)
	Params        FParam `yaml:"default_params"`  // Optional, custom fixed configuration parameters for the Function in the current Flow
}

// KisFuncConfig represents a KisFunction strategy configuration
type KisFuncConfig struct {
	KisType  string        `yaml:"kistype"`
	FName    string        `yaml:"fname"`
	FMode    string        `yaml:"fmode"`
	Source   KisSource     `yaml:"source"`
	Option   KisFuncOption `yaml:"option"`
	connConf *KisConnConfig
}

// NewFuncConfig creates a Function strategy configuration object, used to describe a KisFunction information
func NewFuncConfig(
	funcName string, mode common.KisMode,
	source *KisSource, option *KisFuncOption) *KisFuncConfig {

	config := new(KisFuncConfig)
	config.FName = funcName

	if source == nil {
		defaultSource := KisSource{
			Name: "unNamedSource",
		}
		source = &defaultSource
		log.Logger().InfoF("funcName NewConfig source is nil, funcName = %s, use default unNamed Source.", funcName)
	}
	config.Source = *source

	config.FMode = string(mode)

	/*
		// Functions S and L require the KisConnector parameters to be passed as they need to establish streaming relationships through Connector
		if mode == common.S || mode == common.L {
			if option == nil {
				log.Logger().ErrorF("Function S/L needs option->Cid\n")
				return nil
			} else if option.CName == "" {
				log.Logger().ErrorF("Function S/L needs option->Cid\n")
				return nil
			}
		}
	*/

	if option != nil {
		config.Option = *option
	}

	return config
}

// AddConnConfig WithConn binds Function to Connector
func (fConf *KisFuncConfig) AddConnConfig(cConf *KisConnConfig) error {
	if cConf == nil {
		return fmt.Errorf("KisConnConfig is nil")
	}

	// Function needs to be associated with Connector
	fConf.connConf = cConf

	// Connector needs to be associated with Function
	_ = cConf.WithFunc(fConf)

	// Update CName in Function configuration
	fConf.Option.CName = cConf.CName

	return nil
}

// GetConnConfig gets the Connector configuration
func (fConf *KisFuncConfig) GetConnConfig() (*KisConnConfig, error) {
	if fConf.connConf == nil {
		return nil, fmt.Errorf("KisFuncConfig.connConf not set")
	}

	return fConf.connConf, nil
}

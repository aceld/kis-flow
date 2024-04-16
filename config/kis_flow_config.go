package config

import "github.com/aceld/kis-flow/common"

// KisFlowFunctionParam represents the Id of a Function and carries fixed configuration parameters in a Flow configuration
type KisFlowFunctionParam struct {
	FuncName string `yaml:"fname"`  // Required
	Params   FParam `yaml:"params"` // Optional, custom fixed configuration parameters for the Function in the current Flow
}

// KisFlowConfig represents the object that spans the entire stream computing context environment
type KisFlowConfig struct {
	KisType  string                 `yaml:"kistype"`
	Status   int                    `yaml:"status"`
	FlowName string                 `yaml:"flow_name"`
	Flows    []KisFlowFunctionParam `yaml:"flows"`
}

// NewFlowConfig creates a Flow strategy configuration object, used to describe a KisFlow information
func NewFlowConfig(flowName string, enable common.KisOnOff) *KisFlowConfig {
	config := new(KisFlowConfig)
	config.FlowName = flowName
	config.Flows = make([]KisFlowFunctionParam, 0)

	config.Status = int(enable)

	return config
}

// AppendFunctionConfig adds a Function Config to the current Flow
func (fConfig *KisFlowConfig) AppendFunctionConfig(params KisFlowFunctionParam) {
	fConfig.Flows = append(fConfig.Flows, params)
}

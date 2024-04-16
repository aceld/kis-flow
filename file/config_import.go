package file

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/config"
	"github.com/aceld/kis-flow/flow"
	"github.com/aceld/kis-flow/kis"
	"github.com/aceld/kis-flow/metrics"
	yaml "gopkg.in/yaml.v2"
)

type allConfig struct {
	Flows map[string]*config.KisFlowConfig
	Funcs map[string]*config.KisFuncConfig
	Conns map[string]*config.KisConnConfig
}

// kisTypeFlowConfigure parses Flow configuration file in yaml format
func kisTypeFlowConfigure(all *allConfig, confData []byte, fileName string, kisType interface{}) error {
	flowCfg := new(config.KisFlowConfig)
	if ok := yaml.Unmarshal(confData, flowCfg); ok != nil {
		return fmt.Errorf("%s has wrong format kisType = %s", fileName, kisType)
	}

	// Skip the configuration loading if the Flow status is disabled
	if common.KisOnOff(flowCfg.Status) == common.FlowDisable {
		return nil
	}

	if _, ok := all.Flows[flowCfg.FlowName]; ok {
		return fmt.Errorf("%s set repeat flow_id:%s", fileName, flowCfg.FlowName)
	}

	// Add to the configuration set
	all.Flows[flowCfg.FlowName] = flowCfg

	return nil
}

// kisTypeFuncConfigure parses Function configuration file in yaml format
func kisTypeFuncConfigure(all *allConfig, confData []byte, fileName string, kisType interface{}) error {
	function := new(config.KisFuncConfig)
	if ok := yaml.Unmarshal(confData, function); ok != nil {
		return fmt.Errorf("%s has wrong format kisType = %s", fileName, kisType)
	}
	if _, ok := all.Funcs[function.FName]; ok {
		return fmt.Errorf("%s set repeat function_id:%s", fileName, function.FName)
	}

	// Add to the configuration set
	all.Funcs[function.FName] = function

	return nil
}

// kisTypeConnConfigure parses Connector configuration file in yaml format
func kisTypeConnConfigure(all *allConfig, confData []byte, fileName string, kisType interface{}) error {
	conn := new(config.KisConnConfig)
	if ok := yaml.Unmarshal(confData, conn); ok != nil {
		return fmt.Errorf("%s has wrong format kisType = %s", fileName, kisType)
	}

	if _, ok := all.Conns[conn.CName]; ok {
		return fmt.Errorf("%s set repeat conn_id:%s", fileName, conn.CName)
	}

	// Add to the configuration set
	all.Conns[conn.CName] = conn

	return nil
}

// kisTypeGlobalConfigure parses Global configuration file in yaml format
func kisTypeGlobalConfigure(confData []byte, fileName string, kisType interface{}) error {
	// Global configuration
	if ok := yaml.Unmarshal(confData, config.GlobalConfig); ok != nil {
		return fmt.Errorf("%s is wrong format kisType = %s", fileName, kisType)
	}

	// Start Metrics service
	metrics.RunMetrics()

	return nil
}

// parseConfigWalkYaml recursively parses all configuration files in yaml format and stores the configuration information in allConfig
func parseConfigWalkYaml(loadPath string) (*allConfig, error) {

	all := new(allConfig)

	all.Flows = make(map[string]*config.KisFlowConfig)
	all.Funcs = make(map[string]*config.KisFuncConfig)
	all.Conns = make(map[string]*config.KisConnConfig)

	err := filepath.Walk(loadPath, func(filePath string, info os.FileInfo, err error) error {
		// Validate the file extension
		if suffix := path.Ext(filePath); suffix != ".yml" && suffix != ".yaml" {
			return nil
		}

		// Read file content
		confData, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		confMap := make(map[string]interface{})

		// Validate yaml format
		if err := yaml.Unmarshal(confData, confMap); err != nil {
			return err
		}

		// Check if kisType exists
		var kisType interface{}

		kisType, ok := confMap["kistype"]
		if !ok {
			return fmt.Errorf("%s has no field [kistype]", filePath)
		}

		switch kisType {
		case common.KisIDTypeFlow:
			return kisTypeFlowConfigure(all, confData, filePath, kisType)

		case common.KisIDTypeFunction:
			return kisTypeFuncConfigure(all, confData, filePath, kisType)

		case common.KisIDTypeConnector:
			return kisTypeConnConfigure(all, confData, filePath, kisType)

		case common.KisIDTypeGlobal:
			return kisTypeGlobalConfigure(confData, filePath, kisType)

		default:
			return fmt.Errorf("%s set wrong kistype %s", filePath, kisType)
		}
	})

	if err != nil {
		return nil, err
	}

	return all, nil
}

func buildFlow(all *allConfig, fp config.KisFlowFunctionParam, newFlow kis.Flow, flowName string) error {
	// Load the Functions that the current Flow depends on
	if funcConfig, ok := all.Funcs[fp.FuncName]; !ok {
		return fmt.Errorf("FlowName [%s] need FuncName [%s], But has No This FuncName Config", flowName, fp.FuncName)
	} else {
		// flow add connector
		if funcConfig.Option.CName != "" {
			// Load the Connectors that the current Function depends on
			if connConf, ok := all.Conns[funcConfig.Option.CName]; !ok {
				return fmt.Errorf("FuncName [%s] need ConnName [%s], But has No This ConnName Config", fp.FuncName, funcConfig.Option.CName)
			} else {
				// Function Config associates with Connector Config
				_ = funcConfig.AddConnConfig(connConf)
			}
		}

		// flow add function
		if err := newFlow.AppendNewFunction(funcConfig, fp.Params); err != nil {
			return err
		}
	}

	return nil
}

// ConfigImportYaml recursively parses all configuration files in yaml format
func ConfigImportYaml(loadPath string) error {

	all, err := parseConfigWalkYaml(loadPath)
	if err != nil {
		return err
	}

	for flowName, flowConfig := range all.Flows {

		// Build a new Flow
		newFlow := flow.NewKisFlow(flowConfig)

		for _, fp := range flowConfig.Flows {
			if err := buildFlow(all, fp, newFlow, flowName); err != nil {
				return err
			}
		}

		// Add the flow to FlowPool
		kis.Pool().AddFlow(flowName, newFlow)
	}

	return nil
}

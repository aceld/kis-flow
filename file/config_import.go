package file

import (
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/flow"
	"kis-flow/kis"
	"os"
	"path"
	"path/filepath"
)

type allConfig struct {
	Flows map[string]*config.KisFlowConfig
	Funcs map[string]*config.KisFuncConfig
	Conns map[string]*config.KisConnConfig
}

// kisTypeFlowConfigure 解析Flow配置文件，yaml格式
func kisTypeFlowConfigure(all *allConfig, confData []byte, fileName string, kisType interface{}) error {
	flow := new(config.KisFlowConfig)
	if ok := yaml.Unmarshal(confData, flow); ok != nil {
		return errors.New(fmt.Sprintf("%s has wrong format kisType = %s", fileName, kisType))
	}

	// 如果FLow状态为关闭，则不做配置加载
	if common.KisOnOff(flow.Status) == common.FlowDisable {
		return nil
	}

	if _, ok := all.Flows[flow.FlowName]; ok {
		return errors.New(fmt.Sprintf("%s set repeat flow_id:%s", fileName, flow.FlowName))
	}

	// 加入配置集合中
	all.Flows[flow.FlowName] = flow

	return nil
}

// kisTypeFuncConfigure 解析Function配置文件，yaml格式
func kisTypeFuncConfigure(all *allConfig, confData []byte, fileName string, kisType interface{}) error {
	function := new(config.KisFuncConfig)
	if ok := yaml.Unmarshal(confData, function); ok != nil {
		return errors.New(fmt.Sprintf("%s has wrong format kisType = %s", fileName, kisType))
	}
	if _, ok := all.Funcs[function.FName]; ok {
		return errors.New(fmt.Sprintf("%s set repeat function_id:%s", fileName, function.FName))
	}

	// 加入配置集合中
	all.Funcs[function.FName] = function

	return nil
}

// kisTypeConnConfigure 解析Connector配置文件，yaml格式
func kisTypeConnConfigure(all *allConfig, confData []byte, fileName string, kisType interface{}) error {
	conn := new(config.KisConnConfig)
	if ok := yaml.Unmarshal(confData, conn); ok != nil {
		return errors.New(fmt.Sprintf("%s is wrong format nsType = %s", fileName, kisType))
	}

	if _, ok := all.Conns[conn.CName]; ok {
		return errors.New(fmt.Sprintf("%s set repeat conn_id:%s", fileName, conn.CName))
	}

	// 加入配置集合中
	all.Conns[conn.CName] = conn

	return nil
}

// parseConfigWalkYaml 全盘解析配置文件，yaml格式, 讲配置信息解析到allConfig中
func parseConfigWalkYaml(loadPath string) (*allConfig, error) {

	all := new(allConfig)

	all.Flows = make(map[string]*config.KisFlowConfig)
	all.Funcs = make(map[string]*config.KisFuncConfig)
	all.Conns = make(map[string]*config.KisConnConfig)

	err := filepath.Walk(loadPath, func(filePath string, info os.FileInfo, err error) error {
		// 校验文件后缀是否合法
		if suffix := path.Ext(filePath); suffix != ".yml" && suffix != ".yaml" {
			return nil
		}

		// 读取文件内容
		confData, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}

		confMap := make(map[string]interface{})

		// 校验yaml合法性
		if err := yaml.Unmarshal(confData, confMap); err != nil {
			return err
		}

		// 判断kisType是否存在
		if kisType, ok := confMap["kistype"]; !ok {
			return errors.New(fmt.Sprintf("yaml file %s has no file [kistype]!", filePath))
		} else {
			switch kisType {
			case common.KisIdTypeFlow:
				return kisTypeFlowConfigure(all, confData, filePath, kisType)

			case common.KisIdTypeFunction:
				return kisTypeFuncConfigure(all, confData, filePath, kisType)

			case common.KisIdTypeConnnector:
				return kisTypeConnConfigure(all, confData, filePath, kisType)

			default:
				return errors.New(fmt.Sprintf("%s set wrong kistype %s", filePath, kisType))
			}
		}
	})

	if err != nil {
		return nil, err
	}

	return all, nil
}

func buildFlow(all *allConfig, fp config.KisFlowFunctionParam, newFlow kis.Flow, flowName string) error {
	//加载当前Flow依赖的Function
	if funcConfig, ok := all.Funcs[fp.FuncName]; !ok {
		return errors.New(fmt.Sprintf("FlowName [%s] need FuncName [%s], But has No This FuncName Config", flowName, fp.FuncName))
	} else {
		//flow add connector
		if funcConfig.Option.CName != "" {
			// 加载当前Function依赖的Connector
			if connConf, ok := all.Conns[funcConfig.Option.CName]; !ok {
				return errors.New(fmt.Sprintf("FuncName [%s] need ConnName [%s], But has No This ConnName Config", fp.FuncName, funcConfig.Option.CName))
			} else {
				// Function Config 关联 Connector Config
				_ = funcConfig.AddConnConfig(connConf)
			}
		}

		//flow add function
		if err := newFlow.Link(funcConfig, fp.Params); err != nil {
			return err
		}
	}

	return nil
}

// ConfigImportYaml 全盘解析配置文件，yaml格式
func ConfigImportYaml(loadPath string) error {

	all, err := parseConfigWalkYaml(loadPath)
	if err != nil {
		return err
	}

	for flowName, flowConfig := range all.Flows {

		// 构建一个Flow
		newFlow := flow.NewKisFlow(flowConfig)

		for _, fp := range flowConfig.Flows {
			if err := buildFlow(all, fp, newFlow, flowName); err != nil {
				return err
			}
		}

		//将flow添加到FlowPool中
		kis.Pool().AddFlow(flowName, newFlow)
	}

	return nil
}

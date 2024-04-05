package config

import (
	"errors"
	"fmt"

	"github.com/aceld/kis-flow/common"
)

// KisConnConfig KisConnector 策略配置
type KisConnConfig struct {
	// 配置类型
	KisType string `yaml:"kistype"`
	// 唯一描述标识
	CName string `yaml:"cname"`
	// 基础存储媒介地址
	AddrString string `yaml:"addrs"`
	// 存储媒介引擎类型"Mysql" "Redis" "Kafka"等
	Type common.KisConnType `yaml:"type"`
	// 一次存储的标识：如Redis为Key名称、Mysql为Table名称,Kafka为Topic名称等
	Key string `yaml:"key"`
	// 配置信息中的自定义参数
	Params map[string]string `yaml:"params"`
	// 存储读取所绑定的NsFuncionID
	Load []string `yaml:"load"`
	Save []string `yaml:"save"`
}

// NewConnConfig 创建一个KisConnector策略配置对象, 用于描述一个KisConnector信息
func NewConnConfig(cName string, addr string, t common.KisConnType, key string, param FParam) *KisConnConfig {
	strategy := new(KisConnConfig)
	strategy.CName = cName
	strategy.AddrString = addr

	strategy.Type = t
	strategy.Key = key
	strategy.Params = param

	return strategy
}

// WithFunc Connector与Function进行关系绑定
func (cConfig *KisConnConfig) WithFunc(fConfig *KisFuncConfig) error {

	switch common.KisMode(fConfig.FMode) {
	case common.S:
		cConfig.Save = append(cConfig.Save, fConfig.FName)
	case common.L:
		cConfig.Load = append(cConfig.Load, fConfig.FName)
	default:
		return errors.New(fmt.Sprintf("Wrong KisMode %s", fConfig.FMode))
	}

	return nil
}

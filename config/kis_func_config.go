package config

import (
	"errors"
	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/log"
)

// FParam 在当前Flow中Function定制固定配置参数类型
type FParam map[string]string

// KisSource 表示当前Function的业务源
type KisSource struct {
	Name string   `yaml:"name"` // 本层Function的数据源描述
	Must []string `yaml:"must"` // source必传字段
}

// KisFuncOption 可选配置
type KisFuncOption struct {
	CName         string `yaml:"cname"`           // 连接器Connector名称
	RetryTimes    int    `yaml:"retry_times"`     // 选填,Function调度重试(不包括正常调度)最大次数
	RetryDuration int    `yaml:"return_duration"` // 选填,Function调度每次重试最大时间间隔(单位:ms)
	Params        FParam `yaml:"default_params"`  // 选填,在当前Flow中Function定制固定配置参数
}

// KisFuncConfig 一个KisFunction策略配置
type KisFuncConfig struct {
	KisType  string        `yaml:"kistype"`
	FName    string        `yaml:"fname"`
	FMode    string        `yaml:"fmode"`
	Source   KisSource     `yaml:"source"`
	Option   KisFuncOption `yaml:"option"`
	connConf *KisConnConfig
}

// NewFuncConfig 创建一个Function策略配置对象, 用于描述一个KisFunction信息
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
		log.Logger().Info("funcName NewConfig source is nil, use default unNamed Source.", "funcName", funcName)
	}
	config.Source = *source

	config.FMode = string(mode)

	/*
		// FunctionS 和 L 需要必传KisConnector参数,原因是S和L需要通过Connector进行建立流式关系
		if mode == common.S || mode == common.L {
			if option == nil {
				log.Logger().Error("Funcion S/L need option->Cid\n")
				return nil
			} else if option.CName == "" {
				log.Logger().Error("Funcion S/L need option->Cid\n")
				return nil
			}
		}
	*/

	if option != nil {
		config.Option = *option
	}

	return config
}

func (fConf *KisFuncConfig) AddConnConfig(cConf *KisConnConfig) error {
	if cConf == nil {
		return errors.New("KisConnConfig is nil")
	}

	// Function需要和Connector进行关联
	fConf.connConf = cConf

	// Connector需要和Function进行关联
	_ = cConf.WithFunc(fConf)

	// 更新Function配置中的CName
	fConf.Option.CName = cConf.CName

	return nil
}

func (fConf *KisFuncConfig) GetConnConfig() (*KisConnConfig, error) {
	if fConf.connConf == nil {
		return nil, errors.New("KisFuncConfig.connConf not set")
	}

	return fConf.connConf, nil
}

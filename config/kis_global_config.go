package config

type KisGlobalConfig struct {
	//kistype Global为kisflow的全局配置
	KisType string `yaml:"kistype"`
	//是否启动prometheus监控
	EnableProm bool `yaml:"prometheus_enable"`
	//是否需要kisflow单独启动端口监听
	PrometheusListen bool `yaml:"prometheus_listen"`
	//prometheus取点监听地址
	PrometheusServe string `yaml:"prometheus_serve"`
}

// GlobalConfig 默认全局配置，全部均为关闭
var GlobalConfig = new(KisGlobalConfig)

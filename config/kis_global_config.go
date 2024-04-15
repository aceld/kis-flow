package config

// KisGlobalConfig represents the global configuration for KisFlow
type KisGlobalConfig struct {
	// KisType Global is the global configuration for kisflow
	KisType string `yaml:"kistype"`
	// EnableProm indicates whether to start Prometheus monitoring
	EnableProm bool `yaml:"prometheus_enable"`
	// PrometheusListen indicates whether kisflow needs to start a separate port for listening
	PrometheusListen bool `yaml:"prometheus_listen"`
	// PrometheusServe is the address for Prometheus scraping
	PrometheusServe string `yaml:"prometheus_serve"`
}

// GlobalConfig is the default global configuration, all are set to off
var GlobalConfig = new(KisGlobalConfig)

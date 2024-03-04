package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"kis-flow/common"
	"kis-flow/config"
	"kis-flow/log"
	"net/http"
)

// kisMetrics kisFlow的Prometheus监控指标
type kisMetrics struct {
	//数据数量总量
	DataTotal prometheus.Counter
}

var Metrics *kisMetrics

// RunMetricsService 启动Prometheus监控服务
func RunMetricsService(serverAddr string) error {

	// 注册Prometheus 监控路由路径
	http.Handle(common.METRICS_ROUTE, promhttp.Handler())

	// 启动HttpServer
	err := http.ListenAndServe(serverAddr, nil) //多个进程不可监听同一个端口
	if err != nil {
		log.Logger().ErrorF("RunMetricsService err = %s\n", err)
	}

	return err
}

func InitMetrics() {
	Metrics = new(kisMetrics)

	// DataTotal初始化Counter
	Metrics.DataTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: common.COUNTER_KISFLOW_DATA_TOTAL_NAME,
		Help: common.COUNTER_KISFLOW_DATA_TOTAL_HELP,
	})

	// 注册Metrics
	prometheus.MustRegister(Metrics.DataTotal)
}

// RunMetrics 启动Prometheus指标服务
func RunMetrics() {
	// 初始化Prometheus指标
	InitMetrics()

	if config.GlobalConfig.EnableProm == true && config.GlobalConfig.PrometheusListen == true {
		// 启动Prometheus指标Metrics服务
		go RunMetricsService(config.GlobalConfig.PrometheusServe)
	}
}

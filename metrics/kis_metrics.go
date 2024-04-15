package metrics

import (
	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/config"
	"github.com/aceld/kis-flow/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

// kisMetrics kisFlow的Prometheus监控指标
type kisMetrics struct {
	// 数据数量总量
	DataTotal prometheus.Counter
	// 各Flow处理数据总量
	FlowDataTotal *prometheus.GaugeVec
	// Flow被调度次数
	FlowScheduleCntsToTal *prometheus.GaugeVec
	// Function被调度次数
	FuncScheduleCntsTotal *prometheus.GaugeVec
	// Function执行时间
	FunctionDuration *prometheus.HistogramVec
	// Flow执行时间
	FlowDuration *prometheus.HistogramVec
}

var Metrics *kisMetrics

// RunMetricsService 启动Prometheus监控服务
func RunMetricsService(serverAddr string) error {

	// 注册Prometheus 监控路由路径
	http.Handle(common.METRICS_ROUTE, promhttp.Handler())

	// 启动HttpServer
	err := http.ListenAndServe(serverAddr, nil) // 多个进程不可监听同一个端口
	if err != nil {
		log.Logger().Error("RunMetricsService", "err", err)
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

	// FlowDataTotal初始化GaugeVec
	Metrics.FlowDataTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: common.GANGE_FLOW_DATA_TOTAL_NAME,
			Help: common.GANGE_FLOW_DATA_TOTAL_HELP,
		},
		// 标签名称
		[]string{common.LABEL_FLOW_NAME},
	)

	// FlowScheduleCntsToTal初始化GaugeVec
	Metrics.FlowScheduleCntsToTal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: common.GANGE_FLOW_SCHE_CNTS_NAME,
			Help: common.GANGE_FLOW_SCHE_CNTS_HELP,
		},
		// 标签名称
		[]string{common.LABEL_FLOW_NAME},
	)

	// FuncScheduleCntsTotal初始化GaugeVec
	Metrics.FuncScheduleCntsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: common.GANGE_FUNC_SCHE_CNTS_NAME,
			Help: common.GANGE_FUNC_SCHE_CNTS_HELP,
		},
		// 标签名称
		[]string{common.LABEL_FUNCTION_NAME, common.LABEL_FUNCTION_MODE},
	)

	// FunctionDuration初始化HistogramVec
	Metrics.FunctionDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    common.HISTOGRAM_FUNCTION_DURATION_NAME,
		Help:    common.HISTOGRAM_FUNCTION_DURATION_HELP,
		Buckets: []float64{0.005, 0.01, 0.03, 0.08, 0.1, 0.5, 1.0, 5.0, 10, 100, 1000, 5000, 30000}, // 单位ms,最大半分钟
	},
		[]string{common.LABEL_FUNCTION_NAME, common.LABEL_FUNCTION_MODE},
	)

	// FlowDuration初始化HistogramVec
	Metrics.FlowDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    common.HISTOGRAM_FLOW_DURATION_NAME,
			Help:    common.HISTOGRAM_FLOW_DURATION_HELP,
			Buckets: []float64{0.005, 0.01, 0.03, 0.08, 0.1, 0.5, 1.0, 5.0, 10, 100, 1000, 5000, 30000, 60000}, // 单位ms,最大1分钟
		},
		[]string{common.LABEL_FLOW_NAME},
	)

	// 注册Metrics
	prometheus.MustRegister(Metrics.DataTotal)
	prometheus.MustRegister(Metrics.FlowDataTotal)
	prometheus.MustRegister(Metrics.FlowScheduleCntsToTal)
	prometheus.MustRegister(Metrics.FuncScheduleCntsTotal)
	prometheus.MustRegister(Metrics.FunctionDuration)
	prometheus.MustRegister(Metrics.FlowDuration)
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

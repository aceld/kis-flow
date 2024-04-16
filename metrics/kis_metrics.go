package metrics

import (
	"net/http"

	"github.com/aceld/kis-flow/common"
	"github.com/aceld/kis-flow/config"
	"github.com/aceld/kis-flow/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// KisMetrics kisFlow's Prometheus monitoring metrics
type KisMetrics struct {
	// Total data quantity
	DataTotal prometheus.Counter
	// Total data processed by each Flow
	FlowDataTotal *prometheus.GaugeVec
	// Flow scheduling counts
	FlowScheduleCntsToTal *prometheus.GaugeVec
	// Function scheduling counts
	FuncScheduleCntsTotal *prometheus.GaugeVec
	// Function execution time
	FunctionDuration *prometheus.HistogramVec
	// Flow execution time
	FlowDuration *prometheus.HistogramVec
}

var Metrics *KisMetrics

// RunMetricsService starts the Prometheus monitoring service
func RunMetricsService(serverAddr string) error {

	// Register Prometheus monitoring route path
	http.Handle(common.MetricsRoute, promhttp.Handler())

	// Start HttpServer
	err := http.ListenAndServe(serverAddr, nil) // Multiple processes cannot listen on the same port
	if err != nil {
		log.Logger().ErrorF("RunMetricsService err = %s\n", err)
	}

	return err
}

func InitMetrics() {
	Metrics = new(KisMetrics)

	// Initialize DataTotal Counter
	Metrics.DataTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: common.CounterKisflowDataTotalName,
		Help: common.CounterKisflowDataTotalHelp,
	})

	// Initialize FlowDataTotal GaugeVec
	Metrics.FlowDataTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: common.GamgeFlowDataTotalName,
			Help: common.GamgeFlowDataTotalHelp,
		},
		// Label names
		[]string{common.LabelFlowName},
	)

	// Initialize FlowScheduleCntsToTal GaugeVec
	Metrics.FlowScheduleCntsToTal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: common.GangeFlowScheCntsName,
			Help: common.GangeFlowScheCntsHelp,
		},
		// Label names
		[]string{common.LabelFlowName},
	)

	// Initialize FuncScheduleCntsTotal GaugeVec
	Metrics.FuncScheduleCntsTotal = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: common.GangeFuncScheCntsName,
			Help: common.GangeFuncScheCntsHelp,
		},
		// Label names
		[]string{common.LabelFunctionName, common.LabelFunctionMode},
	)

	// Initialize FunctionDuration HistogramVec
	Metrics.FunctionDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    common.HistogramFunctionDurationName,
		Help:    common.HistogramFunctionDurationHelp,
		Buckets: []float64{0.005, 0.01, 0.03, 0.08, 0.1, 0.5, 1.0, 5.0, 10, 100, 1000, 5000, 30000}, // Unit: ms, maximum half a minute
	},
		[]string{common.LabelFunctionName, common.LabelFunctionMode},
	)

	// Initialize FlowDuration HistogramVec
	Metrics.FlowDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    common.HistogramFlowDurationName,
			Help:    common.HistogramFlowDurationHelp,
			Buckets: []float64{0.005, 0.01, 0.03, 0.08, 0.1, 0.5, 1.0, 5.0, 10, 100, 1000, 5000, 30000, 60000}, // Unit: ms, maximum 1 minute
		},
		[]string{common.LabelFlowName},
	)

	// Register Metrics
	prometheus.MustRegister(Metrics.DataTotal)
	prometheus.MustRegister(Metrics.FlowDataTotal)
	prometheus.MustRegister(Metrics.FlowScheduleCntsToTal)
	prometheus.MustRegister(Metrics.FuncScheduleCntsTotal)
	prometheus.MustRegister(Metrics.FunctionDuration)
	prometheus.MustRegister(Metrics.FlowDuration)
}

// RunMetrics starts the Prometheus metrics service
func RunMetrics() {
	// Initialize Prometheus metrics
	InitMetrics()

	if config.GlobalConfig.EnableProm == true && config.GlobalConfig.PrometheusListen == true {
		// Start Prometheus metrics service
		go RunMetricsService(config.GlobalConfig.PrometheusServe)
	}
}

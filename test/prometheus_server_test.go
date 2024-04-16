package test

import (
	"testing"

	"github.com/aceld/kis-flow/metrics"
)

func TestPrometheusServer(t *testing.T) {

	err := metrics.RunMetricsService("0.0.0.0:20004")
	if err != nil {
		panic(err)
	}
}

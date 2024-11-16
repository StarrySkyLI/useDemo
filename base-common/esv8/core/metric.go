package core

import (
	"github.com/zeromicro/go-zero/core/metric"
)

const esNamespace = "es_client"

var (
	metricClientReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
		Namespace: esNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "esv8 client requests duration(ms).",
		Labels:    []string{"index"},
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 1000},
	})

	metricClientReqErrTotal = metric.NewCounterVec(&metric.CounterVecOpts{
		Namespace: esNamespace,
		Subsystem: "requests",
		Name:      "error_total",
		Help:      "esv8 client requests error count.",
		Labels:    []string{"index", "is_error"},
	})
)

package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go_demo/prom/timeu"
	"strconv"
	"time"
)

const serverNamespace = "http_server"

var (
	//metricServerReqDur = metric.NewHistogramVec(&metric.HistogramVecOpts{
	//	Namespace: serverNamespace,
	//	Subsystem: "requests",
	//	Name:      "duration_ms",
	//	Help:      "http server requests duration(ms).",
	//	Labels:    []string{"path", "method", "code"},
	//	Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 750, 1000},
	//})
	//
	//metricServerReqCodeTotal = metric.NewCounterVec(&metric.CounterVecOpts{
	//	Namespace: serverNamespace,
	//	Subsystem: "requests",
	//	Name:      "code_total",
	//	Help:      "http server requests error count.",
	//	Labels:    []string{"path", "method", "code"},
	//})
	vec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: serverNamespace,
		Subsystem: "requests",
		Name:      "duration_ms",
		Help:      "http server requests duration(ms).",
		Buckets:   []float64{5, 10, 25, 50, 100, 250, 500, 750, 1000},
	}, []string{"path", "method", "code"})
	mycounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "gin_api_count",
		},
		[]string{"method", "path"},
	)
)

func main() {

	prometheus.MustRegister(mycounter, vec)
	r := gin.New()
	r.Use(func(c *gin.Context) {

		startTime := timeu.Now()
		code := strconv.Itoa(c.Writer.Status())
		defer func() {
			mycounter.With(prometheus.Labels{
				"method": c.Request.Method,
				"path":   c.Request.RequestURI,
			}).Inc()
			vec.With(prometheus.Labels{
				"path":   c.Request.RequestURI,
				"method": c.Request.Method,
				"code":   code,
			}).Observe(float64(timeu.Since(startTime).Milliseconds()))
			fmt.Println("******************************************************")
			fmt.Println("start", startTime)
			fmt.Println("timeu.Since(startTime)", timeu.Since(startTime))
			fmt.Println("float64(timeu.Since(startTime).Milliseconds())", float64(timeu.Since(startTime).Milliseconds()))
		}()

		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		time.Sleep(10 * time.Millisecond)
		c.JSON(200, gin.H{"message": "index"})
	})
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	r.Run(":9001")
}

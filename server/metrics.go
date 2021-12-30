package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	metricsNamespace = "mattermost_video_preview"
)

//Metrics is awesome
type Metrics struct {
	registry    *prometheus.Registry
	CPUUsage    *prometheus.GaugeVec
	MemoryUsage *prometheus.GaugeVec
}

//NewMetrics is awesome
func NewMetrics() *Metrics {
	var m Metrics
	m.registry = prometheus.NewRegistry()

	m.registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
		Namespace: metricsNamespace,
	}))
	m.registry.MustRegister(collectors.NewGoCollector())

	m.CPUUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "cpu_usage_percent",
		Help:      "how many percent of the CPU time this process uses",
	},
		[]string{"cpuUsage"},
	)
	m.registry.MustRegister(m.CPUUsage)

	m.MemoryUsage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "cpu_MemoryUsage_percent",
		Help:      "how many percent of the Memory this process is using",
	},
		[]string{"memoryUsage"},
	)
	m.registry.MustRegister(m.MemoryUsage)
	// fmt.Println("This is new message")

	return &m
}

//Handler returns some value
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

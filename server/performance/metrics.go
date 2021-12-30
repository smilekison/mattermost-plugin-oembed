package performance

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
	registry             *prometheus.Registry
	CPUUsage             *prometheus.GaugeOpts
	WebSocketConnections *prometheus.GaugeVec
}

//NewMetrics is awesome
func NewMetrics() *Metrics {
	var m Metrics
	m.registry = prometheus.NewRegistry()

	m.registry.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{
		Namespace: metricsNamespace,
	}))
	m.registry.MustRegister(collectors.NewGoCollector())

	m.WebSocketConnections = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: metricsNamespace,
		Name:      "connections_total",
		Help:      "The total number requests.",
	},
		[]string{"type"},
	)
	m.registry.MustRegister(m.WebSocketConnections)

	return &m
}

//Handler returns some value
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	// "github.com/mattermost/mattermost-plugin-calls/server/performance"
	"github.com/mattermost/mattermost-server/v5/plugin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/shirou/gopsutil/process"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration
	metrics       *Metrics
}

func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.URL.Path, "/metrics") {
		p.metrics.Handler().ServeHTTP(w, r)
		return
	}
	urlBytes, _ := ioutil.ReadAll(r.Body)

	pid := int32(os.Getpid())
	process := process.Process{
		Pid: pid,
	}

	fmt.Println("This is pid: ", pid)

	cpuUsage32, _ := process.CPUPercent()
	cpuUsage := float64(cpuUsage32)
	fmt.Println("cpuUsage: ", cpuUsage)

	memoryUsage32, _ := process.MemoryPercent()
	memoryUsage := float64(memoryUsage32)
	fmt.Println("memoryUsage: ", memoryUsage)

	url := string(urlBytes)

	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// p.metrics.CPUUsage.With(prometheus.Labels{}).Inc()
	// p.metrics.MemoryUsage.With(prometheus.Labels{"memoryUsage": fmt.Sprintf("Inc %f", memoryUsage)}).Add(memoryUsage)
	// p.metrics.CPUUsage.With(prometheus.Labels{"cpuUsage": fmt.Sprintf("Inc %f", cpuUsage)}).Add(cpuUsage)
	resp, err := http.Get(url)
	p.metrics.MemoryUsage.With(prometheus.Labels{"memoryUsage": fmt.Sprintf("Inc %f", memoryUsage)}).Add(memoryUsage)
	p.metrics.CPUUsage.With(prometheus.Labels{"cpuUsage": fmt.Sprintf("Inc %f", cpuUsage)}).Add(cpuUsage)

	// p.metrics.CPUUsage.With(prometheus.Labels{"cpuUsage": fmt.Sprintf("%f", cpuUsage)}).Inc()
	// p.metrics.MemoryUsage.With(prometheus.Labels{"memoryUsage": fmt.Sprintf("Dec %f", memoryUsage)}).Dec()

	// p.metrics.CPUUsage.With(prometheus.Labels{"cpuUsage": fmt.Sprintf("This is cpu Usage Inc: %f", cpuUsage)}).Dec()

	// p.metrics.CPUUsage.With(prometheus.Labels{"cpuUsage": fmt.Sprintf("%f", cpuUsage)}).Add(float64(cpuUsage32))

	if err != nil {
		fmt.Fprintf(w, "Error getting response from %s\n%s", url, err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(w, "Error reading response from %s\n%s", url, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)

}

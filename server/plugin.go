package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"

	// "github.com/mattermost/mattermost-plugin-calls/server/performance"
	"github.com/mattermost/mattermost-server/v5/plugin"
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
	url := string(urlBytes)

	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	// p.metrics.CPUUsage.With(prometheus.Labels{"type": "KVGet"}).Inc()
	// storedData, appErr := p.API.KVGet(key)
	resp, err := http.Get(url)

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

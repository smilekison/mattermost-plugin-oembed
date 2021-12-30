package main

import (
	"sync"

	"./performance"
	"github.com/mattermost/mattermost-server/v5/plugin"
)

func main() {
	plugin.ClientMain(&Plugin{
		MattermostPlugin:  plugin.MattermostPlugin{},
		configurationLock: sync.RWMutex{},
		configuration:     &configuration{},
		metrics:           performance.NewMetrics(),
	})
}

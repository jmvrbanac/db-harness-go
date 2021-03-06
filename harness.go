package harness

import (
	"github.com/jmvrbanac/db-harness-go/mysql"
	"github.com/jmvrbanac/db-harness-go/redis"
	"github.com/jmvrbanac/db-harness-go/utils"
)

// DB Plugin selectors
const (
	Redis = "redis"
	MySQL = "mysql"
)

// DatabaseHarness is the abstract for users to interact with
type DatabaseHarness struct {
	Type    string
	Options map[string]string
	plugins map[string]utils.Plugin
}

// New creates a new DatabaseHarness
func New(dbType string, options map[string]string) *DatabaseHarness {
	harness := DatabaseHarness{
		Type:    dbType,
		Options: options,
		plugins: map[string]utils.Plugin{
			Redis: redis.New(),
			MySQL: mysql.New(),
		},
	}
	return &harness
}

// Start initializes and starts the harness
func (h *DatabaseHarness) Start() {
	p := h.GetPlugin()
	p.Initialize(h.Options)
	p.Start()

}

// Stop signals to the harness to shutdown
func (h *DatabaseHarness) Stop() {
	p := h.GetPlugin()
	p.Stop()
	p.Cleanup()
}

// GetInfo retieves the Database values provided by the active plugin
func (h *DatabaseHarness) GetInfo() utils.DatabaseInfo {
	return h.GetPlugin().GetInfo()
}

// GetPlugin returns the active Plugin
func (h *DatabaseHarness) GetPlugin() utils.Plugin {
	return h.plugins[h.Type]
}

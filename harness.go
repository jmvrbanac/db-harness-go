package harness

import "github.com/jmvrbanac/db-harness-go/redis"

const (
	// Redis is a DB Plugin Type
	Redis = "redis"
)

// Plugin is ...
type Plugin interface {
	Initialize(map[string]string)
	Start()
	Stop()
	Cleanup()
}

// DatabaseHarness is the abstract for users to interact with
type DatabaseHarness struct {
	Type    string
	Options map[string]string
	plugins map[string]Plugin
}

// New creates a new DatabaseHarness
func New(dbType string, options map[string]string) *DatabaseHarness {
	harness := DatabaseHarness{
		Type:    dbType,
		Options: options,
		plugins: map[string]Plugin{
			Redis: redis.New(),
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

// GetPlugin returns the active Plugin
func (h *DatabaseHarness) GetPlugin() Plugin {
	return h.plugins[h.Type]
}

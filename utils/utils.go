package utils

// DsnStringBuilder is a Plugin defined function to handle building a connection string
type DsnStringBuilder func() string

// Dsn is the connection string object consumed by users
type Dsn struct {
	Host       string
	Port       int64
	Proto      string
	ConnectURI DsnStringBuilder
}

// Plugin is the core interface for Databases
type Plugin interface {
	Initialize(map[string]string)
	Start()
	Stop()
	Cleanup()
	GetDsn() Dsn
}

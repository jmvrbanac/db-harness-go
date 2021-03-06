package redis

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/jmvrbanac/db-harness-go/utils"
)

// Redis is the plugin data type
type Redis struct {
	cmd *exec.Cmd
	cfg map[string]string
}

// New creates a new Redis Plugin
func New() *Redis {
	r := &Redis{
		cfg: map[string]string{
			"bind":        "0.0.0.0",
			"port":        "6379",
			"dir":         "/tmp/go-harness/data",
			"dbfilename":  "test.rdb",
			"maxclients":  "100",
			"maxmemory":   "100000000",
			"tcp-backlog": "69",
		},
	}
	return r
}

func (r *Redis) getConfigPath() string {
	return path.Join(r.cfg["dir"], "redis.conf")
}

func (r *Redis) createConfig() {
	os.MkdirAll(r.cfg["dir"], os.ModePerm)

	f, err := os.Create(r.getConfigPath())
	if err != nil {
		log.Fatal("Couldn't create Redis Config")
	}

	// Write out config
	for k, v := range r.cfg {
		f.WriteString(k + " " + v + "\n")
	}

	f.Close()
}

// Initialize prepares the harness to be executed
func (r *Redis) Initialize(options map[string]string) {
	// Update Default Config Options
	if options != nil {
		for k, v := range options {
			r.cfg[k] = v
		}
	}

	r.createConfig()
	r.cmd = exec.Command("redis-server", r.getConfigPath())
}

// Start executes the Harness
func (r *Redis) Start() {
	stdout, _ := r.cmd.StdoutPipe()
	err := r.cmd.Start()
	if err != nil {
		log.Fatal("Couldn't start harness ", err)
	}

	// Block until we know the server is up
	reader := bufio.NewReader(stdout)
	var line []byte

	for true {
		line, _, _ = reader.ReadLine()
		if strings.Contains(string(line), "The server is now ready") {
			break
		}
	}
}

// Stop signals to the process to shutdown
func (r *Redis) Stop() {
	if r.cmd.Process == nil {
		return
	}

	r.cmd.Process.Signal(os.Interrupt)
	r.cmd.Wait()
}

// Cleanup removes any temporary files that may have been created for the harness
func (r *Redis) Cleanup() {
	os.RemoveAll("/tmp/go-harness")
}

// GetInfo returns a new DatabaseInfo from current configuration
func (r *Redis) GetInfo() utils.DatabaseInfo {
	port, _ := strconv.ParseInt(r.cfg["port"], 10, 64)
	host := r.cfg["bind"]

	d := utils.DatabaseInfo{
		Host:     host,
		Port:     port,
		Proto:    "tcp",
		Database: "0",
		ConnectString: func() string {
			return fmt.Sprintf(
				"redis://%s:%d",
				host,
				port,
			)
		},
	}
	return d
}

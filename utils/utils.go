package utils

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

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

// FindFile looks up a file across the search paths and returns the first path found
func FindFile(filename string, searchPaths []string) string {
	for _, searchPath := range searchPaths {
		files, err := ioutil.ReadDir(searchPath)
		if err != nil {
			continue
		}

		for _, file := range files {
			if file.Name() == filename {
				return path.Join(searchPath, file.Name())
			}
		}
	}

	return ""
}

// WaitForFile is a blocking call to halt execution until a file exists
func WaitForFile(path string) {
	_, err := os.Stat(path)

	for err != nil {
		time.Sleep(100 * time.Millisecond)
		_, err = os.Stat(path)
	}
}

// CheckFatal cleans up harness and logs out a Fatal message
func CheckFatal(p Plugin, e error, msg string) {
	if e != nil {
		p.Stop()
		p.Cleanup()

		log.Fatal(msg)
	}
}

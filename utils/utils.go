package utils

import (
	"io/ioutil"
	"log"
	"path"
)

// ConnStringBuilder is a Plugin defined function to handle building a connection string
type ConnStringBuilder func() string

// DatabaseInfo is the connection string object consumed by users
type DatabaseInfo struct {
	User          string
	Password      string
	Host          string
	Port          int64
	Proto         string
	Database      string
	ConnectString ConnStringBuilder
}

// Plugin is the core interface for Databases
type Plugin interface {
	Initialize(map[string]string)
	Start()
	Stop()
	Cleanup()
	GetInfo() DatabaseInfo
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

// CheckFatal cleans up harness and logs out a Fatal message
func CheckFatal(p Plugin, e error, msg string) {
	if e != nil {
		p.Stop()
		p.Cleanup()

		log.Fatal(msg)
	}
}

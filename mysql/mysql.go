package mysql

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"

	"database/sql"

	// SQL Hooks
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmvrbanac/db-harness-go/utils"
	"github.com/renstrom/dedent"
)

// MySQL is the plugin data type
type MySQL struct {
	cmd         *exec.Cmd
	cfg         map[string]string
	installExec string
	mysqldExec  string
}

// New creates a new MySQL Plugin
func New() *MySQL {
	m := &MySQL{
		cfg: map[string]string{
			"host":     "0.0.0.0",
			"port":     "3306",
			"username": "tester",
			"password": "changeMe",
			"dir":      "/tmp/go-harness",
			"database": "test",
		},
	}

	return m
}

func (m *MySQL) resolvePaths() {
	searchPaths := []string{
		"/usr/bin",
		"/usr/local/bin",
		"/usr/sbin",
		"/sbin",
		"/bin",
	}

	m.installExec = utils.FindFile("mysql_install_db", searchPaths)
	m.mysqldExec = utils.FindFile("mysqld", searchPaths)

	if m.installExec == "" {
		log.Fatal("Couldn't find mysql_install_db. Check your installation")
	}

	if m.mysqldExec == "" {
		log.Fatal("Couldn't find mysqld. Check your installation")
	}
}

func (m *MySQL) createConfig(base string, tmpdir string, datadir string, confdir string) {
	socket := path.Join(tmpdir, "mysql.sock")
	pid := path.Join(tmpdir, "mysql.pid")

	os.MkdirAll(base, os.ModePerm)
	os.MkdirAll(tmpdir, os.ModePerm)
	os.MkdirAll(datadir, os.ModePerm)
	os.MkdirAll(confdir, os.ModePerm)

	cf := fmt.Sprintf(
		dedent.Dedent(`
			[mysqld]
			port=%s
			datadir=%s
			tmpdir=%s
			socket=%s
			pid-file=%s
		`),
		m.cfg["port"],
		datadir,
		tmpdir,
		socket,
		pid,
	)

	f, err := os.Create(path.Join(confdir, "my.cf"))
	if err != nil {
		log.Fatal("couldn't create MySQL config")
	}

	f.Write([]byte(cf))
	f.Close()
}

// Initialize sets up the harness and DB to be started
func (m *MySQL) Initialize(options map[string]string) {
	m.resolvePaths()

	base := m.cfg["dir"]
	tmpdir := path.Join(base, "tmp")
	datadir := path.Join(base, "var")
	confdir := path.Join(base, "etc")

	m.createConfig(base, tmpdir, datadir, confdir)

	// Initialize DB and Data directories
	initCmd := exec.Command(
		m.installExec,
		fmt.Sprintf("--defaults-file=%s", path.Join(confdir, "my.cf")),
		fmt.Sprintf("--datadir=%s", datadir),
	)

	_, err := initCmd.Output()
	if err != nil {
		log.Fatalf("Couldn't initialize MySQL database: %s", err)
	}

	// Setup DB Command
	m.cmd = exec.Command(
		m.mysqldExec,
		fmt.Sprintf("--defaults-file=%s", path.Join(confdir, "my.cf")),
		fmt.Sprintf("--datadir=%s", datadir),
	)
}

// Start executes the Harness
func (m *MySQL) Start() {
	err := m.cmd.Start()
	if err != nil {
		log.Fatal("Couldn't start harness ", err)
	}

	utils.WaitForFile(path.Join(m.cfg["dir"], "tmp", "mysql.pid"))

	// Create Testing Database and User
	db, err := sql.Open("mysql", "root@tcp(0.0.0.0:3306)/")
	utils.CheckFatal(m, err, "Couldn't create create db connection")

	_, err = db.Exec(fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS %s",
		m.cfg["database"],
	))
	utils.CheckFatal(m, err, "Couldn't create create testing db")

	_, err = db.Exec(fmt.Sprintf(
		"CREATE USER '%s'@'localhost' IDENTIFIED BY '%s'",
		m.cfg["username"],
		m.cfg["password"],
	))
	utils.CheckFatal(m, err, "Couldn't create create testing user")

	db.Close()
}

// Stop signals to the process to shutdown
func (m *MySQL) Stop() {
	if m.cmd.Process == nil {
		return
	}

	m.cmd.Process.Kill()
	m.cmd.Wait()
}

// Cleanup removes any temporary files that may have been created for the harness
func (m *MySQL) Cleanup() {
	os.RemoveAll("/tmp/go-harness")
}

// GetInfo returns a new DatabaseInfo from current configuration
func (m *MySQL) GetInfo() utils.DatabaseInfo {
	port, _ := strconv.ParseInt(m.cfg["port"], 10, 64)

	d := utils.DatabaseInfo{
		Host:     m.cfg["host"],
		Port:     port,
		Proto:    "tcp",
		Database: m.cfg["database"],
		ConnectURI: func() string {
			return fmt.Sprintf(
				"%s:%s@tcp(%s:%s)/%s",
				m.cfg["username"],
				m.cfg["password"],
				m.cfg["host"],
				m.cfg["port"],
				m.cfg["database"],
			)
		},
	}
	return d
}

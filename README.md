# Go Database Harness

[![Build Status](https://travis-ci.org/jmvrbanac/db-harness-go.svg?branch=master)](https://travis-ci.org/jmvrbanac/db-harness-go) [![codecov](https://codecov.io/gh/jmvrbanac/db-harness-go/branch/master/graph/badge.svg)](https://codecov.io/gh/jmvrbanac/db-harness-go)[![GoDoc](https://godoc.org/github.com/jmvrbanac/db-harness-go?status.svg)](https://godoc.org/github.com/jmvrbanac/db-harness-go)

This library is inspired by the `testing.redis` and `testing.mysql` packages
in Python. It is designed to act as a testing harness for Database servers
such as Redis, MySQL, Postgres, etc. The primary purpose of a database test
harness is that it allows for users to build a temporary database to run
their tests against; which is useful for testing migrations or reducing mocking.

## Plugin(s)

* Redis
* MySQL
* Postgres (Not completed)
* Rethink (Not completed)

**Note:** Make sure to stop the harness after your testing (including
errors). Failure to do so will result in a orphaned server instance that you'll
have to manually kill.


## Example

**Startup with Defaults**
```Go
import "github.com/jmvrbanac/db-harness-go"

// Build a new Harness
h := harness.New(harness.Redis, nil)

// Initialize and start the harness
h.Start()

// ... Do your testing here

// Shutdown and cleanup the harness
h.Stop()
```

**Set a different port**
```Go
import "github.com/jmvrbanac/db-harness-go"

// Setup Options
options := map[string]string {
    "port": "2222",
}

// Build a new Harness
h := harness.New(harness.Redis, options)

// Initialize and start the harness
h.Start()

// ... Do your testing here

// Shutdown and cleanup the harness
h.Stop()
```

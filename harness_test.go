package harness_test

import (
	"net"

	"github.com/jmvrbanac/db-harness-go"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func canConnect(addr string) bool {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return false
	}

	conn.Close()
	return true
}

var _ = Describe("Redis Harness", func() {
	It("Should be able to run without options", func() {
		h := harness.New(harness.Redis, nil)

		h.Start()
		Expect(canConnect("0.0.0.0:6379")).To(BeTrue())

		h.Stop()
		Expect(canConnect("0.0.0.0:6379")).To(BeFalse())
	})

	It("Should be able to accept a port option", func() {
		options := map[string]string{
			"port": "2222",
		}
		h := harness.New(harness.Redis, options)

		h.Start()
		Expect(canConnect("0.0.0.0:2222")).To(BeTrue())

		h.Stop()
		Expect(canConnect("0.0.0.0:2222")).To(BeFalse())
	})
})

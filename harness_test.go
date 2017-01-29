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

var _ = Describe("DB Harnesses", func() {
	// Describe("MySQL", func() {
	// 	It("Should be able to run without options", func() {
	// 		h := harness.New(harness.MySQL, nil)
	// 		h.Start()

	// 		Expect(canConnect("0.0.0.0:3306")).To(BeTrue())

	// 		info := h.GetInfo()
	// 		Expect(info.ConnectURI()).To(Equal("tester:changeMe@tcp(0.0.0.0:3306)/test"))

	// 		h.Stop()
	// 	})
	// })

	Describe("Redis", func() {
		It("Should be able to run without options", func() {
			h := harness.New(harness.Redis, nil)

			h.Start()
			Expect(h.GetInfo().ConnectURI()).To(Equal("redis://0.0.0.0:6379"))
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
})

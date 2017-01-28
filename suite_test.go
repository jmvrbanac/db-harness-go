package harness_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestDbHarness(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "DB Harness Suite")
}

package transmissionrss_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestTrss(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Trss Suite")
}

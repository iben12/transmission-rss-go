package transmissionrss_test

import (
	"io/ioutil"
	"testing"

	trss "github.com/iben12/transmission-rss-go/trss"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
)

func TestTrss(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Trss Suite")
}

var _ = BeforeSuite(func() {
	trss.Logger = zerolog.New(ioutil.Discard).With().Timestamp().Logger()
})

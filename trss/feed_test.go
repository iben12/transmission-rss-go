package transmissionrss_test

import (
	"net/http"

	helpers "github.com/iben12/transmission-rss-go/tests/helpers"
	"github.com/iben12/transmission-rss-go/tests/mocks"
	"github.com/iben12/transmission-rss-go/trss"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Feed", func() {

	It("should parse valid XML", func() {
		url := "/feed/1234"
		handler := func(rw http.ResponseWriter, req *http.Request) {
			Expect(url).To(Equal(req.URL.String()))
			rw.Write([]byte(mocks.ValidXml))
		}

		server := helpers.CreateServer(handler)

		defer server.Close()

		feed := new(transmissionrss.Feeds)

		rssAddress := server.URL + url

		items, _ := feed.FetchItems(rssAddress)

		Expect(len(items)).To(Equal(1))
		Expect(items[0].ShowTitle).To(Equal("Million Dollar Listing: Los Angeles"))
	})

	It("should error on invalid XML", func() {
		handler := func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(mocks.InvalidXml))
		}

		server := helpers.CreateServer(handler)

		defer server.Close()

		feed := new(transmissionrss.Feeds)

		_, err := feed.FetchItems(server.URL)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("cannot parse XML"))
	})
})

package transmissionrss_test

import (
	"net/http"
	"net/http/httptest"

	helpers "github.com/iben12/transmission-rss-go/tests/helpers"
	"github.com/iben12/transmission-rss-go/tests/mocks"
	"github.com/iben12/transmission-rss-go/trss"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Feed", func() {
	var server *httptest.Server

	AfterEach(func() {
		server.Close()
	})

	It("should parse valid XML", func() {
		url := "/feed/1234"
		handler := func(rw http.ResponseWriter, req *http.Request) {
			Expect(url).To(Equal(req.URL.String()))
			rw.Write([]byte(mocks.ValidXml))
		}

		server = helpers.CreateServer(handler)

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

		feed := new(transmissionrss.Feeds)

		_, err := feed.FetchItems(server.URL)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("cannot parse XML"))
	})

	It("should error if server is unavailable", func() {
		feed := new(transmissionrss.Feeds)

		_, err := feed.FetchItems("http://localhost")

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("Get \"http://localhost\": dial tcp [::1]:80: connect: connection refused"))
	})

	It("should error if request fails", func() {
		handler := func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(500)
			rw.Write([]byte("HTTP 500 error"))
		}

		server := helpers.CreateServer(handler)
		feed := new(transmissionrss.Feeds)

		_, err := feed.FetchItems(server.URL)

		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal("500 Internal Server Error"))
	})
})

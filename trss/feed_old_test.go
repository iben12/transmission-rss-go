package transmissionrss_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iben12/transmission-rss-go/tests/mocks"
	"github.com/iben12/transmission-rss-go/trss"
	"github.com/stretchr/testify/assert"
)

func TestFeedParse(t *testing.T) {
	assert := assert.New(t)

	t.Run("Valid XML", func(t *testing.T) {
		url := "/feed/1234"
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal(url, req.URL.String())
			rw.Write([]byte(mocks.ValidXml))
		}))

		defer server.Close()

		feed := new(transmissionrss.Feeds)

		rssAddress := server.URL + url

		items, _ := feed.FetchItems(rssAddress)

		expectedLength := 1

		assert.Equal(len(items), expectedLength)
		assert.Equal(items[0].ShowTitle, "Million Dollar Listing: Los Angeles")
	})

	t.Run("Invalid XML", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.Write([]byte(mocks.InvalidXml))
		}))

		defer server.Close()

		feed := new(transmissionrss.Feeds)

		_, err := feed.FetchItems(server.URL)

		if assert.Error(err) {
			assert.Equal(err.Error(), "cannot parse XML")
		}
	})

}

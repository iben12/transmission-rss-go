package transmissionrss_test

import (
	"encoding/json"
	"errors"

	// "fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	// helper "github.com/iben12/transmission-rss-go/tests/helper"
	"github.com/iben12/transmission-rss-go/tests/mocks"
	trss "github.com/iben12/transmission-rss-go/trss"
)

type ApiEndpoint func(w http.ResponseWriter, r *http.Request)

func sendApiRequest(endpoint ApiEndpoint, method string, url string) (*http.Response, []byte) {
	req := httptest.NewRequest(method, url, nil)
	w := httptest.NewRecorder()

	endpoint(w, req)

	res := w.Result()

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	return res, body
}

var _ = Describe("Api", func() {
	var api *trss.Api
	var rssAddress string
	var mockFeeds *mocks.MockFeeds
	var mockEpisodes *mocks.MockEpisodes
	var mockTrs *mocks.MockTransmissionService

	BeforeEach(func() {
		api = new(trss.Api)
		rssAddress = "https://example.rss/feed"
		os.Setenv("RSS_FEED_ADDRESS", rssAddress)
		DeferCleanup(func() {
			os.Unsetenv("RSS_FEED_ADDRESS")
		})
		mockFeeds = &mocks.MockFeeds{}
		trss.FeedService = mockFeeds
		mockEpisodes = &mocks.MockEpisodes{}
		trss.EpisodeService = mockEpisodes
		mockTrs = &mocks.MockTransmissionService{}
		trss.TrsService = mockTrs
	})

	Context("Feeds endpoint", func() {
		It("should return feeds", func() {
			_, body := sendApiRequest(api.Feeds, http.MethodGet, "/")

			var result []interface{}
			err := json.Unmarshal(body, &result)
			Expect(err).To(BeNil())

			expectedResult := []interface{}{map[string]interface{}{"name": "ShowRSS", "url": rssAddress}}

			Expect(result).To(Equal(expectedResult))
		})
	})

	Context("Episodes endpoint", func() {
		It("should return episodes", func() {
			episodes := []trss.Episode{
				mocks.EpisodeExample,
			}

			mockEpisodes.MockAll = func() ([]trss.Episode, error) {
				return episodes, nil
			}

			_, body := sendApiRequest(api.Episodes, http.MethodGet, "/")

			var resultData []trss.Episode
			jsonErr := json.Unmarshal(body, &resultData)

			Expect(jsonErr).To(BeNil())

			Expect(resultData).To(Equal(episodes))
		})

		It("shoud returns HTTP 500 if query fails", func() {
			episodes := []trss.Episode{}

			mockEpisodes.MockAll = func() ([]trss.Episode, error) {
				return episodes, errors.New("Database error")
			}

			res, body := sendApiRequest(api.Episodes, http.MethodGet, "/")

			Expect(res.StatusCode).To(Equal(500))

			var resultData map[string]string
			jsonErr := json.Unmarshal(body, &resultData)

			Expect(jsonErr).To(BeNil())
			Expect(resultData).To(HaveKeyWithValue("error", "Could not read episodes"))
		})
	})

	Context("Download endpoint", func() {
		It("should respond whith downloaded items", func() {

		})

		It("should resturn error", func() {

		})
	})
})

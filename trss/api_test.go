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
	"gorm.io/gorm"

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
			mockFeeds.MockFetchItems = func(r string) ([]trss.FeedItem, error) {
				feedItems := []trss.FeedItem{mocks.FeedItemExample}

				return feedItems, nil
			}

			findMockData := map[string]mocks.FindMockData{
				"1": {Episode: false, Err: gorm.ErrRecordNotFound},
			}

			mocks.CreateEpisodeFindMock(mockEpisodes, findMockData)

			mockTrs.MockAddTorrent = func(e trss.Episode) error {
				return nil
			}

			mockEpisodes.MockAddEpisode = func(e *trss.Episode) error {
				e.ID = 2
				return nil
			}

			res, body := sendApiRequest(api.Download, http.MethodGet, "/")

			GinkgoWriter.Println(string(body))
			Expect(res.StatusCode).To(Equal(200))

			type ResponseType struct {
				Episodes []trss.Episode
				Errors   []string
			}
			responseData := &ResponseType{}
			json.Unmarshal(body, responseData)

			expectedEpisode := mocks.EpisodeExample
			expectedEpisode.ID = 2

			Expect(responseData.Episodes).To(ContainElement(expectedEpisode))
			Expect(responseData.Errors).To(BeEmpty())
		})

		It("should return error if transmission add fails", func() {
			mockFeeds.MockFetchItems = func(r string) ([]trss.FeedItem, error) {
				feedItems := []trss.FeedItem{mocks.FeedItemExample}

				return feedItems, nil
			}

			findMockData := map[string]mocks.FindMockData{
				"1": {Episode: false, Err: gorm.ErrRecordNotFound},
			}

			mocks.CreateEpisodeFindMock(mockEpisodes, findMockData)

			mockTrs.MockAddTorrent = func(e trss.Episode) error {
				return errors.New("Transmission error")
			}

			mockEpisodes.MockAddEpisode = func(e *trss.Episode) error {
				e.ID = 2
				return nil
			}

			res, body := sendApiRequest(api.Download, http.MethodGet, "/")

			GinkgoWriter.Println(string(body))
			Expect(res.StatusCode).To(Equal(200))

			type ResponseType struct {
				Episodes []trss.Episode
				Errors   []string
			}
			responseData := &ResponseType{}
			json.Unmarshal(body, responseData)

			expectedEpisode := mocks.EpisodeExample
			expectedEpisode.ID = 2

			Expect(responseData.Episodes).To(BeEmpty())
			Expect(responseData.Errors).To(ContainElement("Transmission error"))
		})

		It("should return error if feed fetch fails", func() {
			mockFeeds.MockFetchItems = func(r string) ([]trss.FeedItem, error) {
				return nil, errors.New("Feed error")
			}

			findMockData := map[string]mocks.FindMockData{
				"1": {Episode: false, Err: gorm.ErrRecordNotFound},
			}

			mocks.CreateEpisodeFindMock(mockEpisodes, findMockData)

			mockTrs.MockAddTorrent = func(e trss.Episode) error {
				return errors.New("Transmission error")
			}

			mockEpisodes.MockAddEpisode = func(e *trss.Episode) error {
				e.ID = 2
				return nil
			}

			res, body := sendApiRequest(api.Download, http.MethodGet, "/")

			GinkgoWriter.Println(string(body))
			Expect(res.StatusCode).To(Equal(500))

			responseData := map[string]string{}
			json.Unmarshal(body, &responseData)

			expectedEpisode := mocks.EpisodeExample
			expectedEpisode.ID = 2

			Expect(responseData["episodes"]).To(BeEmpty())
			Expect(responseData["error"]).To(Equal("Could not parse feed"))
		})
	})

	Context("Cleanup endpoint", func() {
		It("returns cleaned torrent titles", func() {
			mockTrs.MockCleanFinished = func() ([]string, error) {
				return []string{"Episode title"}, nil
			}

			res, body := sendApiRequest(api.Clean, http.MethodGet, "/")

			Expect(res.StatusCode).To(Equal(200))

			responseData := []string{}

			json.Unmarshal(body, &responseData)

			Expect(responseData).To(ContainElement("Episode title"))
		})

		It("returns error if clean fails", func() {
			mockTrs.MockCleanFinished = func() ([]string, error) {
				return nil, errors.New("Remove error")
			}

			res, body := sendApiRequest(api.Clean, http.MethodGet, "/")

			Expect(res.StatusCode).To(Equal(500))

			responseData := map[string]string{}

			json.Unmarshal(body, &responseData)

			Expect(responseData["error"]).To(Equal("Could not remove torrents"))
		})
	})
})

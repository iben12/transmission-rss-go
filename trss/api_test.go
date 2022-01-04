package transmissionrss_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/iben12/transmission-rss-go/tests/mocks"
	trss "github.com/iben12/transmission-rss-go/trss"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestApi(t *testing.T) {
	assert := assert.New(t)
	rssAddress := "https://example.rss/feed"
	os.Setenv("RSS_FEED_ADDRESS", rssAddress)
	mockFeeds := &mocks.MockFeeds{}
	trss.FeedService = mockFeeds
	mockEpisodes := &mocks.MockEpisodes{}
	trss.EpisodeService = mockEpisodes
	mockTrs := &mocks.MockTransmissionService{}
	trss.TrsService = mockTrs

	t.Run("Feeds endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		api := new(trss.Api)

		api.Feeds(w, req)

		res := w.Result()

		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)

		expectedResponse := fmt.Sprintf("[{\"name\":\"ShowRSS\",\"url\":\"%v\"}]\n", rssAddress)
		assert.Equal(expectedResponse, string(data))
	})

	t.Run("Episodes endpoint", func(t *testing.T) {
		episodes := []trss.Episode{
			mocks.EpisodeExample,
		}

		mockEpisodes.MockAll = func() ([]trss.Episode, error) {
			return episodes, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		api := new(trss.Api)

		api.Episodes(w, req)

		res := w.Result()

		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)

		var resultData []trss.Episode
		jsonErr := json.Unmarshal(data, &resultData)

		if jsonErr != nil {
			t.Error(jsonErr)
		}

		assert.Equal(episodes, resultData)
	})

	t.Run("Episodes endpoint fail", func(t *testing.T) {
		episodes := []trss.Episode{}

		mockEpisodes.MockAll = func() ([]trss.Episode, error) {
			return episodes, errors.New("Database error")
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		api := new(trss.Api)

		api.Episodes(w, req)

		res := w.Result()

		assert.Equal(500, res.StatusCode)

		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)

		var resultData map[string]string
		jsonErr := json.Unmarshal(data, &resultData)

		if jsonErr != nil {
			t.Error(jsonErr)
		}

		expectedResponse := map[string]string{"error": "Could not read episodes"}

		assert.Equal(expectedResponse, resultData)
	})

	t.Run("Download endpoint", func(t *testing.T) {
		mockFeeds.MockFetchItems = func(r string) ([]trss.FeedItem, error) {
			feedItems := []trss.FeedItem{mocks.FeedItemExample}

			return feedItems, nil
		}

		mockEpisodes.MockFindEpisode = func(e *trss.Episode) (trss.Episode, error) {
			return trss.Episode{}, gorm.ErrRecordNotFound
		}

		mockEpisodes.MockDownloadEpisode = func(e trss.Episode) error {
			return nil
		}

		mockEpisodes.MockAddEpisode = func(e *trss.Episode) error {
			e.ID = 2
			return nil
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		api := new(trss.Api)

		api.Download(w, req)

		res := w.Result()

		assert.Equal(200, res.StatusCode)

		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)

		type ResponseType struct {
			Episodes []trss.Episode
			Errors   []string
		}

		expectedEpisode := mocks.EpisodeExample
		expectedEpisode.ID = 2
		expectedEpisodes := []trss.Episode{expectedEpisode}
		expectedErrors := []string(nil)

		responseData := &ResponseType{}

		json.Unmarshal(data, responseData)

		assert.Equal(expectedEpisodes, responseData.Episodes)
		assert.Equal(expectedErrors, responseData.Errors)
	})

	t.Run("Download endpoint fail", func(t *testing.T) {
		mockFeeds.MockFetchItems = func(r string) ([]trss.FeedItem, error) {
			feedItems := []trss.FeedItem{mocks.FeedItemExample}

			return feedItems, nil
		}

		mockEpisodes.MockFindEpisode = func(e *trss.Episode) (trss.Episode, error) {
			return trss.Episode{}, gorm.ErrRecordNotFound
		}

		mockEpisodes.MockDownloadEpisode = func(e trss.Episode) error {
			return errors.New("Transmission error")
		}

		mockEpisodes.MockAddEpisode = func(e *trss.Episode) error {
			e.ID = 2
			return nil
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		api := new(trss.Api)

		api.Download(w, req)

		res := w.Result()

		assert.Equal(200, res.StatusCode)

		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)

		type ResponseType struct {
			Episodes []trss.Episode
			Errors   []string
		}

		expectedEpisodes := []trss.Episode(nil)
		expectedErrors := []string{"Transmission error"}

		responseData := &ResponseType{}

		json.Unmarshal(data, responseData)

		assert.Equal(expectedEpisodes, responseData.Episodes)
		assert.Equal(expectedErrors, responseData.Errors)
	})

	t.Run("Clean endpoint", func(t *testing.T) {
		mockTrs.MockCleanFinished = func() ([]string, error) {
			return []string{"Episode title"}, nil
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		api := new(trss.Api)

		api.Clean(w, req)

		res := w.Result()

		assert.Equal(200, res.StatusCode)

		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)

		responseData := []string{}

		json.Unmarshal(data, &responseData)

		expectedResponse := []string{"Episode title"}

		assert.Equal(expectedResponse, responseData)
	})

	t.Run("Clean endpoint fail", func(t *testing.T) {
		mockTrs.MockCleanFinished = func() ([]string, error) {
			return []string{}, errors.New("Remove error")
		}

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		api := new(trss.Api)

		api.Clean(w, req)

		res := w.Result()

		assert.Equal(500, res.StatusCode)

		defer res.Body.Close()
		data, _ := ioutil.ReadAll(res.Body)

		responseData := map[string]string{}

		json.Unmarshal(data, &responseData)

		expectedResponse := map[string]string{"error": "Could not remove torrents"}

		assert.Equal(expectedResponse, responseData)
	})
}

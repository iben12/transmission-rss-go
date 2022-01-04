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
)

func TestApi(t *testing.T) {
	assert := assert.New(t)
	rssAddress := "https://example.rss/feed"
	os.Setenv("RSS_FEED_ADDRESS", rssAddress)
	mockFeeds := &mocks.MockFeeds{}
	trss.FeedService = mockFeeds
	mockEpisodes := &mocks.MockEpisodes{}
	trss.EpisodeService = mockEpisodes

	t.Run("Feeds endpoint", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		api := new(trss.Api)

		api.Feeds(w, req)

		res := w.Result()

		defer res.Body.Close()
		data, err := ioutil.ReadAll(res.Body)

		expectedResponse := fmt.Sprintf("[{\"name\":\"ShowRSS\",\"url\":\"%v\"}]\n", rssAddress)
		assert.Nil(err)
		assert.Equal(expectedResponse, string(data))
	})

	t.Run("Episodes endpoint", func(t *testing.T) {
		episodes := []trss.Episode{
			{ShowId: "1", EpisodeId: "2", Title: "From test", ShowTitle: "Testdriver"},
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
		data, err := ioutil.ReadAll(res.Body)

		var resultData []trss.Episode
		jsonErr := json.Unmarshal(data, &resultData)

		if jsonErr != nil {
			t.Error(jsonErr)
		}

		assert.Nil(err)
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
		data, err := ioutil.ReadAll(res.Body)

		var resultData map[string]string
		jsonErr := json.Unmarshal(data, &resultData)

		if jsonErr != nil {
			t.Error(jsonErr)
		}

		expectedResponse := map[string]string{"error": "Could not read episodes"}

		assert.Nil(err)
		assert.Equal(expectedResponse, resultData)
	})
}

package transmissionrss_test

import (
	"errors"
	"testing"

	"github.com/iben12/transmission-rss-go/tests/mocks"
	trss "github.com/iben12/transmission-rss-go/trss"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestDownload(t *testing.T) {
	assert := assert.New(t)
	mockEpisodes := &mocks.MockEpisodes{}

	t.Run("Episodes download", func(t *testing.T) {
		expectedEpisodeId := "22"
		feedItems := []trss.FeedItem{
			{ShowId: "1", EpisodeId: "12"},
			{ShowId: "2", EpisodeId: expectedEpisodeId},
		}

		mockEpisodes.MockFindEpisode = func(e *trss.Episode) (trss.Episode, error) {
			if e.ShowId == "1" {
				return *e, nil // episode exists
			} else {
				return trss.Episode{}, gorm.ErrRecordNotFound // episode does not exist
			}
		}

		mockEpisodes.MockDownloadEpisode = func(e trss.Episode) error {
			return nil
		}

		mockEpisodes.MockAddEpisode = func(e *trss.Episode) error {
			e.ID = 2
			return nil
		}

		downloaded, _ := trss.Download(feedItems, mockEpisodes)

		expectedLength := 1
		assert.Equal(len(downloaded), expectedLength)
		assert.Equal(downloaded[0].EpisodeId, expectedEpisodeId)
	})

	t.Run("Episodes fail", func(t *testing.T) {
		expectedEpisodeId := "22"
		feedItems := []trss.FeedItem{
			{ShowId: "1", EpisodeId: "12"},
			{ShowId: "2", EpisodeId: expectedEpisodeId},
		}

		mockEpisodes.MockFindEpisode = func(e *trss.Episode) (trss.Episode, error) {
			return trss.Episode{}, gorm.ErrRecordNotFound // episode does not exist
		}

		transmissionError := errors.New("Transmission error")

		mockEpisodes.MockDownloadEpisode = func(e trss.Episode) error {
			if e.ShowId == "1" {
				return transmissionError
			}
			return nil
		}

		dbError := errors.New("Database error")

		mockEpisodes.MockAddEpisode = func(e *trss.Episode) error {
			if e.ShowId == "2" {
				return dbError
			}
			e.ID = 2
			return nil
		}

		downloaded, errs := trss.Download(feedItems, mockEpisodes)

		expectedLength := 0
		assert.Equal(len(downloaded), expectedLength)

		expectedErrors := []string{transmissionError.Error(), dbError.Error()}

		assert.Equal(errs, expectedErrors)
	})
}
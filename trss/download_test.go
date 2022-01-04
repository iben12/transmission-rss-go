package transmissionrss_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/iben12/transmission-rss-go/tests/mocks"
	trss "github.com/iben12/transmission-rss-go/trss"
	"gorm.io/gorm"
)

func TestDownload(t *testing.T) {
	mockHandler := &mocks.MockEpisodeHandler{}

	t.Run("Episodes download", func(t *testing.T) {
		expectedEpisodeId := "22"
		feedItems := []trss.FeedItem{
			{ShowId: "1", EpisodeId: "12"},
			{ShowId: "2", EpisodeId: expectedEpisodeId},
		}

		mockHandler.MockFindEpisode = func(e *trss.Episode) (trss.Episode, error) {
			if e.ShowId == "1" {
				return *e, nil // episode exists
			} else {
				return trss.Episode{}, gorm.ErrRecordNotFound // episode does not exist
			}
		}

		mockHandler.MockDownloadEpisode = func(e trss.Episode) error {
			return nil
		}

		mockHandler.MockAddEpisode = func(e *trss.Episode) error {
			e.ID = 2
			return nil
		}

		downloaded, _ := trss.Download(feedItems, mockHandler)

		expectedLength := 1
		if len(downloaded) != expectedLength {
			t.Errorf("Expected %v items, but got %v", expectedLength, len(downloaded))
		}

		if downloaded[0].EpisodeId != expectedEpisodeId {
			t.Errorf("Expected EpisodeID to be %v, but got %v", expectedEpisodeId, downloaded[0].EpisodeId)
		}
	})

	t.Run("Episodes fail", func(t *testing.T) {
		expectedEpisodeId := "22"
		feedItems := []trss.FeedItem{
			{ShowId: "1", EpisodeId: "12"},
			{ShowId: "2", EpisodeId: expectedEpisodeId},
		}

		mockHandler.MockFindEpisode = func(e *trss.Episode) (trss.Episode, error) {
			return trss.Episode{}, gorm.ErrRecordNotFound // episode does not exist
		}

		transmissionError := errors.New("Transmission error")

		mockHandler.MockDownloadEpisode = func(e trss.Episode) error {
			if e.ShowId == "1" {
				return transmissionError
			}
			return nil
		}

		dbError := errors.New("Database error")

		mockHandler.MockAddEpisode = func(e *trss.Episode) error {
			if e.ShowId == "2" {
				return dbError
			}
			e.ID = 2
			return nil
		}

		downloaded, errs := trss.Download(feedItems, mockHandler)

		expectedLength := 0
		if len(downloaded) != expectedLength {
			t.Errorf("Expected %v items, but got %v", expectedLength, len(downloaded))
		}

		expectedErrors := []string{transmissionError.Error(), dbError.Error()}

		if !reflect.DeepEqual(expectedErrors, errs) {
			t.Errorf("Expected errors do not match actual %v %v", expectedErrors, errs)
		}
	})
}

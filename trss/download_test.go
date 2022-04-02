package transmissionrss_test

import (
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/iben12/transmission-rss-go/tests/mocks"
	trss "github.com/iben12/transmission-rss-go/trss"
	"gorm.io/gorm"
)

var _ = Describe("Download", func() {
	var mockEpisodes mocks.MockEpisodes
	var expectedEpisodeId string
	var feedItems []trss.FeedItem

	BeforeEach(func() {
		mockEpisodes = mocks.MockEpisodes{}

		expectedEpisodeId = "22"
		feedItems = []trss.FeedItem{
			{ShowId: "1", EpisodeId: "12"},
			{ShowId: "2", EpisodeId: expectedEpisodeId},
		}
	})

	It("should download episodes", func() {
		findMockData := map[string]mocks.FindMockData{
			"1": {Episode: true, Err: nil},
			"2": {Episode: false, Err: gorm.ErrRecordNotFound},
		}

		mocks.CreateFindMock(&mockEpisodes, findMockData)

		mockEpisodes.MockDownloadEpisode = func(e trss.Episode) error {
			return nil
		}

		mockEpisodes.MockAddEpisode = func(e *trss.Episode) error {
			e.ID = 2
			return nil
		}

		downloaded, _ := trss.Download(feedItems, &mockEpisodes)

		Expect(len(downloaded)).To(Equal(1))
		Expect(downloaded[0].EpisodeId).To(Equal(expectedEpisodeId))
	})

	It("should return error if download fails", func() {
		findMockData := map[string]mocks.FindMockData{
			"1": {Episode: false, Err: gorm.ErrRecordNotFound},
			"2": {Episode: false, Err: gorm.ErrRecordNotFound},
		}

		mocks.CreateFindMock(&mockEpisodes, findMockData)

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

		downloaded, errs := trss.Download(feedItems, &mockEpisodes)

		Expect(downloaded).To(HaveLen(0))
		Expect(errs).To(ContainElements(transmissionError.Error(), dbError.Error()))
	})
})

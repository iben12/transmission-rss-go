package mocks

import (
	trss "github.com/iben12/transmission-rss-go/trss"
)

type MockEpisodes struct {
	MockAddEpisode  func(e *trss.Episode) error
	MockFindEpisode func(e *trss.Episode) (trss.Episode, error)
	MockAll         func() ([]trss.Episode, error)
}

func (h *MockEpisodes) AddEpisode(e *trss.Episode) error {
	return h.MockAddEpisode(e)
}

func (h *MockEpisodes) FindEpisode(e *trss.Episode) (trss.Episode, error) {
	return h.MockFindEpisode(e)
}

func (h *MockEpisodes) All() ([]trss.Episode, error) {
	return h.MockAll()
}

type FindMockData struct {
	Episode bool
	Err     error
}

func CreateEpisodeFindMock(mockEpisodes *MockEpisodes, data map[string]FindMockData) *MockEpisodes {
	mockEpisodes.MockFindEpisode = func(e *trss.Episode) (trss.Episode, error) {
		var episode *trss.Episode
		if data[e.ShowId].Episode {
			episode = e
		} else {
			episode = &trss.Episode{}
		}
		return *episode, data[e.ShowId].Err
	}

	return mockEpisodes
}

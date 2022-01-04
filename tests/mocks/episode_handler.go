package mocks

import (
	trss "github.com/iben12/transmission-rss-go/trss"
)

type MockEpisodes struct {
	MockAddEpisode      func(e *trss.Episode) error
	MockFindEpisode     func(e *trss.Episode) (trss.Episode, error)
	MockAll             func() ([]trss.Episode, error)
	MockDownloadEpisode func(e trss.Episode) error
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

func (h *MockEpisodes) DownloadEpisode(e trss.Episode) error {
	return h.MockDownloadEpisode(e)
}

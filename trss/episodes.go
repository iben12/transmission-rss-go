package transmissionrss

import (
	"gorm.io/gorm"
)

type Episode struct {
	gorm.Model
	ShowId    string `json:"show_id"`
	EpisodeId string `json:"episode_id"`
	ShowTitle string `json:"show_title"`
	Title     string `json:"title"`
	Link      string `json:"link"`
}

type Episodes interface {
	AddEpisode(episode *Episode) error
	FindEpisode(episode *Episode) (Episode, error)
	All() ([]Episode, error)
	DownloadEpisode(episode Episode) error
}

type EpisodeHandler struct {
	db           *gorm.DB
	transmission *Trs
}

func (h *EpisodeHandler) AddEpisode(episode *Episode) error {
	result := h.db.Create(&episode)
	return result.Error
}

func (h *EpisodeHandler) FindEpisode(episodeToFind *Episode) (Episode, error) {
	episode := Episode{}
	result := h.db.Where(episodeToFind).First(&episode)

	return episode, result.Error
}

func (h *EpisodeHandler) All() ([]Episode, error) {
	episodes := []Episode{}

	result := h.db.Find(&episodes)

	return episodes, result.Error
}

func (h *EpisodeHandler) DownloadEpisode(episode Episode) error {
	return h.transmission.AddDownload(episode)
}

func NewEpisodeHanlder() Episodes {
	dbConnection := new(DB).getConnection()
	transmissionClient := &Trs{}
	return &EpisodeHandler{db: dbConnection, transmission: transmissionClient}
}

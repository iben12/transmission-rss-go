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

type EpisodeHandler interface {
	AddEpisode(episode *Episode) error
	FindEpisode(episode *Episode) (Episode, error)
	All() ([]Episode, error)
	DownloadEpisode(episode Episode) error
}

type Episodes struct {
	db           *gorm.DB
	transmission *Trs
}

func (h *Episodes) AddEpisode(episode *Episode) error {
	result := h.db.Create(&episode)
	return result.Error
}

func (h *Episodes) FindEpisode(episodeToFind *Episode) (Episode, error) {
	episode := Episode{}
	result := h.db.Where(episodeToFind).First(&episode)

	return episode, result.Error
}

func (h *Episodes) All() ([]Episode, error) {
	episodes := []Episode{}

	result := h.db.Find(&episodes)

	return episodes, result.Error
}

func (h *Episodes) DownloadEpisode(episode Episode) error {
	return h.transmission.AddDownload(episode)
}

func NewEpisodeHanlder() EpisodeHandler {
	dbConnection := new(DB).getConnection()
	transmissionClient := &Trs{}
	return &Episodes{db: dbConnection, transmission: transmissionClient}
}

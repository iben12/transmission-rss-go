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
	// DownloadEpisode(episode Episode) error
}

type Episodes struct {
	Db           *gorm.DB
	transmission TransmissionService
}

func (h *Episodes) AddEpisode(episode *Episode) error {
	result := h.Db.Create(&episode)
	return result.Error
}

func (h *Episodes) FindEpisode(episodeToFind *Episode) (Episode, error) {
	episode := Episode{}
	result := h.Db.Where(episodeToFind).First(&episode)

	return episode, result.Error
}

func (h *Episodes) All() ([]Episode, error) {
	episodes := []Episode{}

	result := h.Db.Find(&episodes)

	return episodes, result.Error
}

func NewEpisodes() EpisodeHandler {
	dbConnection := new(DB).getConnection()
	transmissionClient := NewTrs()
	return &Episodes{Db: dbConnection, transmission: transmissionClient}
}

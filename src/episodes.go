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

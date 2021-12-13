package transmissionrss

import (
	"gorm.io/gorm"
)

type Episode struct {
	gorm.Model
	ShowId    string
	EpisodeId string
	ShowTitle string
	Title     string
	Link      string
}

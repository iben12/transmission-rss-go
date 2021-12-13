package transmissionrss

import (
	"gorm.io/gorm"
)

func download(feedItems []FeedItem, db *gorm.DB) (downloaded []Episode, errs []error) {
	var episodesAdded []Episode

	for _, feedItem := range feedItems {
		episode := Episode{}
		result := db.Where(&Episode{ShowId: feedItem.ShowId, EpisodeId: feedItem.EpisodeId}).First(&episode)

		if result.Error != nil && result.Error.Error() == "record not found" {
			episode := Episode{
				Model:     gorm.Model{},
				ShowId:    feedItem.ShowId,
				EpisodeId: feedItem.EpisodeId,
				ShowTitle: feedItem.ShowTitle,
				Title:     feedItem.Title,
				Link:      feedItem.Link,
			}

			trs := Trs{}

			added := trs.addDownload(episode)

			if added {
				saved := db.Create(&episode)

				if saved.Error == nil {
					episodesAdded = append(episodesAdded, episode)
				} else {
					errs = append(errs, saved.Error)
				}
			}
		}
	}

	return episodesAdded, errs
}

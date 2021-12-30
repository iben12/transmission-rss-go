package transmissionrss

import (
	"errors"
	"fmt"
	"sync"

	"gorm.io/gorm"
)

var (
	TransmissionClient *Trs
	SlackClient        *SlackNotification
)

func init() {
	TransmissionClient = &Trs{}
	SlackClient = new(SlackNotification)
}

func download(feedItems []FeedItem, db *gorm.DB) (downloaded []Episode, errs []string) {
	episodesChannel := make(chan Episode, len(feedItems))
	errorChannel := make(chan error, len(feedItems))

	var wg sync.WaitGroup

	for _, feedItem := range feedItems {
		episode := Episode{}
		result := db.Where(&Episode{ShowId: feedItem.ShowId, EpisodeId: feedItem.EpisodeId}).First(&episode)

		if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
			wg.Add(1)
			go processEpisode(feedItem, episodesChannel, errorChannel, &wg, db)
		}
	}

	wg.Wait()

	close(episodesChannel)
	close(errorChannel)

	var episodesAdded []Episode

	for episode := range episodesChannel {
		episodesAdded = append(episodesAdded, episode)
	}

	for err := range errorChannel {
		errs = append(errs, err.Error())
	}

	if len(episodesAdded) > 0 {
		notify(episodesAdded)
	}

	fmt.Println(errs)

	return episodesAdded, errs
}

func notify(episodes []Episode) {
	title := "TransmissionRSS: New episode(s)"
	body := "Added episodes:"

	for _, episode := range episodes {
		body += "\n" + episode.Title
	}

	SlackClient.Send(title, body)
}

func processEpisode(
	feedItem FeedItem,
	episodesChannel chan Episode,
	errorChannel chan error,
	wg *sync.WaitGroup,
	db *gorm.DB) {

	defer wg.Done()

	episode := Episode{
		Model:     gorm.Model{},
		ShowId:    feedItem.ShowId,
		EpisodeId: feedItem.EpisodeId,
		ShowTitle: feedItem.ShowTitle,
		Title:     feedItem.Title,
		Link:      feedItem.Link,
	}

	transmissionError := TransmissionClient.AddDownload(episode)

	if transmissionError == nil {
		saved := db.Create(&episode)

		if saved.Error == nil {
			episodesChannel <- episode
		} else {
			errorChannel <- saved.Error
		}
	} else {
		errorChannel <- transmissionError
	}
}

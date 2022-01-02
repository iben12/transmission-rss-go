package transmissionrss

import (
	"errors"
	"fmt"
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

func download(feedItems []FeedItem, db *gorm.DB) ([]Episode, []string) {
	var (
		episodesAdded []Episode
		errs          []string
		resultChans   []chan Episode
		errorChans    []chan error
	)

	for _, feedItem := range feedItems {
		episode := Episode{}
		result := db.Where(&Episode{ShowId: feedItem.ShowId, EpisodeId: feedItem.EpisodeId}).First(&episode)

		if result.Error != nil && errors.Is(result.Error, gorm.ErrRecordNotFound) {
			resultChan := make(chan Episode, 1)
			errChan := make(chan error, 1)
			go processEpisode(feedItem, resultChan, errChan, db)

			errorChans = append(errorChans, errChan)
			resultChans = append(resultChans, resultChan)
		}
	}

	for _, resultChan := range resultChans {
		result := <-resultChan
		if result.ID != 0 {
			episodesAdded = append(episodesAdded, result)
		}
	}

	for _, errorChan := range errorChans {
		err := <-errorChan
		if err != nil {
			errs = append(errs, err.Error())
		}
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
	result chan Episode,
	err chan error,
	db *gorm.DB,
) {
	defer close(result)
	defer close(err)

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
			result <- episode
		} else {
			err <- saved.Error
		}
	} else {
		err <- transmissionError
	}
}

package transmissionrss

import (
	"errors"

	"gorm.io/gorm"
)

func Download(feedItems []FeedItem, episodes EpisodeHandler) ([]Episode, []string) {
	var (
		episodesAdded []Episode
		errs          []string
		resultChans   []chan Episode
		errorChans    []chan error
	)

	for _, feedItem := range feedItems {
		episodeToFind := &Episode{ShowId: feedItem.ShowId, EpisodeId: feedItem.EpisodeId}
		_, err := episodes.FindEpisode(episodeToFind)

		if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			resultChan := make(chan Episode, 1)
			errChan := make(chan error, 1)
			go processEpisode(feedItem, resultChan, errChan, episodes)

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

	if len(errs) > 0 {
		Logger.Warn().
			Str("action", "download error").
			Strs("errors", errs).
			Msg("Download errors")
	}

	return episodesAdded, errs
}

func processEpisode(
	feedItem FeedItem,
	result chan Episode,
	err chan error,
	episodeHandler EpisodeHandler,
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

	transmissionError := episodeHandler.DownloadEpisode(episode)

	if transmissionError == nil {
		dbError := episodeHandler.AddEpisode(&episode)

		if dbError == nil {
			result <- episode
		} else {
			err <- dbError
		}
	} else {
		err <- transmissionError
	}
}

package transmissionrss

import (
	"errors"
	"github.com/antchfx/xmlquery"
	"io/ioutil"
	"net/http"
	"strings"
)

type FeedItem struct {
	Title     string
	ShowTitle string
	EpisodeId string
	ShowId    string
	Link      string
}

type FeedHandler interface {
	FetchItems(r string) (items []FeedItem, err error)
}

type Feeds struct{}

func (f *Feeds) FetchItems(rssAddress string) (items []FeedItem, err error) {
	xml, err := f.fetchRss(rssAddress)

	if err != nil {
		Logger.Error().
			Str("action", "fetch feed").
			Str("url", rssAddress).
			Err(err)

		return nil, err
	}

	var feedItems []FeedItem

	feed, parseError := xmlquery.Parse(strings.NewReader(xml))

	if parseError != nil {
		Logger.Error().
			Str("action", "parse feed").
			Err(parseError)

		return nil, errors.New("cannot parse XML")
	}

	for _, item := range xmlquery.Find(feed, "//item") {
		title := item.SelectElement("//title").InnerText()
		showTitle := item.SelectElement("//tv:show_name").InnerText()
		showId := item.SelectElement("//tv:show_id").InnerText()
		episodeId := item.SelectElement("//tv:episode_id").InnerText()
		link := item.SelectElement("//link").InnerText()

		feedItems = append(feedItems, FeedItem{title, showTitle, episodeId, showId, link})
	}

	return feedItems, nil
}

func (f *Feeds) fetchRss(rssAddress string) (string, error) {
	resp, err := http.Get(rssAddress)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), nil
}

func NewFeeds() FeedHandler {
	return new(Feeds)
}

package transmissionrss

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/antchfx/xmlquery"
)

type Feed struct{}

type FeedItem struct {
	Title     string
	ShowTitle string
	EpisodeId string
	ShowId    string
	Link      string
}

func (f *Feed) parse(xml string) (items []FeedItem, err error) {

	var feedItems []FeedItem

	feed, parseError := xmlquery.Parse(strings.NewReader(xml))

	if parseError != nil {
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

func (f *Feed) fetchRss(rssAddress string) (string, error) {
	resp, err := http.Get(rssAddress)
	if err != nil {
		return "", fmt.Errorf("can't fetch RSS feed: %s", rssAddress)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), nil
}

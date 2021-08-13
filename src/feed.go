package main

import (
	"fmt"
	"github.com/antchfx/xmlquery"
	"io/ioutil"
	"net/http"
	"strings"
)

type Feed struct{}

type FeedItem struct {
	Title     string
	ShowTitle string
	EpisodeId string
	ShowId    string
	Link      string
}

func (f *Feed) fetch(rssAddress string) {
	resp, err := http.Get(rssAddress)
	if err != nil {
		panic(fmt.Sprintf("Can't fetch RSS feed: %s", rssAddress))
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	var feedItems []FeedItem

	feed, parseError := xmlquery.Parse(strings.NewReader(string(body)))
	if parseError != nil {
		panic("Couldn't parse XML feed")
	}
	for _, item := range xmlquery.Find(feed, "//item") {
		title := item.SelectElement("//title").InnerText()
		showTitle := item.SelectElement("//tv:show_name").InnerText()
		showId := item.SelectElement("//tv:show_id").InnerText()
		episodeId := item.SelectElement("//tv:episode_id").InnerText()
		link := item.SelectElement("//link").InnerText()

		feedItems = append(feedItems, FeedItem{title, showTitle, episodeId, showId, link})
	}

	fmt.Println("Title:", feedItems[0].Title)
	fmt.Println("Episode count:", len(feedItems))
}

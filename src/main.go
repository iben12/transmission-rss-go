package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	// "github.com/hekmon/transmissionrpc"
	// "os"
)

func main() {
	fmt.Println("Starting")

	feed := Feed{}
	feed.fetch("https://showrss.info/show/951.rss")

	trs := Trs{}

	finishedTorrents := trs.getFinished()
	fmt.Println("Found", len(finishedTorrents), "finished torrents.")

	DB := DB{}
	db := DB.connect()

	episode := Episode{ShowId: "12345", ShowTitle: "Wire", EpisodeId: "56789", Title: "S01E01-Hello", Link: "magnet:dfghjkhjwidfuicuds"}
	result := db.Create(&episode)
	if result.Error != nil {
		fmt.Printf("Cloud not create episode in DB: %s\n", episode.Title)
	}
	fmt.Printf("Created episode, ID: %d\n", episode.ID)

	episodes := []Episode{}
	db.Find(&episodes)

	fmt.Println("Episode count: ", len(episodes))
	fmt.Println("Title of first episode is:", episodes[0].Title)
}

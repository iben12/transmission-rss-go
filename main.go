package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/iben12/transmission-rss-go/src"
	"log"
	"net/http"
)

func handleRequests() {
	static := http.FileServer(http.Dir("./static"))

	router := mux.NewRouter().StrictSlash(true)

	api := new(transmissionrss.Api)

	router.HandleFunc("/api/feeds", api.Feeds)
	router.HandleFunc("/api/episodes", api.Episodes)
	router.HandleFunc("/api/download", api.Download)
	router.HandleFunc("/api/cleanup", api.Clean)
	router.PathPrefix("/").Handler(static)

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {

	fmt.Println("Server starting")
	handleRequests()
}

// func main() {
// 	fmt.Println("Starting")

// 	feed := Feed{}
// 	feed.fetch("https://showrss.info/show/951.rss")

// 	trs := Trs{}

// 	finishedTorrents := trs.getFinished()
// 	fmt.Println("Found", len(finishedTorrents), "finished torrents.")

// 	DB := DB{}
// 	db := DB.connect()

// 	episode := Episode{ShowId: "12345", ShowTitle: "Wire", EpisodeId: "56789", Title: "S01E01-Hello", Link: "magnet:dfghjkhjwidfuicuds"}
// 	result := db.Create(&episode)
// 	if result.Error != nil {
// 		fmt.Printf("Cloud not create episode in DB: %s\n", episode.Title)
// 	}
// 	fmt.Printf("Created episode, ID: %d\n", episode.ID)

// 	episodes := []Episode{}
// 	db.Find(&episodes)

// 	fmt.Println("Episode count: ", len(episodes))
// 	fmt.Println("Title of first episode is:", episodes[0].Title)

// 	var slack SlackNotification

// 	slack.Send("hello", "bello")
// }

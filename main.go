package main

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/iben12/transmission-rss-go/trss"
)

func handleRequests() {
	static := http.FileServer(http.Dir("./static"))

	router := mux.NewRouter().StrictSlash(true)

	api := transmissionrss.NewApi()

	router.HandleFunc("/api/feeds", api.Feeds)
	router.HandleFunc("/api/episodes", api.Episodes)
	router.HandleFunc("/api/download", api.Download)
	router.HandleFunc("/api/cleanup", api.Clean)
	router.PathPrefix("/").Handler(static)

	err := http.ListenAndServe("127.0.0.1:8080", router)
	transmissionrss.Logger.Fatal().Err(err)
}

func main() {
	transmissionrss.Logger.Info().
		Str("action", "start server").
		Msg("Server starting")

	handleRequests()
}

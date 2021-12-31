package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/iben12/transmission-rss-go/trss"
	"log"
	"net/http"
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

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {

	fmt.Println("Server starting")
	handleRequests()
}

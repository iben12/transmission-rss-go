package main

import (
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/iben12/transmission-rss-go/trss"

	_ "github.com/joho/godotenv/autoload"
)

func handleRequests() {
	static := http.FileServer(http.Dir("./static"))

	router := mux.NewRouter().StrictSlash(true)

	api := transmissionrss.NewApi()

	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, "{\"status\": \"OK\"}")
	})
	router.HandleFunc("/api/feeds", api.Feeds)
	router.HandleFunc("/api/episodes", api.Episodes)
	router.HandleFunc("/api/download", api.Download)
	router.HandleFunc("/api/cleanup", api.Clean)
	router.PathPrefix("/").Handler(static)

	if os.Getenv("REQUEST_LOGGING") == "true" {
		router.Use(loggingMiddleware)
	}

	err := http.ListenAndServe("0.0.0.0:8080", router)
	transmissionrss.Logger.Fatal().Err(err)
}

func loggingMiddleware(next http.Handler) http.Handler {
	excludedPaths := strings.Split(os.Getenv("LOGGING_EXCLUDE_PATHS"), ",")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !contains(excludedPaths, r.RequestURI) {
			transmissionrss.Logger.Info().
				Str("action", "request").
				Str("URI", r.RequestURI).
				Str("Method", r.Method).
				Msg("Incoming request")
		}

		next.ServeHTTP(w, r)
	})
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func main() {
	transmissionrss.Logger.Info().
		Str("action", "start server").
		Msg("Server starting")

	handleRequests()
}

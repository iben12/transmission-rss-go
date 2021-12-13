package transmissionrss

import (
	"encoding/json"
	_ "github.com/joho/godotenv/autoload"
	"io"
	"net/http"
	"os"
)

type Api struct{}

func (a *Api) Episodes(w http.ResponseWriter, r *http.Request) {
	db := new(DB).getConnection()

	episodes := []Episode{}
	db.Find(&episodes)

	json.NewEncoder(w).Encode(episodes)
}

func (a *Api) Download(w http.ResponseWriter, r *http.Request) {
	rssAddress := os.Getenv("RSS_FEED_ADDRESS")

	feed := new(Feed)
	xml, err1 := feed.fetchRss(rssAddress)
	if err1 != nil {
		w.WriteHeader(500)
		io.WriteString(w, "{\"error\": \"Could not fetch feed\"}")
	}

	feedItems, err2 := feed.parse(xml)

	if err2 != nil {
		w.WriteHeader(500)
		io.WriteString(w, "{\"error\": \"Could not parse feed\"}")
	}

	db := new(DB).getConnection()

	episodes, errs := download(feedItems, db)

	if errs != nil {
		w.WriteHeader(500)
		io.WriteString(w, "{\"error\": \"Could not download items\"}")
	}

	json.NewEncoder(w).Encode(episodes)
}

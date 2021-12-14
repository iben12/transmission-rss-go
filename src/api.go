package transmissionrss

import (
	"encoding/json"
	_ "github.com/joho/godotenv/autoload"
	"io"
	"net/http"
	"os"
)

type Api struct{}

func (a *Api) Feeds(w http.ResponseWriter, r *http.Request) {
	type Feed struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}

	feeds := []Feed{{Name: "ShowRSS", Url: os.Getenv("RSS_FEED_ADDRESS")}}

	json.NewEncoder(w).Encode(feeds)
}

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
		return
	}

	feedItems, err2 := feed.parse(xml)

	if err2 != nil {
		w.WriteHeader(500)
		io.WriteString(w, "{\"error\": \"Could not parse feed\"}")
		return
	}

	db := new(DB).getConnection()

	downloaded, errs := download(feedItems, db)

	if errs != nil {
		w.WriteHeader(500)
		io.WriteString(w, "{\"error\": \"Could not download items\"}")
		return
	}

	json.NewEncoder(w).Encode(downloaded)
}

func (a *Api) Clean(w http.ResponseWriter, r *http.Request) {
	trs := new(Trs)

	ids, titles := trs.getFinished()

	if len(ids) > 0 {
		err := trs.remove(ids)

		if err != nil {
			w.WriteHeader(500)
			io.WriteString(w, "{\"error\": \"Could not remove torrents\"}")
			return
		}

		title := "TransmissionRSS: Removed episode(s)"
		body := "Removed episodes:"

		for _, episode := range titles {
			body += "\n" + episode
		}

		slackNotification := new(SlackNotification)

		slackNotification.Send(title, body)
	}

	json.NewEncoder(w).Encode(titles)
}

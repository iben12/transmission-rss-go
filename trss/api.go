package transmissionrss

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	episodes EpisodeHandler
)

func NewApi() *Api {
	episodes = NewEpisodeHanlder()

	return new(Api)
}

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
	episodeList, err := episodes.All()

	if err != nil {
		w.WriteHeader(500)
		io.WriteString(w, "{\"error\": \"Could not read episodes\"}")
		return
	}

	json.NewEncoder(w).Encode(episodeList)
}

func (a *Api) Download(w http.ResponseWriter, r *http.Request) {
	rssAddress := os.Getenv("RSS_FEED_ADDRESS")

	feed := new(Feed)

	feedItems, fetchError := feed.FetchItems(rssAddress)

	if fetchError != nil {
		w.WriteHeader(500)
		io.WriteString(w, "{\"error\": \"Could not parse feed\"}")
		return
	}

	downloaded, errs := Download(feedItems, episodes)

	if len(downloaded) > 0 {
		notify(downloaded)
		Logger.Info().
			Str("action", "downloaded").
			Int("count", len(downloaded)).
			Msg("Downloaded episodes")
	}

	type Response struct {
		Errors   []string
		Episodes []Episode
	}

	response := Response{
		Errors:   errs,
		Episodes: downloaded,
	}

	json.NewEncoder(w).Encode(response)
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

func notify(episodes []Episode) {
	SlackClient := &SlackNotification{}
	title := "TransmissionRSS: New episode(s)"
	body := "Added episodes:"

	for _, episode := range episodes {
		body += "\n" + episode.Title
	}

	SlackClient.Send(title, body)
}

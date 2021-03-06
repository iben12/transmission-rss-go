package transmissionrss

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Notification interface {
	Send(title string, body string) error
}

type SlackNotification struct{}

type SlackTextBlock struct {
	Text string `json:"text"`
	Type string `json:"type"`
}

type SlackBlock struct {
	Type string `json:"type"`
	Text struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"text"`
}

type SlackPayload struct {
	Channel string       `json:"channel"`
	Text    string       `json:"text"`
	Blocks  []SlackBlock `json:"blocks"`
}

func (s *SlackNotification) Send(title string, body string) error {
	payload := s.renderPayload(title, body)
	json, _ := json.Marshal(payload)

	resp, err := http.Post(os.Getenv("SLACK_URL"), "application/json", bytes.NewBuffer(json))

	if err != nil {
		return err
	}

	defer resp.Body.Close()
	respBytes, _ := ioutil.ReadAll(resp.Body)

	if string(respBytes) != "ok" {
		return errors.New(string(respBytes))
	}

	return nil
}

func (s *SlackNotification) renderPayload(title string, body string) (p SlackPayload) {
	payload := SlackPayload{
		Channel: os.Getenv("SLACK_CHANNEL"),
		Text:    title,
		Blocks: []SlackBlock{
			{
				Type: "section",
				Text: SlackTextBlock{
					Text: fmt.Sprintf(":arrow_up_down: *%s*\n%s", title, body),
					Type: "mrkdwn",
				},
			},
		},
	}

	return payload
}

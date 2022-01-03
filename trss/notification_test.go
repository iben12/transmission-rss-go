package transmissionrss_test

import (
	"bytes"
	"github.com/iben12/transmission-rss-go/tests/mocks"
	"github.com/iben12/transmission-rss-go/trss"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func TestNotification(t *testing.T) {
	os.Setenv("SLACK_URL", "https://hooks.slack.com/services/testurl")
	transmissionrss.Client = &mocks.MockHttpClient{}

	t.Run("Notification gets through", func(t *testing.T) {
		response := ioutil.NopCloser(bytes.NewReader([]byte("ok")))
		mocks.GetHttpDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 200,
				Body:       response,
			}, nil
		}

		notification := new(transmissionrss.SlackNotification)
		err := notification.Send("hello", "bello")

		if err != nil {
			t.Error("Expected no error, but got:", err)
		}
	})

	t.Run("Notification fails", func(t *testing.T) {
		errorMessage := "invalid payload"
		response := ioutil.NopCloser(bytes.NewReader([]byte(errorMessage)))
		mocks.GetHttpDoFunc = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       response,
			}, nil
		}

		notification := new(transmissionrss.SlackNotification)
		err := notification.Send("hello", "bello")

		if err == nil {
			t.Error("Expected error, but got nil")
		}

		if err.Error() != errorMessage {
			t.Errorf("Expected error '%s', but got '%s'", errorMessage, err.Error())
		}
	})

}

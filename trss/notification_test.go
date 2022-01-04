package transmissionrss_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/iben12/transmission-rss-go/tests/mocks"
	"github.com/iben12/transmission-rss-go/trss"
	"github.com/stretchr/testify/assert"
)

func TestNotification(t *testing.T) {
	assert := assert.New(t)
	os.Setenv("SLACK_URL", "https://hooks.slack.com/services/testurl")
	transmissionrss.HttpClient = &mocks.MockHttpClient{}

	t.Run("Notification gets through", func(t *testing.T) {
		response := ioutil.NopCloser(bytes.NewReader([]byte("ok")))
		var requestBody []byte
		mocks.MockHttpDo = func(req *http.Request) (*http.Response, error) {
			defer req.Body.Close()
			requestBody, _ = ioutil.ReadAll(req.Body)
			return &http.Response{
				StatusCode: 200,
				Body:       response,
			}, nil
		}

		notification := new(transmissionrss.SlackNotification)
		err := notification.Send("hello", "bello")

		assert.Nil(err)

		assert.Equal(string(requestBody), `{"channel":"","text":"hello","blocks":[{"type":"section","text":{"text":":arrow_up_down: *hello*\nbello","type":"mrkdwn"}}]}`)
	})

	t.Run("Notification fails", func(t *testing.T) {
		errorMessage := "invalid payload"
		response := ioutil.NopCloser(bytes.NewReader([]byte(errorMessage)))
		mocks.MockHttpDo = func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: 400,
				Body:       response,
			}, nil
		}

		notification := new(transmissionrss.SlackNotification)
		err := notification.Send("hello", "bello")

		if assert.Error(err) {
			assert.Equal(err.Error(), errorMessage)
		}
	})

}

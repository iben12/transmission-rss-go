package transmissionrss_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	helpers "github.com/iben12/transmission-rss-go/tests/helpers"
	trss "github.com/iben12/transmission-rss-go/trss"
	"github.com/stretchr/testify/assert"
)

func TestNotification(t *testing.T) {
	assert := assert.New(t)

	t.Run("Notification gets through", func(t *testing.T) {
		expectedRequestBody := `{"channel":"","text":"hello","blocks":[{"type":"section","text":{"text":":arrow_up_down: *hello*\nbello","type":"mrkdwn"}}]}`
		// server := httptest.NewServer(http.HandlerFunc())

		handler := func(rw http.ResponseWriter, req *http.Request) {
			assert.Equal("/services/testurl", req.URL.String())

			defer req.Body.Close()
			body, _ := ioutil.ReadAll(req.Body)
			assert.Equal(expectedRequestBody, string(body))

			rw.Write([]byte("ok"))
		}

		server := helpers.CreateServer(handler)

		os.Setenv("SLACK_URL", server.URL+"/services/testurl")

		notification := new(trss.SlackNotification)
		err := notification.Send("hello", "bello")

		assert.Nil(err)
	})

	t.Run("Notification fails", func(t *testing.T) {
		errorMessage := "invalid payload"
		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(400)
			rw.Write([]byte(errorMessage))
		}))

		os.Setenv("SLACK_URL", server.URL+"/services/testurl")

		notification := new(trss.SlackNotification)
		err := notification.Send("hello", "bello")

		if assert.Error(err) {
			assert.Equal(err.Error(), errorMessage)
		}
	})

}

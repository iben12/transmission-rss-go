package transmissionrss_test

import (
	"io/ioutil"
	"net/http"
	"os"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	helpers "github.com/iben12/transmission-rss-go/tests/helpers"
	trss "github.com/iben12/transmission-rss-go/trss"
)

var _ = Describe("Notification", func() {
	AfterEach(func() {
		os.Clearenv()
	})

	It("sends notification", func() {
		expectedRequestBody := `{"channel":"","text":"hello","blocks":[{"type":"section","text":{"text":":arrow_up_down: *hello*\nbello","type":"mrkdwn"}}]}`

		handler := func(rw http.ResponseWriter, req *http.Request) {
			Expect(req.URL.String()).To(Equal("/services/testurl"))
			defer req.Body.Close()
			body, _ := ioutil.ReadAll(req.Body)

			Expect(string(body)).To(Equal(expectedRequestBody))

			rw.Write([]byte("ok"))
		}

		server := helpers.CreateServer(handler)

		os.Setenv("SLACK_URL", server.URL+"/services/testurl")

		notification := new(trss.SlackNotification)
		err := notification.Send("hello", "bello")

		Expect(err).NotTo(HaveOccurred())
	})

	It("errors if notification fails", func() {
		errorMessage := "invalid payload"

		handler := func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(400)
			rw.Write([]byte(errorMessage))
		}

		server := helpers.CreateServer(handler)

		os.Setenv("SLACK_URL", server.URL+"/services/testurl")

		notification := new(trss.SlackNotification)
		err := notification.Send("hello", "bello")

		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(errorMessage))
	})
})

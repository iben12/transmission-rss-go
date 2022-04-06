package transmissionrss_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"time"

	"github.com/hekmon/transmissionrpc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	helpers "github.com/iben12/transmission-rss-go/tests/helpers"
	trss "github.com/iben12/transmission-rss-go/trss"
)

var _ = Describe("Transmission", func() {
	var client *trss.Trs
	var server *httptest.Server
	Context("check version", func() {
		AfterEach(func() {
			server.Close()
		})

		It("is OK", func() {
			arguments := fmt.Sprintf(`{"rpc-version-minimum": %v, "rpc-version": 17}`, transmissionrpc.RPCVersion)
			response := &TransmissionResponse{
				arguments: arguments,
				status:    200,
				result:    "success",
			}
			client, server = setUpTransmissionTestServer([]*TransmissionResponse{response})

			err := client.CheckVersion()

			Expect(err).NotTo(HaveOccurred())
		})

		It("is not OK", func() {
			arguments := fmt.Sprintf(`{"rpc-version-minimum": %v, "rpc-version": 17}`, transmissionrpc.RPCVersion+1)
			response := &TransmissionResponse{
				arguments: arguments,
				status:    200,
				result:    "success",
			}
			client, server = setUpTransmissionTestServer([]*TransmissionResponse{response})

			err := client.CheckVersion()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("remote transmission RPC version (v17) is incompatible with the transmission library (v15): remote needs at least v16"))
		})

		It("server fails", func() {
			response := &TransmissionResponse{
				arguments: "{}",
				status:    500,
				result:    "failure",
			}
			client, server = setUpTransmissionTestServer([]*TransmissionResponse{response})

			err := client.CheckVersion()

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("can't get session values: 'session-get' rpc method failed: HTTP error 500: Internal Server Error"))
		})
	})

	Context("add torrent", func() {
		It("succeeds", func() {
			episode := &trss.Episode{
				ShowId:    "1",
				EpisodeId: "2",
				ShowTitle: "Show",
				Title:     "Episode",
				Link:      "url",
			}
			arguments := `{"torrent-added":{"name":"Torrent Name","id": 2}}`
			response := &TransmissionResponse{
				arguments: arguments,
				status:    200,
				result:    "success",
			}
			client, server = setUpTransmissionTestServer([]*TransmissionResponse{response})

			err := client.AddTorrent(*episode)

			Expect(err).NotTo(HaveOccurred())
		})

		It("is duplicate", func() {
			episode := &trss.Episode{
				ShowId:    "1",
				EpisodeId: "2",
				ShowTitle: "Show",
				Title:     "Episode",
				Link:      "url",
			}
			arguments := `{"torrent-duplicate":{"name":"Torrent Name","id": 2}}`
			response := &TransmissionResponse{
				arguments: arguments,
				status:    200,
				result:    "success",
			}
			client, server = setUpTransmissionTestServer([]*TransmissionResponse{response})

			err := client.AddTorrent(*episode)

			Expect(err).NotTo(HaveOccurred())
		})

		It("request fails", func() {
			episode := &trss.Episode{
				ShowId:    "1",
				EpisodeId: "2",
				ShowTitle: "Show",
				Title:     "Episode",
				Link:      "url",
			}
			response := &TransmissionResponse{
				arguments: "{}",
				status:    500,
				result:    "failure",
			}
			client, server = setUpTransmissionTestServer([]*TransmissionResponse{response})

			err := client.AddTorrent(*episode)

			Expect(err).To(HaveOccurred())
			Expect(err).To(MatchError("'torrent-add' rpc method failed: HTTP error 500: Internal Server Error"))
		})
	})

	Context("clean finished", func() {
		It("succeeds", func() {
			arguments := `{"torrents":[{"name": "Torrent name", "id": 5, "isFinished": true}]}`
			response1 := &TransmissionResponse{
				arguments: arguments,
				status:    200,
				result:    "success",
			}
			response2 := &TransmissionResponse{
				arguments: arguments,
				status:    200,
				result:    "success",
			}
			client, server = setUpTransmissionTestServer([]*TransmissionResponse{response1, response2})

			titles, err := client.CleanFinished()

			Expect(err).NotTo(HaveOccurred())
			Expect(titles).To(ContainElement("Torrent name"))
		})

		It("nothing to clean", func() {
			arguments := `{"torrents":[{"name": "Torrent name", "id": 5, "isFinished": false}]}`
			response := &TransmissionResponse{
				arguments: arguments,
				status:    200,
				result:    "success",
			}
			client, server = setUpTransmissionTestServer([]*TransmissionResponse{response})

			titles, err := client.CleanFinished()

			Expect(err).NotTo(HaveOccurred())
			Expect(titles).To(BeEmpty())
		})

		It("server fails for getAllTorrents", func() {
			response := &TransmissionResponse{
				arguments: "{}",
				status:    500,
				result:    "failure",
			}
			client, server = setUpTransmissionTestServer([]*TransmissionResponse{response})

			titles, err := client.CleanFinished()

			Expect(err).To(HaveOccurred())
			Expect(titles).To(BeEmpty())
		})

		It("server fails for remove torrents", func() {
			arguments := `{"torrents":[{"name": "Torrent name", "id": 5, "isFinished": true}]}`
			response1 := &TransmissionResponse{
				arguments: arguments,
				status:    200,
				result:    "success",
			}
			response2 := &TransmissionResponse{
				arguments: "{}",
				status:    500,
				result:    "failure",
			}
			client, server = setUpTransmissionTestServer([]*TransmissionResponse{response1, response2})

			titles, err := client.CleanFinished()

			Expect(err).To(HaveOccurred())
			Expect(titles).To(BeEmpty())
		})
	})
})

type TransmissionResponse struct {
	arguments string
	status    int
	result    string
}

func setUpTransmissionTestServer(responses []*TransmissionResponse) (*trss.Trs, *httptest.Server) {
	counter := 0
	server := helpers.CreateServer(func(w http.ResponseWriter, r *http.Request) {
		Expect(r.URL.String()).To(Equal("/transmission/rpc"))
		body, _ := ioutil.ReadAll(r.Body)
		GinkgoWriter.Println(string(body))
		GinkgoWriter.Printf("counter is: %v", counter)
		response := responses[counter]
		counter++
		var request map[string]int
		json.Unmarshal(body, &request)
		tag := fmt.Sprint(request["tag"])
		responseString := fmt.Sprintf(`{"arguments": %v, "result": "%v", "tag": %v}`, response.arguments, response.result, tag)
		w.WriteHeader(response.status)
		w.Write([]byte(responseString))
	})
	testUrl, _ := url.Parse(server.URL)
	port, _ := strconv.ParseUint(testUrl.Port(), 10, 0)
	const timeout time.Duration = 12 * time.Second
	transmissionbt, err := transmissionrpc.New(
		testUrl.Hostname(),
		"username",
		"secret",
		&transmissionrpc.AdvancedConfig{
			HTTPTimeout: timeout,
			Port:        uint16(port),
		})

	Expect(err).NotTo(HaveOccurred())

	return &trss.Trs{Client: transmissionbt, AddPaused: true}, server
}

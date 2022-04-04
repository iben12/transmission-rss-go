package transmissionrss_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/hekmon/transmissionrpc"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	helpers "github.com/iben12/transmission-rss-go/tests/helpers"
	trss "github.com/iben12/transmission-rss-go/trss"
)

var (
	client *trss.Trs
)

var _ = Describe("Transmission", func() {
	Context("check version", func() {
		It("is OK", func() {
			arguments := fmt.Sprintf(`{"rpc-version-minimum": %v, "rpc-version": 17}`, transmissionrpc.RPCVersion)
			setUpTransmissionTestServer(arguments)

			result := client.CheckVersion()

			Expect(result).To(BeTrue())
		})

		It("is not OK", func() {
			arguments := fmt.Sprintf(`{"rpc-version-minimum": %v, "rpc-version": 17}`, transmissionrpc.RPCVersion+1)
			setUpTransmissionTestServer(arguments)

			result := client.CheckVersion()

			Expect(result).To(BeFalse())
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
			setUpTransmissionTestServer(arguments)

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
			setUpTransmissionTestServer(arguments)

			err := client.AddTorrent(*episode)

			Expect(err).NotTo(HaveOccurred())
		})
	})

	Context("clean finished", func() {
		It("succeeds", func() {
			arguments := `{"torrents":[{"name": "Torrent name", "id": 5, "isFinished": true}]}`
			setUpTransmissionTestServer(arguments)

			titles, err := client.CleanFinished()

			Expect(err).NotTo(HaveOccurred())
			Expect(titles).To(ContainElement("Torrent name"))
		})

		It("nothing to clean", func() {
			arguments := `{"torrents":[{"name": "Torrent name", "id": 5, "isFinished": false}]}`
			setUpTransmissionTestServer(arguments)

			titles, err := client.CleanFinished()

			Expect(err).NotTo(HaveOccurred())
			Expect(titles).To(BeEmpty())
		})
	})
})

func setUpTransmissionTestServer(responseArguments string) {
	server := helpers.CreateServer(func(w http.ResponseWriter, r *http.Request) {
		Expect(r.URL.String()).To(Equal("/transmission/rpc"))
		body, _ := ioutil.ReadAll(r.Body)
		GinkgoWriter.Println(string(body))
		var request map[string]int
		json.Unmarshal(body, &request)
		tag := fmt.Sprint(request["tag"])
		responseString := fmt.Sprintf(`{"arguments": %v, "result": "success", "tag": %v}`, responseArguments, tag)
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

	client = &trss.Trs{Client: *transmissionbt, AddPaused: true}
}

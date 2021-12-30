package transmissionrss

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hekmon/transmissionrpc"
)

type Trs struct{}

var client transmissionrpc.Client

var addPaused, _ = strconv.ParseBool(os.Getenv("ADD_PAUSED"))

func init() {
	const timeout time.Duration = 12 * time.Second
	transmissionbt, err := transmissionrpc.New(
		os.Getenv("TRANSMISSION_HOST"),
		os.Getenv("TRANSMISSION_USER"),
		os.Getenv("TRANSMISSION_PASSWORD"),
		&transmissionrpc.AdvancedConfig{
			HTTPTimeout: timeout,
		})
	if err != nil {
		panic(err)
	} else {
		client = *transmissionbt
	}
}

func (trs *Trs) CheckVersion() (ok bool) {
	ok, serverVersion, serverMinimumVersion, err := client.RPCVersion()
	if err != nil {
		panic(err)
	}
	if !ok {
		panic(fmt.Sprintf("Remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): remote needs at least v%d",
			serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion))
	}
	fmt.Printf("Remote transmission RPC version (v%d) is compatible with our transmissionrpc library (v%d)\n",
		serverVersion, transmissionrpc.RPCVersion)

	return true
}

func (trs *Trs) getAllTorrents(fields []string) (t []*transmissionrpc.Torrent) {
	torrents, err := client.TorrentGet(fields, nil)
	if err != nil {
		panic(err)
	} else {
		return torrents
	}
}

func (trs *Trs) getFinished() (ids []int64, titles []string) {
	torrents := trs.getAllTorrents([]string{"id", "isFinished", "name"})
	finishedTorrentIds := []int64{}
	finishedTorrentTitles := []string{}
	for _, torrent := range torrents {
		if *torrent.IsFinished {
			finishedTorrentIds = append(finishedTorrentIds, *torrent.ID)
			finishedTorrentTitles = append(finishedTorrentTitles, *torrent.Name)
		}
	}

	return finishedTorrentIds, finishedTorrentTitles
}

func (trs *Trs) AddDownload(episode Episode) error {
	torrentToAdd := &transmissionrpc.TorrentAddPayload{Filename: &episode.Link, Paused: &addPaused}

	_, err := client.TorrentAdd(torrentToAdd)

	return err
}

func (trs *Trs) remove(ids []int64) error {
	payload := &transmissionrpc.TorrentRemovePayload{IDs: ids, DeleteLocalData: false}
	err := client.TorrentRemove(payload)

	return err
}

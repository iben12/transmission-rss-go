package transmissionrss

import (
	"fmt"
	"github.com/hekmon/transmissionrpc"
	"os"
	"time"
)

type Trs struct{}

var client transmissionrpc.Client

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

func (trs *Trs) getFinished() (t []int64) {
	torrents := trs.getAllTorrents([]string{"id", "isFinished"})
	finishedTorrents := []int64{}
	for _, torrent := range torrents {
		if *torrent.IsFinished {
			finishedTorrents = append(finishedTorrents, *torrent.ID)
		}
	}

	return finishedTorrents
}

func (trs *Trs) addDownload(episode Episode) bool {
	paused := true

	torrentToAdd := &transmissionrpc.TorrentAddPayload{Filename: &episode.Link, Paused: &paused}

	_, err := client.TorrentAdd(torrentToAdd)

	return err == nil
}

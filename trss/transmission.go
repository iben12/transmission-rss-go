package transmissionrss

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hekmon/transmissionrpc"
)

type TransmissionService interface {
	CheckVersion() bool
	AddTorrent(e Episode) error
	CleanFinished() ([]string, error)
}

type Trs struct {
	client    transmissionrpc.Client
	addPaused bool
}

func NewTrs() *Trs {
	paused, _ := strconv.ParseBool(os.Getenv("ADD_PAUSED"))
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
		return &Trs{
			client:    *transmissionbt,
			addPaused: paused,
		}
	}
}

func (trs *Trs) CheckVersion() bool {
	ok, serverVersion, serverMinimumVersion, err := trs.client.RPCVersion()
	if err != nil {
		Logger.Error().
			Str("action", "transmission check version").
			Err(err).Msg("")
	}
	if !ok {
		Logger.Fatal().
			Str("action", "transmission check version").
			Err(fmt.Errorf("remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): remote needs at least v%d",
				serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion))
	}

	Logger.Info().
		Str("action", "transmission check version").
		Msg(fmt.Sprintf("Remote transmission RPC version (v%d) is compatible with our transmissionrpc library (v%d)",
			serverVersion, transmissionrpc.RPCVersion))

	return true
}

func (trs *Trs) AddTorrent(episode Episode) error {
	torrentToAdd := &transmissionrpc.TorrentAddPayload{Filename: &episode.Link, Paused: &trs.addPaused}

	_, err := trs.client.TorrentAdd(torrentToAdd)

	return err
}

func (trs *Trs) CleanFinished() ([]string, error) {
	ids, titles := trs.getFinished()

	err := trs.remove(ids)

	if err != nil {
		Logger.Error().
			Str("action", "remove torrents").
			Err(err).Msg("")
	}

	return titles, err
}

func (trs *Trs) getAllTorrents(fields []string) (t []*transmissionrpc.Torrent) {
	torrents, err := trs.client.TorrentGet(fields, nil)
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

func (trs *Trs) remove(ids []int64) error {
	payload := &transmissionrpc.TorrentRemovePayload{IDs: ids, DeleteLocalData: false}
	err := trs.client.TorrentRemove(payload)

	return err
}

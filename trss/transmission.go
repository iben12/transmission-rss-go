package transmissionrss

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/hekmon/transmissionrpc"
)

type TransmissionService interface {
	CheckVersion() error
	AddTorrent(e Episode) error
	CleanFinished() ([]string, error)
}

type Trs struct {
	Client    *transmissionrpc.Client
	AddPaused bool
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
		Logger.Fatal().Str("action", "Transmission connect").Err(err)
		return nil
	}

	trs := &Trs{
		Client:    transmissionbt,
		AddPaused: paused,
	}

	trsErr := trs.CheckVersion()

	if trsErr != nil {
		Logger.Fatal().Err(trsErr)
		return nil
	} else {
		return &Trs{
			Client:    transmissionbt,
			AddPaused: paused,
		}
	}
}

func (trs *Trs) CheckVersion() error {
	ok, serverVersion, serverMinimumVersion, err := trs.Client.RPCVersion()
	if err != nil {
		Logger.Error().
			Str("action", "transmission check version").
			Err(err).Msg("")
		return err
	}
	if !ok {
		err := fmt.Errorf("remote transmission RPC version (v%d) is incompatible with the transmission library (v%d): remote needs at least v%d",
			serverVersion, transmissionrpc.RPCVersion, serverMinimumVersion)
		Logger.Error().
			Str("action", "transmission check version").
			Err(err).Msg("")
		return err
	}

	Logger.Info().
		Str("action", "transmission check version").
		Msg(fmt.Sprintf("Remote transmission RPC version (v%d) is compatible with our transmissionrpc library (v%d)",
			serverVersion, transmissionrpc.RPCVersion))

	return nil
}

func (trs *Trs) AddTorrent(episode Episode) error {
	dir := fmt.Sprintf("/videos/Series/%v", episode.ShowTitle)
	torrentToAdd := &transmissionrpc.TorrentAddPayload{Filename: &episode.Link, DownloadDir: &dir, Paused: &trs.AddPaused}

	_, err := trs.Client.TorrentAdd(torrentToAdd)

	return err
}

func (trs *Trs) CleanFinished() ([]string, error) {
	ids, titles, err := trs.getFinished()

	if err != nil {
		Logger.Error().
			Str("action", "get finished torrents").
			Err(err).Msg("")
		return []string{}, err
	}

	if len(ids) == 0 {
		Logger.Info().
			Str("action", "get finished torrents").
			Msg("no finished torrents")
		return []string{}, nil
	}

	removeErr := trs.remove(ids)

	if removeErr != nil {
		Logger.Error().
			Str("action", "remove torrents").
			Err(removeErr).Msg("")
		return []string{}, removeErr
	}

	return titles, nil
}

func (trs *Trs) getAllTorrents(fields []string) ([]*transmissionrpc.Torrent, error) {
	torrents, err := trs.Client.TorrentGet(fields, nil)
	if err != nil {
		Logger.Error().Str("action", "get all torrents").
			Err(err).Msg("")
		return nil, err
	} else {
		return torrents, nil
	}
}

func (trs *Trs) getFinished() ([]int64, []string, error) {
	torrents, err := trs.getAllTorrents([]string{"id", "isFinished", "name"})
	if err != nil {
		return nil, nil, err
	}
	finishedTorrentIds := []int64{}
	finishedTorrentTitles := []string{}
	for _, torrent := range torrents {
		if *torrent.IsFinished {
			finishedTorrentIds = append(finishedTorrentIds, *torrent.ID)
			finishedTorrentTitles = append(finishedTorrentTitles, *torrent.Name)
		}
	}

	return finishedTorrentIds, finishedTorrentTitles, nil
}

func (trs *Trs) remove(ids []int64) error {
	payload := &transmissionrpc.TorrentRemovePayload{IDs: ids, DeleteLocalData: false}
	err := trs.Client.TorrentRemove(payload)

	return err
}

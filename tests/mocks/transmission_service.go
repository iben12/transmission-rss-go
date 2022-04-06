package mocks

import (
	trss "github.com/iben12/transmission-rss-go/trss"
)

type MockTransmissionService struct {
	MockCheckVersion  func() error
	MockAddTorrent    func(e trss.Episode) error
	MockCleanFinished func() ([]string, error)
}

func (t *MockTransmissionService) CheckVersion() error {
	return t.MockCheckVersion()
}

func (t *MockTransmissionService) AddTorrent(e trss.Episode) error {
	return t.MockAddTorrent(e)
}

func (t *MockTransmissionService) CleanFinished() ([]string, error) {
	return t.MockCleanFinished()
}

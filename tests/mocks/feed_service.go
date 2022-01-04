package mocks

import (
	trss "github.com/iben12/transmission-rss-go/trss"
)

type MockFeeds struct {
	MockFetchItems func(r string) ([]trss.FeedItem, error)
}

func (f *MockFeeds) FetchItems(r string) ([]trss.FeedItem, error) {
	return f.MockFetchItems(r)
}

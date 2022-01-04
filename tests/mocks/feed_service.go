package mocks

import (
	trss "github.com/iben12/transmission-rss-go/trss"
)

type MockFeeds struct {
	MockFetchItems func(r string) (items []trss.FeedItem, err error)
}

func (f *MockFeeds) FetchItems(r string) (items []trss.FeedItem, err error) {
	return f.MockFetchItems(r)
}

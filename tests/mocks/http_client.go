package mocks

import "net/http"

type MockHttpClient struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	// GetHttpDoFunc fetches the mock client's `Do` func
	GetHttpDoFunc func(req *http.Request) (*http.Response, error)
)

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return GetHttpDoFunc(req)
}

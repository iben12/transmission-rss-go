package mocks

import "net/http"

type MockHttpClient struct {
	// DoFunc func(req *http.Request) (*http.Response, error)
}

var (
	// MockHttpDo fetches the mock client's `Do` func
	MockHttpDo func(req *http.Request) (*http.Response, error)
)

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return MockHttpDo(req)
}

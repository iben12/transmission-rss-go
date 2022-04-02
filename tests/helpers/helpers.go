package test_helpers

import (
	"net/http"
	"net/http/httptest"
)

func CreateServer(handler func(rw http.ResponseWriter, req *http.Request)) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(handler))

	return server
}

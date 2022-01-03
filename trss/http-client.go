package transmissionrss

import (
	"net/http"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var (
	HttpClient HTTPClient
)

func init() {
	HttpClient = &http.Client{}
}

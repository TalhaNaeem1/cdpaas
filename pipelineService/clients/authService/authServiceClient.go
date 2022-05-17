//go:generate mockgen -destination=mocks/mock_httpclient.go -package=mock_httpclient . HttpClient
package authService

import (
	"io"
	"net/http"
)

type AuthServiceClient interface {
	AuthServiceQuerier
}

type AuthServiceHttpClient struct {
	*RequestMaker
}

func NewClient(httpClient HttpClient) AuthServiceClient {
	return &AuthServiceHttpClient{
		RequestMaker: NewRequestMaker(httpClient),
	}
}

type HttpClient interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
	Get(url string) (resp *http.Response, err error)
}
type RequestMaker struct {
	client HttpClient
}

func NewRequestMaker(rm HttpClient) *RequestMaker {
	return &RequestMaker{client: rm}
}

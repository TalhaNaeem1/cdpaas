package airbyte

import (
	"io"
	"net/http"
)

type AirByteClient interface {
	AirByteQuerier
}

type AirByteHttpClient struct {
	*RequestMaker
}

func NewClient(httpClient *http.Client) AirByteClient {
	return &AirByteHttpClient{
		RequestMaker: NewRequestMaker(httpClient),
	}
}

type HttpClient interface {
	Post(url string, contentType string, body io.Reader) (resp *http.Response, err error)
}
type RequestMaker struct {
	client HttpClient
}

func NewRequestMaker(rm HttpClient) *RequestMaker {
	return &RequestMaker{client: rm}
}

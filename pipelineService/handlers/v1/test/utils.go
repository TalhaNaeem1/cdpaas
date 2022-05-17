package test

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

//var GlobalContext gin.Context

//ReqResBodyMatcher checks if the request and response are deeply equal, request object should be marshalled in json.
func ReqResBodyMatcher(t *testing.T, expected *bytes.Buffer, actual []byte) {
	var req map[string]interface{}

	err := json.Unmarshal(actual, &req)
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(expected)
	require.NoError(t, err)

	var res map[string]interface{}
	err = json.Unmarshal(data, &res)

	require.NoError(t, err)
	require.Equal(t, req, res)
}

// MakeHttpRequest requests the http server and return the response in response recorder.
func MakeHttpRequest(r http.Handler, requestType string, path string, query map[string]string, body []byte) (*httptest.ResponseRecorder, error) {
	var requestBody io.Reader = nil
	if len(body) != 0 {
		requestBody = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(requestType, path, requestBody)
	if err != nil {
		return nil, err
	}

	recorder := httptest.NewRecorder()

	if query != nil {
		q := req.URL.Query()
		for k, v := range query {
			q.Add(k, v)
		}

		req.URL.RawQuery = q.Encode()
	}

	MockAddAuthorization(req)

	r.ServeHTTP(recorder, req)

	return recorder, nil
}

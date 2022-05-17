package airbyte

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/pkg/errors"
	"pipelineService/utils"
)

func (airByteClient *RequestMaker) sendRequest(airByteURL string, reqBody *bytes.Buffer) ([]byte, error) {
	logger := utils.GetLogger()

	var body []byte

	if reqBody == nil {
		reqBody = new(bytes.Buffer)
	}

	res, err := airByteClient.client.Post(airByteURL, "application/json", reqBody)
	if err != nil {
		logger.Error("request to airbyte failed")

		return body, err
	} else if res.StatusCode != http.StatusOK {
		body, _= ioutil.ReadAll(res.Body)
		logger.Error(string(body))
		err = errors.New("request to airbyte was not successful")
		logger.Error(err.Error())

		return body, err
	}
	defer res.Body.Close()

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("failed to read response body from airbyte")

		return body, err
	}

	return body, nil
}

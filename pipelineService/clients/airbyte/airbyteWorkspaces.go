package airbyte

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	"pipelineService/env"
	"pipelineService/models/v1"
	"pipelineService/utils"
)

func (airByteClient *RequestMaker) CreateWorkspace(workspace models.WorkspaceRequest) (models.WorkspaceAPIResponse, error){
		logger := utils.GetLogger()

		airByteURL := fmt.Sprintf("%s/api/v1/workspaces/create", env.Env.AirByteAddress)

		var response models.WorkspaceAPIResponse

		jsonData, err := json.Marshal(workspace)
		if err != nil {
			logger.Error("failed to convert request body to json")

			return response, err
		}

		body, err := airByteClient.sendRequest(airByteURL, bytes.NewBuffer(jsonData))
		if err != nil {
			logger.Error("failed to read response body from airbyte")

			return response, err
		}

		err = json.Unmarshal(body, &response)
		if err != nil {
			logger.Error("failed to read response body from airbyte")

			return response, err
		}

		return response, nil
}

func (airByteClient *RequestMaker) GetWorkspaceID() (string, error) {
	logger := utils.GetLogger()

	airByteURL := fmt.Sprintf("%s/api/v1/workspaces/list", env.Env.AirByteAddress)

	body, err := airByteClient.sendRequest(airByteURL, nil)
	if err != nil {
		logger.Error("failed to read response body from airbyte")

		return "", err
	}

	var response models.Workspaces

	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Error("failed to read response body from airbyte")

		return "", err
	}

	//Our Airbyte environment will have only one workspace
	if len(response.Workspaces) == 0 {
		err = errors.New("no workspace retrieved from airbyte")
		logger.Error(err.Error())

		return "", err
	}

	workspace := response.Workspaces[0]

	return workspace.WorkspaceId, nil
}

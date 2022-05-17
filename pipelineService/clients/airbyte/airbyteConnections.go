package airbyte

import (
	"bytes"
	"encoding/json"
	"fmt"

	"pipelineService/env"
	"pipelineService/models/v1"
	"pipelineService/utils"
)

func (airByteClient *RequestMaker) GetConnectionDetails(requestBody map[string]interface{}) (models.ConnectionMeta, error) {
	logger := utils.GetLogger()
	logger.Info("GetConnectionDetails from airbyte endpoint called")

	var response models.ConnectionMeta

	body, err := airByteClient.GetConnection(requestBody)

	if err != nil {
		logger.Error("get connection failed")

		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Error("failed to read response body from airbyte")

		return response, err
	}

	return response, nil
}

func (airByteClient *RequestMaker) GetConnection(requestBody map[string]interface{}) ([]byte, error) {
	logger := utils.GetLogger()
	logger.Info("GetConnection from airbyte endpoint called")

	airByteURL := fmt.Sprintf("%s/api/v1/web_backend/connections/get", env.Env.AirByteAddress)

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		logger.Error("failed to convert request body to json")

		return nil, err
	}

	body, err := airByteClient.sendRequest(airByteURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("failed to read response body from airbyte")

		return nil, err
	}

	return body, nil
}

func (airByteClient *RequestMaker) CreateConnection(
	requestBody models.CreatePipelineAirbyteRequest) (models.CreatePipelineAirbyteResponse, error) {
	logger := utils.GetLogger()

	airByteURL := fmt.Sprintf("%s/api/v1/web_backend/connections/create", env.Env.AirByteAddress)

	var response models.CreatePipelineAirbyteResponse

	jsonData, err := json.Marshal(requestBody)
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

func (airByteClient *RequestMaker) UpdateConnection(
	requestBody models.UpdatePipelineAirByteRequest) (models.CreatePipelineAirbyteResponse, error) {
	logger := utils.GetLogger()

	airByteURL := fmt.Sprintf("%s/api/v1/web_backend/connections/update", env.Env.AirByteAddress)

	var response models.CreatePipelineAirbyteResponse

	jsonData, err := json.Marshal(requestBody)
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

func (airByteClient *RequestMaker) SyncConnectionManually(requestBody map[string]interface{}) (models.ManualConnectionSyncResponse, error) {
	logger := utils.GetLogger()
	logger.Info("SyncConnectionManually from airbyte endpoint called")

	airByteURL := fmt.Sprintf("%s/api/v1/connections/sync", env.Env.AirByteAddress)

	var response models.ManualConnectionSyncResponse

	jsonData, err := json.Marshal(requestBody)
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

func (airByteClient *RequestMaker) FetchSyncHistory(request models.SyncHistoryRequest) (models.SyncHistoryResponse, error) {
	logger := utils.GetLogger()
	logger.Info("FetchSyncHistory from airByte endpoint called")

	airByteURL := fmt.Sprintf("%s/api/v1/jobs/list", env.Env.AirByteAddress)

	var response models.SyncHistoryResponse

	jsonData, err := json.Marshal(request)
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

func (airByteClient *RequestMaker) GetJobLogs(jobID int) (models.JobLogs, error) {
	logger := utils.GetLogger()
	logger.Info("GetJobLogs from airByte endpoint called")

	airByteURL := fmt.Sprintf("%s/api/v1/jobs/get", env.Env.AirByteAddress)

	var response models.JobLogs

	requestBody := map[string]interface{}{
		"id": jobID,
	}

	jsonData, err := json.Marshal(requestBody)
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

func (airByteClient *RequestMaker) GetConnectionSchema(connectionID string) (models.ConnectionSourceSchema, error) {
	logger := utils.GetLogger()
	logger.Info("GetConnectionSchema from airByte endpoint called")

	airByteURL := fmt.Sprintf("%s/api/v1/web_backend/connections/get", env.Env.AirByteAddress)

	var response models.ConnectionSourceSchema

	requestBody := map[string]interface{}{
		"connectionId":         connectionID,
		"withRefreshedCatalog": false,
	}

	jsonData, err := json.Marshal(requestBody)
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

func (airByteClient *RequestMaker) GetConnectionSummary(connectionID string) (models.ConnectionSummaryAirByte, error) {
	logger := utils.GetLogger()
	logger.Info("GetConnectionSummary from airByte endpoint called")

	airByteURL := fmt.Sprintf("%s/api/v1/connections/get", env.Env.AirByteAddress)

	var response models.ConnectionSummaryAirByte

	requestBody := map[string]interface{}{
		"connectionId": connectionID,
	}

	jsonData, err := json.Marshal(requestBody)
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

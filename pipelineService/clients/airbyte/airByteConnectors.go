package airbyte

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"pipelineService/env"
	"pipelineService/models/v1"
	"pipelineService/utils"
)

func (airByteClient *RequestMaker) CreateDestinationConnectorOnAirByte(
	requestBody models.CreateDestinationConnectorRequestAirbyte) (models.CreateDestinationConnectorResponseAirbyte, error) {
	logger := utils.GetLogger()

	airByteURL := fmt.Sprintf("%s/api/v1/destinations/create", env.Env.AirByteAddress)

	var response models.CreateDestinationConnectorResponseAirbyte

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

func (airByteClient *RequestMaker) GetDestinationDefinitions() (models.DestinationDefinitions, error) {
	logger := utils.GetLogger()

	airByteURL := fmt.Sprintf("%s/api/v1/destination_definitions/list", env.Env.AirByteAddress)

	var response models.DestinationDefinitions

	body, err := airByteClient.sendRequest(airByteURL, nil)
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

func (airByteClient *RequestMaker) GetDestinationSpecification(destinationDefinitionID string) (models.DestinationSpecification, error) {
	logger := utils.GetLogger()

	airByteURL := fmt.Sprintf("%s/api/v1/destination_definition_specifications/get", env.Env.AirByteAddress)

	requestBody := map[string]string{
		"destinationDefinitionId": destinationDefinitionID,
	}

	var response models.DestinationSpecification

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

func (airByteClient *RequestMaker) CreateSourceConnectorOnAirByte(
	requestBody models.CreateSourceConnectorRequestAirbyte) (models.CreateSourceConnectorResponseAirbyte, error) {
	logger := utils.GetLogger()

	airByteURL := fmt.Sprintf("%s/api/v1/sources/create", env.Env.AirByteAddress)

	var response models.CreateSourceConnectorResponseAirbyte

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

func (airByteClient *RequestMaker) EditSourceConnectorOnAirByte(requestBody models.EditSourceConnectorRequestAirByte) (models.CreateSourceConnectorResponseAirbyte, error) {
	logger := utils.GetLogger()

	airByteURL := fmt.Sprintf("%s/api/v1/sources/update", env.Env.AirByteAddress)

	var response models.CreateSourceConnectorResponseAirbyte

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		logger.Error("failed to convert request body to json")

		return response, err
	}

	body, err := airByteClient.sendRequest(airByteURL, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Error("failed to read response body from airByte")

		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Error("failed to read response body from airByte")

		return response, err
	}

	return response, nil
}

func (airByteClient *RequestMaker) GetSourceDefinitions() (models.SourceDefinitions, error) {
	logger := utils.GetLogger()

	airByteURL := fmt.Sprintf("%s/api/v1/source_definitions/list", env.Env.AirByteAddress)

	var response models.SourceDefinitions

	body, err := airByteClient.sendRequest(airByteURL, nil)
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

func (airByteClient *RequestMaker) GetConfiguredSource(sourceId string) (models.ConfiguredSource, error) {
	logger := utils.GetLogger()

	airByteURL := fmt.Sprintf("%s/api/v1/sources/get", env.Env.AirByteAddress)

	requestBody := map[string]string{
		"sourceId": sourceId,
	}

	var response models.ConfiguredSource

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

func (airByteClient *RequestMaker) GetSourceSpecification(sourceDefinitionID string) (models.SourceSpecification, error) {
	logger := utils.GetLogger()

	airByteURL := fmt.Sprintf("%s/api/v1/source_definition_specifications/get", env.Env.AirByteAddress)

	requestBody := map[string]string{
		"sourceDefinitionId": sourceDefinitionID,
	}

	var response models.SourceSpecification

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

func (airByteClient *RequestMaker) DiscoverSourceSchema(sourceId string) (models.SourceSchema, error) {
	logger := utils.GetLogger()

	airByteURL := fmt.Sprintf("%s/api/v1/sources/discover_schema", env.Env.AirByteAddress)

	requestBody := map[string]string{
		"sourceId": sourceId,
	}

	var response models.SourceSchema

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

func (airByteClient *RequestMaker) CheckDestinationConnection(requestBody map[string]interface{}) error {
	logger := utils.GetLogger()
	logger.Info("CheckDestinationConnection on AirByte called")

	airByteURL := fmt.Sprintf("%s/api/v1/scheduler/destinations/check_connection", env.Env.AirByteAddress)

	var response models.CheckConnection

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	body, err := airByteClient.sendRequest(airByteURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Status == "failed" {
		return errors.New(response.Message)
	}

	return nil
}

func (airByteClient *RequestMaker) CheckSourceConnection(requestBody map[string]interface{}) error {
	logger := utils.GetLogger()
	logger.Info("CheckSourceConnection on AirByte called")

	airByteURL := fmt.Sprintf("%s/api/v1/scheduler/sources/check_connection", env.Env.AirByteAddress)

	var response models.CheckConnection

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	body, err := airByteClient.sendRequest(airByteURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if response.Status == "failed" {
		return errors.New(response.Message)
	}

	return nil
}

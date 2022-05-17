package authService

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"pipelineService/env"
	"pipelineService/models/v1"
	"pipelineService/utils"
)

func (authServiceCient *RequestMaker) GetUserByID(ownerID int) (models.UserDetails, error) {
	logger := utils.GetLogger()
	logger.Info("GetUserByID from authService endpoint called")

	authServiceURL := fmt.Sprintf("%s/auth-service/api/v1/accounts/internal/user-from-id?user_id=%d", env.Env.AuthServiceAddress, ownerID)

	var response models.UserDetails

	res, err := authServiceCient.client.Get(authServiceURL)
	if err != nil {
		logger.Error("request to authService failed")

		return response, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error("failed to read response body from authService")

		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Error("failed to read response body from authService")

		return response, err
	}

	return response, nil
}

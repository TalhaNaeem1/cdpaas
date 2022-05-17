package authService

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"pipelineService/env"
	"pipelineService/utils"
)

func (authServiceClient *RequestMaker) ValidateSession(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("Middleware to validate sessionID called")

	sessionID, err := ctx.Cookie("sessionid")
	if err != nil {
		logger.Error(err.Error())

		msg := "sessionid cookie not found"

		utils.BuildResponseAndAbort(ctx, http.StatusUnauthorized, utils.ERROR, msg, nil)

		return
	}

	authServiceURL := fmt.Sprintf(
		"%s/auth-service/api/v1/accounts/user-info/?session_id=%s",
		env.Env.AuthServiceAddress,
		sessionID)

	res, err := authServiceClient.client.Get(authServiceURL)

	if err != nil {
		logger.Error(err.Error())

		msg := "request to auth-service failed"

		utils.BuildResponseAndAbort(ctx, http.StatusUnauthorized, utils.ERROR, msg, nil)

		return
	} else if res.StatusCode != http.StatusOK {
		logger.Error(fmt.Sprintf("validate session request failed with status code: %d", res.StatusCode))
		msg := "session validation failed"
		utils.BuildResponseAndAbort(ctx, http.StatusUnauthorized, utils.ERROR, msg, nil)

		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Error(err.Error())

		msg := "failed to read response"

		utils.BuildResponseAndAbort(ctx, http.StatusUnauthorized, utils.ERROR, msg, nil)

		return
	}

	var response UserInfo

	err = json.Unmarshal(body, &response)
	if err != nil {
		logger.Error(err.Error())

		msg := "failed to read response body from auth-service"

		utils.BuildResponseAndAbort(ctx, http.StatusUnauthorized, utils.ERROR, msg, nil)

		return
	}

	ctx.Set("userID", response.Payload.User.Id)
	ctx.Set("workspaceID", response.Payload.User.Workspace.Id)
	ctx.Set("airbyteWorkspaceID", response.Payload.User.Workspace.AirbyteWorkspaceId)

	logger.Info("session validation successful")
}

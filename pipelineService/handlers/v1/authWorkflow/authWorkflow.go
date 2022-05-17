package authWorkflow

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/cadence/client"
	"pipelineService/clients/cadenceClient"
	"pipelineService/env"
	"pipelineService/models/v1"
	"pipelineService/utils"
)

type Server struct {
	Router        *gin.Engine
	RouterGroup   *gin.RouterGroup
	CadenceClient cadenceClient.CadStore
}

// EmailPin Sends verification pin via email
// @Summary Sends Verification pin
// @Description Sends verification pin via email
// @Tags auth-workflows/internal
// @Produce  json
// @Param email body models.EmailTemplate true "Receiver Credentials"
// @Success 200 {object} models.Response
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /auth-workflows/internal/pin/ [post].
func (server *Server) EmailPin(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("EmailPin endpoint called")

	var receiver models.EmailTemplate
	if err := ctx.ShouldBindJSON(&receiver); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	wID, _ := uuid.NewV1()

	workflowOptions := client.StartWorkflowOptions{
		ID:                              wID.String(),
		TaskList:                        env.Env.TaskListName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	c, cancel := context.WithTimeout(context.Background(), 20*time.Second)

	defer cancel()

	err := server.CadenceClient.TriggerSendEmailWorkflow(c, receiver, workflowOptions)

	if err != nil {
		msg := "couldn't Trigger Send Email Workflow"

		logger.Error(err.Error())

		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, msg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", nil)

	logger.Info("EmailPin is in progress")
}

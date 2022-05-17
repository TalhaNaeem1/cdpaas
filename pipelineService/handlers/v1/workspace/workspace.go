package workspace

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pipelineService/clients/airbyte"
	"pipelineService/clients/authService"
	"pipelineService/models/v1"
	"pipelineService/services/db"
	"pipelineService/utils"
)

type Server struct {
	Store       db.Store
	Router      *gin.Engine
	RouterGroup *gin.RouterGroup
	Airbyte     airbyte.AirByteQuerier
	AuthService authService.AuthServiceQuerier
}

// CreateWorkspaceOnAirByte creates a workspace on AirByte
// @Summary Create Workspace on AirByte and returns the workspace_id.
// @Description Creates a workspace on AirByte
// @Tags workspaces/internal
// @Produce  json
// @Param createWorkspace body models.WorkspaceRequest true "Workspace Info"
// @Success 201 {object} models.WorkspaceResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /workspaces/internal/ [post].
func (server *Server) CreateWorkspaceOnAirByte(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("CreateWorkspaceOnAirByte endpoint called")


	var workspace models.WorkspaceRequest
	if err := ctx.ShouldBindJSON(&workspace); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}


	createdWorkspace, err := server.Airbyte.CreateWorkspace(workspace)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusInternalServerError, utils.ERROR, err.Error(), nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusCreated, utils.SUCCESS, "", createdWorkspace)
	logger.Info("CreateWorkspaceOnAirByte endpoint returned")
}



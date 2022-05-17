package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pipelineService/services/db"
	"pipelineService/utils"
)

type Server struct {
	Store       db.Store
	Router      *gin.Engine
	RouterGroup *gin.RouterGroup
}

// GetHealth returns health of the server
// @Summary Get Health
// @Description Get health of the server which tell either server is up or down
// @Tags health
// @Produce  json
// @Success 200
// @Failure 400
// @Failure 500
// @Router /health/ [get].
func (server *Server) GetHealth(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetHealth endpoint called")

	status := "UP"
	utils.BuildResponse(ctx, http.StatusOK, status, "", nil)
	logger.Info("GetHealth endpoint returned successfully")
}

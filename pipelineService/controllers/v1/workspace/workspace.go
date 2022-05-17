package workspace

import (
	"github.com/gin-gonic/gin"
	"pipelineService/clients/airbyte"
	"pipelineService/clients/authService"
	"pipelineService/handlers/v1/workspace"
	"pipelineService/services/db"
)

func registerRoutes(server *workspace.Server) {
	workSpaceRoutes := server.RouterGroup.Group("workspaces/internal")
	{
		workSpaceRoutes.POST("/", server.CreateWorkspaceOnAirByte)
	}
}
func CreateNewServer(dbStore db.Store, airbyteClient airbyte.AirByteClient,
	authServiceClient authService.AuthServiceClient, router *gin.Engine, rg *gin.RouterGroup) {
	server := &workspace.Server{
		Store:       dbStore,
		Router:      router,
		RouterGroup: rg,
		Airbyte:     airbyteClient,
		AuthService: authServiceClient,
	}
	registerRoutes(server)
}

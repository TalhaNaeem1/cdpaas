package authWorkflow

import (
	"github.com/gin-gonic/gin"
	"pipelineService/clients/cadenceClient"
	"pipelineService/handlers/v1/authWorkflow"
)

func registerRoutes(server *authWorkflow.Server) {
	dataProductRoutes := server.RouterGroup.Group("auth-workflows")
	{
		dataProductRoutes.POST("internal/pin/", server.EmailPin)
	}
}

func CreateNewServer(router *gin.Engine, rg *gin.RouterGroup, cc cadenceClient.CadStore) {
	server := &authWorkflow.Server{
		Router:        router,
		RouterGroup:   rg,
		CadenceClient: cc,
	}
	registerRoutes(server)
}

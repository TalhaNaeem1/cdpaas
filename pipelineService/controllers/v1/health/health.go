package health

import (
	"github.com/gin-gonic/gin"
	"pipelineService/handlers/v1/health"
	"pipelineService/services/db"
)

func registerRoutes(server *health.Server) {
	healthRoutes := server.RouterGroup.Group("health")
	{
		healthRoutes.GET("/", server.GetHealth)
	}
}

func CreateNewServer(dbStore db.Store, router *gin.Engine, rg *gin.RouterGroup) {
	server := &health.Server{
		Store:       dbStore,
		Router:      router,
		RouterGroup: rg,
	}
	registerRoutes(server)
}

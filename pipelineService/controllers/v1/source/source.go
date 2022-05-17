package source

import (
	"github.com/gin-gonic/gin"
	"pipelineService/clients/airbyte"
	"pipelineService/clients/authService"
	"pipelineService/handlers/v1/source"
	"pipelineService/services/db"
)

func registerRoutes(server *source.Server) {
	sourceRoutes := server.RouterGroup.Group("sources")
	{
		sourceRoutes.POST("/", server.ConfigureSourceOnAirbyte)
		sourceRoutes.PUT("/:id/", server.EditSourceOnAirByte)
		sourceRoutes.GET("/", server.GetSupportedSources)
		sourceRoutes.GET("/:id/", server.GetConfiguredSource)
		sourceRoutes.GET("/:id/summary/", server.GetConnectionSummary)
		sourceRoutes.GET("/specification/", server.GetSourceSpecification)
		sourceRoutes.GET("/discover/schema/", server.DiscoverSourceSchema)
	}
}
func CreateNewServer(dbStore db.Store, airbyteClient airbyte.AirByteClient,
	authServiceClient authService.AuthServiceClient, router *gin.Engine, rg *gin.RouterGroup) {
	server := &source.Server{
		Store:       dbStore,
		Router:      router,
		RouterGroup: rg,
		Airbyte:     airbyteClient,
		AuthService: authServiceClient,
	}
	registerRoutes(server)
}

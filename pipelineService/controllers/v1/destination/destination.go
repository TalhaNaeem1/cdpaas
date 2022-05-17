package destination

import (
	"github.com/gin-gonic/gin"
	"pipelineService/clients/airbyte"
	"pipelineService/clients/authService"
	"pipelineService/handlers/v1/destination"
	"pipelineService/services/db"
)

func registerRoutes(server *destination.Server) {
	destinationRoutes := server.RouterGroup.Group("destinations")
	{
		destinationRoutes.POST("/", server.ConfigureDestinationOnAirbyte)
		destinationRoutes.GET("/", server.GetSupportedDestinations)
		destinationRoutes.GET("/specification/", server.GetDestinationSpecification)
		destinationRoutes.GET("/configured/", server.GetConfiguredDestinations)
		destinationRoutes.GET("/:id/summary/", server.GetDestinationSummary)
	}
}
func CreateNewServer(dbStore db.Store, airbyteClient airbyte.AirByteClient,
	authServiceClient authService.AuthServiceClient, router *gin.Engine, rg *gin.RouterGroup) {
	server := &destination.Server{
		Store:       dbStore,
		Router:      router,
		RouterGroup: rg,
		Airbyte:     airbyteClient,
		AuthService: authServiceClient,
	}
	registerRoutes(server)
}

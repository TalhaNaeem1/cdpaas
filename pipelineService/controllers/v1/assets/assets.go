package assets

import (
	"github.com/gin-gonic/gin"
	"pipelineService/clients/airbyte"
	"pipelineService/clients/authService"
	"pipelineService/handlers/v1/assets"
	"pipelineService/services/db"
)

func registerRoutes(server *assets.Server) {
	assetRoutes := server.RouterGroup.Group("assets")
	{
		assetRoutes.GET("/:id/preview/", server.PreviewAsset)
		assetRoutes.GET("/:id/transformed/preview/", server.PreviewTransformedAsset)
		assetRoutes.GET("/pipeline/:id/", server.GetPipelineAssets)
		assetRoutes.GET("/products/:id/transformed/", server.GetTransformedAssets)
	}
}
func CreateNewServer(dbStore db.Store, airbyteClient airbyte.AirByteClient,
	authServiceClient authService.AuthServiceClient, router *gin.Engine, rg *gin.RouterGroup) {
	server := &assets.Server{
		Store:       dbStore,
		Router:      router,
		RouterGroup: rg,
		Airbyte:     airbyteClient,
		AuthService: authServiceClient,
	}
	registerRoutes(server)
}

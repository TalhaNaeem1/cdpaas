package dataProduct

import (
	"github.com/gin-gonic/gin"
	"pipelineService/clients/airbyte"
	"pipelineService/clients/authService"
	"pipelineService/clients/cadenceClient"
	"pipelineService/handlers/v1/dataProduct"
	"pipelineService/services/db"
)

func registerRoutes(server *dataProduct.Server) {
	dataProductRoutes := server.RouterGroup.Group("data-products")
	{
		dataProductRoutes.POST("/", server.CreateDataProduct)
		dataProductRoutes.GET("/", server.GetAllDataProducts)
		dataProductRoutes.GET("/:id/", server.GetDataProduct)
		dataProductRoutes.POST("/:id/add-pipeline/", server.AddPipeline)
		dataProductRoutes.PUT("/:id/", server.UpdateDataProduct)

		dataProductRoutes.POST("/transformations/:id/", server.ApplyTransformations)
		dataProductRoutes.GET("/transformations/:id/", server.GetTransformationDetails)
		dataProductRoutes.PUT("/transformations/:id/", server.UpdateTransformations)
	}

	dataProductRoutes = server.RouterGroup.Group("data-products/internal")
	{
		dataProductRoutes.GET("/", server.GetProductDetails)
		dataProductRoutes.POST("/transformations/assets/", server.SyncTransformedAssets)
	}
}

func CreateNewServer(dbStore db.Store, airbyteClient airbyte.AirByteClient, router *gin.Engine,
	authServiceClient authService.AuthServiceClient, rg *gin.RouterGroup, cc cadenceClient.CadStore) {
	server := &dataProduct.Server{
		Store:         dbStore,
		Router:        router,
		RouterGroup:   rg,
		Airbyte:       airbyteClient,
		CadenceClient: cc,
		AuthService:   authServiceClient,
	}
	registerRoutes(server)
}

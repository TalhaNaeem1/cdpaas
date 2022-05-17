package pipeline

import (
	"github.com/gin-gonic/gin"
	"pipelineService/clients/airbyte"
	"pipelineService/clients/authService"
	"pipelineService/clients/cadenceClient"
	"pipelineService/handlers/v1/pipeline"
	"pipelineService/services/db"
)

func registerRoutes(server *pipeline.Server) {
	pipelineRoutes := server.RouterGroup.Group("pipelines")
	{
		pipelineRoutes.POST("/", server.CreatePipeline)
		pipelineRoutes.PUT("/:id/", server.UpdatePipeline)
		pipelineRoutes.GET("/", server.GetAllPipelines)
		pipelineRoutes.GET("/:id/", server.GetPipeline)
		pipelineRoutes.DELETE("/:id/", server.TriggerDeletePipeline)
		pipelineRoutes.POST("/connections/", server.CreatePipelineConnection)
		pipelineRoutes.PUT("/connections/:id/", server.UpdatePipelineConnection)
		pipelineRoutes.GET("/connections/sync/logs/:job_id/", server.GetJobLogsFromAirByte)
		pipelineRoutes.GET("/connections/:connection_id/schema/", server.GetSourceSchemaFromAirByteConnection)
		pipelineRoutes.POST("/connections/:connection_id/sync/", server.RunManualSyncOnAirByte)
		pipelineRoutes.GET("/connections/:connection_id/sync/history/", server.FetchSyncHistoryFromAirByte)
	}

	pipelineRoutes = server.RouterGroup.Group("pipelines/internal")
	{
		pipelineRoutes.GET("/connections/", server.GetAllConnections)
		pipelineRoutes.PATCH("/connections/", server.UpdateConnections)
		pipelineRoutes.POST("/connections/", server.CreatePipelineConnectionOnAirbyte)
		pipelineRoutes.PUT("/connections/:id/", server.UpdatePipelineConnectionOnAirByte)

		pipelineRoutes.GET("/:id/", server.GetPipelineSourceAndConnectionID)
		pipelineRoutes.DELETE("/:id/", server.DeletePipeline)
		pipelineRoutes.PATCH("/:id/", server.UpdatePipelineStatus)

		pipelineRoutes.POST("/schema/", server.CreatePipelineSchema)
		pipelineRoutes.DELETE("/schema/:id/", server.DeletePipelineSchema)
		pipelineRoutes.GET("/:id/schema/", server.GetPipelineSchema)
		pipelineRoutes.POST("/pipeline_assets/", server.CreatePipelineAssets)
		pipelineRoutes.PATCH("/assets/enable/", server.EnablePipelineAssets)
	}
}

func CreateNewServer(dbStore db.Store, airbyteClient airbyte.AirByteClient,
	authServiceClient authService.AuthServiceClient, router *gin.Engine, rg *gin.RouterGroup, cc cadenceClient.CadStore) {
	server := &pipeline.Server{
		Store:         dbStore,
		Router:        router,
		RouterGroup:   rg,
		Airbyte:       airbyteClient,
		CadenceClient: cc,
		AuthService:   authServiceClient,
	}
	registerRoutes(server)
}

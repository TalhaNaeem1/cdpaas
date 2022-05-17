package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"pipelineService/clients/airbyte"
	"pipelineService/clients/authService"
	cadenceclient "pipelineService/clients/cadenceClient"
	"pipelineService/controllers/v1/assets"
	"pipelineService/controllers/v1/authWorkflow"
	"pipelineService/controllers/v1/dataProduct"
	"pipelineService/controllers/v1/destination"
	"pipelineService/controllers/v1/health"
	"pipelineService/controllers/v1/pipeline"
	"pipelineService/controllers/v1/source"
	"pipelineService/controllers/v1/workspace"
	"pipelineService/docs"
	"pipelineService/env"
	"pipelineService/services/db"
	"pipelineService/utils"
)

func main() {
	logger := utils.GetLogger()
	logger.Info("Starting Pipeline Service")
	setupSwaggerDocumentation()

	cadenceClient, err := cadenceclient.GetNewCadenceClient()
	if err != nil {
		logger.Error(err.Error())

		return
	}

	cadStore := cadenceclient.NewStore(&cadenceClient)

	database := db.GetConnection()
	dbStore := db.NewStore(database)

	router := gin.New()

	//router.Use(cors.Default())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		Output: utils.LogFileWriter,
	}))

	httpClient := http.DefaultClient
	airByteClient := airbyte.NewClient(httpClient)
	authServiceClient := authService.NewClient(httpClient)

	pipelineServiceGrp := router.Group("pipeline-service/api/v1")

	authWorkflow.CreateNewServer(router, pipelineServiceGrp, cadStore)
	health.CreateNewServer(dbStore, router, pipelineServiceGrp)
	dataProduct.CreateNewServer(dbStore, airByteClient, router, authServiceClient, pipelineServiceGrp, nil)
	pipeline.CreateNewServer(dbStore, airByteClient, authServiceClient, router, pipelineServiceGrp, cadStore)
	source.CreateNewServer(dbStore, airByteClient, authServiceClient, router, pipelineServiceGrp)
	destination.CreateNewServer(dbStore, airByteClient, authServiceClient, router, pipelineServiceGrp)
	workspace.CreateNewServer(dbStore, airByteClient, authServiceClient, router, pipelineServiceGrp)
	assets.CreateNewServer(dbStore, airByteClient, authServiceClient, router, pipelineServiceGrp)

	// register swagger documentation endpoint
	pipelineServiceGrp.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	httpServer := &http.Server{
		Addr:    ":" + env.Env.ServerPort,
		Handler: router,
	}

	startHttpServer(httpServer)
}

func setupSwaggerDocumentation() {
	docs.SwaggerInfo.Title = "CdPaas - Pipeline Service API Documentation"
	docs.SwaggerInfo.Description = "This swagger documentation contains API documentation for Pipeline Service" +
		"REST endpoints. Endpoints are categorized into two types i.e. External endpoints and Internal endpoints." +
		"External endpoints are the ones which will be consumed by External clients," +
		"whereas Internal endpoints are to be consumed by Internal services."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/pipeline-service/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http"}
}

func startHttpServer(server *http.Server) {
	logger := utils.GetLogger()
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Error(err.Error())
	}
}

package test

import (
	"github.com/gin-gonic/gin"
	mock_airbyte "pipelineService/clients/airbyte/mocks"
	"pipelineService/clients/authService"
	"pipelineService/controllers/v1/dataProduct"
	"pipelineService/controllers/v1/destination"
	"pipelineService/controllers/v1/health"
	"pipelineService/controllers/v1/pipeline"
	"pipelineService/controllers/v1/source"
	mock_store "pipelineService/services/db/mocks"
)

const BaseURL = "/pipeline-service/api/v1/"

type PackageName string

const (
	HEALTH       PackageName = "health"
	DATA_PRODUCT PackageName = "dataProduct"
	SOURCE       PackageName = "source"
	PIPELINE     PackageName = "pipeline"
	DESTINATION  PackageName = "destination"
)

// NewTestServer returns a router.
func NewTestServer(packageName PackageName, mockStore *mock_store.MockStore,
	mockAirByteClient *mock_airbyte.MockAirByteQuerier, AuthServiceClient authService.AuthServiceClient) *gin.Engine {
	router := gin.New()
	pipelineServiceGrp := router.Group(BaseURL)

	switch packageName {
	case HEALTH:
		health.CreateNewServer(mockStore, router, pipelineServiceGrp)

		return router

	case DATA_PRODUCT:
		dataProduct.CreateNewServer(mockStore, mockAirByteClient, router, AuthServiceClient, pipelineServiceGrp, nil)

		return router

	case SOURCE:
		source.CreateNewServer(mockStore, mockAirByteClient, AuthServiceClient, router, pipelineServiceGrp)

		return router

	case PIPELINE:
		pipeline.CreateNewServer(mockStore, mockAirByteClient, AuthServiceClient, router, pipelineServiceGrp, nil)

		return router

	case DESTINATION:
		destination.CreateNewServer(mockStore, mockAirByteClient, AuthServiceClient, router, pipelineServiceGrp)

		return router
	}

	return nil
}

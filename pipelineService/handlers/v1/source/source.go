package source

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"pipelineService/clients/airbyte"
	"pipelineService/clients/authService"
	"pipelineService/models/v1"
	"pipelineService/services/db"
	"pipelineService/utils"
)

type Server struct {
	Store       db.Store
	Router      *gin.Engine
	RouterGroup *gin.RouterGroup
	Airbyte     airbyte.AirByteQuerier
	AuthService authService.AuthServiceQuerier
}

// ConfigureSourceOnAirbyte configures and creates a source connector on Airbyte
// @Summary Create Source on Airbyte and store related information in local database. A connection is also created with the source in the database.
// @Description Configures and creates a source connector on Airbyte
// @Tags source
// @Produce  json
// @Param configureSourceData body models.CreateSourceConnectorRequestAPI true "Source Details"
// @Success 201 {object} models.CreateSourceConnectorResponseAPI
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /sources/ [post].
func (server *Server) ConfigureSourceOnAirbyte(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("ConfigureSourceOnAirbyte endpoint called")

	userID, workspaceID, airbyteWorkspaceID := utils.GetUserAndWorkspaceIDFromContext(ctx)

	var configureSourceData models.CreateSourceConnectorRequestAPI
	if err := ctx.ShouldBindJSON(&configureSourceData); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	var pipelineConnection models.PipelineConnection

	pipelineConnection, err := server.Store.GetPipelineConnection(configureSourceData.Pipeline)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Source and Connection creation")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	if pipelineConnection == (models.PipelineConnection{}) {
		requestBody := map[string]interface{}{
			"sourceDefinitionId":      configureSourceData.AirbyteSourceDefinitionId,
			"connectionConfiguration": configureSourceData.ConnectionConfiguration,
		}

		err = server.Airbyte.CheckSourceConnection(requestBody)
		if err != nil {
			logger.Error(err.Error())
			utils.BuildResponse(ctx, http.StatusInternalServerError, utils.ERROR, err.Error(), nil)

			return
		}

		airbyteRequest := models.CreateSourceConnectorRequestAirbyte{
			WorkspaceId:                  airbyteWorkspaceID,
			CreateSourceConnectorRequest: configureSourceData.CreateSourceConnectorRequest,
		}

		createSourceResponse, err := server.Airbyte.CreateSourceConnectorOnAirByte(airbyteRequest)
		if err != nil {
			logger.Error(err.Error())
			utils.BuildResponse(ctx, http.StatusInternalServerError, utils.ERROR, err.Error(), nil)

			return
		}

		createdSource := models.Source{
			SourceName:                createSourceResponse.SourceName,
			AirbyteSourceID:           createSourceResponse.AirbyteSourceId,
			AirbyteSourceDefinitionID: createSourceResponse.AirbyteSourceDefinitionId,
			Owner:                     userID,
			WorkspaceID:               workspaceID,
		}
		createdConnection := models.Connection{
			PipelineID: configureSourceData.Pipeline,
		}

		source, connection, err := server.Store.CreateConnectionAndSourceAgainstAPipeline(createdSource, createdConnection)
		if err != nil {
			logger.Error(err.Error())
			statusCode, errMsg := utils.ParseDBError(err, "Source and Connection creation")
			utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

			return
		}

		pipelineConnection.SourceID = source.SourceID
		pipelineConnection.SourceName = source.SourceName
		pipelineConnection.ConnectionID = connection.ConnectionID
		pipelineConnection.PipelineID = connection.PipelineID
	}

	utils.BuildResponse(ctx, http.StatusCreated, utils.SUCCESS, "", pipelineConnection)
	logger.Info("ConfigureSourceOnAirbyte endpoint returned")
}

// EditSourceOnAirByte Edits a source connector on AirByte
// @Summary Edit Source on AirByte
// @Description Edit a source connector on AirByte
// @Tags source
// @Produce  json
// @Param id path string true "Source ID"
// @Param EditSourceData body models.EditSourceConnectorRequest true "Source Details"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /sources/{id}/ [put].
func (server *Server) EditSourceOnAirByte(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("EditSourceOnAirByte endpoint called")

	sourceID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	var source models.Source

	source, err = server.Store.GetSource(sourceID.String())
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Source")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	var editSourceData models.EditSourceConnectorRequest
	if err := ctx.ShouldBindJSON(&editSourceData); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	requestBody := map[string]interface{}{
		"sourceDefinitionId":      source.AirbyteSourceDefinitionID,
		"connectionConfiguration": editSourceData.ConnectionConfiguration,
	}

	err = server.Airbyte.CheckSourceConnection(requestBody)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusInternalServerError, utils.ERROR, err.Error(), nil)

		return
	}

	editSourceDataAirByte := models.EditSourceConnectorRequestAirByte{
		AirByteSourceID:         source.AirbyteSourceID,
		ConnectionConfiguration: editSourceData.ConnectionConfiguration,
		Name:                    source.SourceName,
	}

	_, err = server.Airbyte.EditSourceConnectorOnAirByte(editSourceDataAirByte)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusInternalServerError, utils.ERROR, err.Error(), nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", "Source Edited Successfully")

	logger.Info("GetSupportedSources endpoint returned")
}

// GetSupportedSources return all the sources supported by cdpaas
// @Summary Get All Supported Sources
// @Description Return all the sources supported by cdpaas
// @Tags source
// @Produce  json
// @Success 200 {object} models.SupportedSourcesResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /sources/ [get].
func (server *Server) GetSupportedSources(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetSupportedSources endpoint called")

	sources, err := server.Store.GetSupportedSources()
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Supported Sources")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", sources)

	logger.Info("GetSupportedSources endpoint returned")
}

// GetConfiguredSource Gets a source connector details from AirByte
// @Summary Get Source from AirByte
// @Description Get a Source from AirByte
// @Tags source
// @Produce  json
// @Param id path string true "Source ID"
// @Success 200 {object} models.ConfiguredSourceResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /sources/{id}/ [get].
func (server *Server) GetConfiguredSource(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetConfiguredSource endpoint called")

	var configuredSource models.ConfiguredSource

	sourceID := ctx.Param("id")

	source, err := server.Store.GetSource(sourceID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Source")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	if source.AirbyteSourceID == "" {
		errMsg := "missing AirByte SourceID"
		logger.Error(errMsg)
		utils.BuildResponse(ctx, http.StatusInternalServerError, utils.ERROR, errMsg, nil)

		return
	}

	configuredSource, err = server.Airbyte.GetConfiguredSource(source.AirbyteSourceID)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusInternalServerError, utils.ERROR, err.Error(), nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", configuredSource)

	logger.Info("GetConfiguredSource endpoint returned")
}

// GetSourceSpecification return the specification of a given source
// @Summary Get Source Specification
// @Description Return the specification of a given source
// @Tags source
// @Produce  json
// @Param source query string true "Source Name"
// @Success 200 {object} models.SourceSpecification
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /sources/specification/ [get].
func (server *Server) GetSourceSpecification(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetSourceSpecification endpoint called")

	sourceName := ctx.Query("source")

	logger.Info("Fetching Source Definitions from AirByte")

	sourceDefinitions, err := server.Airbyte.GetSourceDefinitions()
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Source Definitions")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	var sourceDefinitionID string

	for _, sourceDefinition := range sourceDefinitions.SourceDefinitions {
		if sourceDefinition.Name == sourceName {
			sourceDefinitionID = sourceDefinition.SourceDefinitionID

			break
		}
	}

	logger.Info(fmt.Sprintf("Fetching Source Specification from AirByte for source: %s", sourceName))

	sourceSpecification, err := server.Airbyte.GetSourceSpecification(sourceDefinitionID)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", sourceSpecification)

	logger.Info("GetSourceSpecification endpoint returned")
}

// DiscoverSourceSchema discover and return the source schema from AirByte
// @Summary Discover Source Schema
// @Description Discover and Return the source schema from AirByte
// @Tags source
// @Produce  json
// @Param source_id query string true "Source ID"
// @Success 200 {object} models.SourceSchemaResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /sources/discover/schema/ [get].
func (server *Server) DiscoverSourceSchema(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("DiscoverSourceSchema endpoint called")

	sourceID := ctx.Query("source_id")

	var source models.Source
	source, err := server.Store.GetSource(sourceID)

	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Source Schema")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	sourceSchema, err := server.Airbyte.DiscoverSourceSchema(source.AirbyteSourceID)
	if err != nil {
		logger.Error(err.Error())

		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, errors.New("couldn't get discover source schema").Error(), nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", sourceSchema)

	logger.Info("DiscoverSourceSchema endpoint returned")
}

// GetConnectionSummary Returns Summary of a connection
// @Summary Returns Summary of a connection
// @Description Returns Summary of a connection
// @Tags source
// @Produce  json
// @Param id path string true "Source ID"
// @Success 200 {object} models.ConnectionSummaryResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /sources/{id}/summary/ [get].
func (server *Server) GetConnectionSummary(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetConnectionSummary endpoint called")

	sourceID := ctx.Param("id")

	connectionSummary, err := server.Store.GetSourceAndConnectionDetails(sourceID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Connection")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	connectionSummaryResponseAirByte, err := server.Airbyte.GetConnectionSummary(connectionSummary.AirbyteConnectionID)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	authResponse, err := server.AuthService.GetUserByID(connectionSummary.Owner)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	connectionSummaryResponse := models.ConnectionSummaryResponse{
		AirByteSummary:    connectionSummaryResponseAirByte,
		SourceName:        connectionSummary.SourceName,
		Owner:             authResponse.Payload.UserInfo,
		ConfigurationDate: connectionSummary.CreatedAt,
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", connectionSummaryResponse)
	logger.Info("GetConnectionSummary successfully returned")
}

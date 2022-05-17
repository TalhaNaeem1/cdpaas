package pipeline

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"go.uber.org/cadence/client"
	"pipelineService/clients/airbyte"
	"pipelineService/clients/authService"
	"pipelineService/clients/cadenceClient"
	"pipelineService/env"
	"pipelineService/models/v1"
	"pipelineService/services/db"
	"pipelineService/utils"
)

type Server struct {
	Store         db.Store
	Router        *gin.Engine
	RouterGroup   *gin.RouterGroup
	Airbyte       airbyte.AirByteQuerier
	CadenceClient cadenceClient.CadStore
	AuthService   authService.AuthServiceQuerier
}

var wg sync.WaitGroup

func (server *Server) getPipelinesData(pipeline models.PipelinesMetaData, pipelineCh chan models.PipelinesMetaData) {
	logger := utils.GetLogger()

	defer wg.Done()

	requestBody := make(map[string]interface{})

	var (
		connectionMeta models.ConnectionMeta
		authResponse   models.UserDetails
	)

	//parse the airbyte connectionID
	if pipeline.AirbyteConnectionID != "" {
		airbyteConnectionID, err := uuid.FromString(pipeline.AirbyteConnectionID)
		if err != nil {
			logger.Error(err.Error())
			pipelineCh <- pipeline

			return
		}

		requestBody["connectionId"] = airbyteConnectionID
		requestBody["withRefreshedCatalog"] = false

		connectionMeta, err = server.Airbyte.GetConnectionDetails(requestBody)
		if err != nil {
			logger.Error(err.Error())
			pipelineCh <- pipeline

			return
		}

		pipeline.AirbyteLastRun = connectionMeta.LatestSyncJobCreatedAt
		pipeline.AirbyteStatus = connectionMeta.LatestSyncJobStatus
	}

	authResponse, err := server.AuthService.GetUserByID(pipeline.OwnerID)
	if err != nil {
		logger.Error(err.Error())
		pipelineCh <- pipeline

		return
	}

	pipeline.Owner = authResponse.Payload.UserInfo
	pipelineCh <- pipeline
}

// CreatePipeline returns newly created pipeline
// @Summary Create Pipeline
// @Description Creates the pipeline and links it to the specified data product
// @Tags pipelines
// @Produce  json
// @Param pipeline body models.Pipeline true "Pipeline Info"
// @Success 201 {object} models.PipelineResponse
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/ [post].
func (server *Server) CreatePipeline(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("CreatePipeline endpoint called")

	userID, workspaceID, _ := utils.GetUserAndWorkspaceIDFromContext(ctx)

	var pipeline models.Pipeline
	if err := ctx.ShouldBindJSON(&pipeline); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	pipeline.Owner = userID
	pipeline.WorkspaceID = workspaceID

	newPipeline, err := server.Store.CreatePipeline(pipeline)

	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Pipeline")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusCreated, utils.SUCCESS, "", newPipeline)

	logger.Info("CreatePipeline endpoint returned successfully")
}

// UpdatePipeline updates already created pipeline
// @Summary Updates a Pipeline
// @Description Edits the pipeline name and governance
// @Tags pipelines
// @Produce  json
// @Param pipeline body models.UpdatePipeline true "Pipeline Info"
// @Param id path string true "Pipeline ID"
// @Success 200 {object} models.PipelineResponse
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/{id}/ [put].
func (server *Server) UpdatePipeline(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("UpdatePipeline endpoint called")

	pipelineID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	var updatePipeline models.UpdatePipeline
	if err := ctx.ShouldBindJSON(&updatePipeline); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	updatePipeline.PipelineID = pipelineID

	updatedPipeline, err := server.Store.UpdatePipeline(updatePipeline)

	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Pipeline")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", updatedPipeline)

	logger.Info("UpdatePipeline endpoint returned successfully")
}

// GetAllPipelines returns all the pipelines of a specific product
// @Summary Get pipelines
// @Description Get all the pipelines
// @Tags pipelines
// @Produce  json
// @Success 200 {object} models.PipelineMetaDataResponse
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/ [get].
func (server *Server) GetAllPipelines(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetAllPipelines endpoint called")

	_, workspaceID, _ := utils.GetUserAndWorkspaceIDFromContext(ctx)

	var pipelinesMetaData []models.PipelinesMetaData

	pipelinesMetaData, err := server.Store.GetAllPipelines(workspaceID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Pipeline")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	pipelineChannel := make(chan models.PipelinesMetaData, len(pipelinesMetaData))

	for i := range pipelinesMetaData {
		wg.Add(1)

		go server.getPipelinesData(pipelinesMetaData[i], pipelineChannel) //pass err channel

		pipelinesData := <-pipelineChannel
		pipelinesMetaData[i] = pipelinesData
	}

	wg.Wait()

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", pipelinesMetaData)
	logger.Info("GetAllPipeline endpoint returned successfully")
}

// GetPipeline returns a pipeline
// @Summary Returns a pipeline
// @Description Returns a pipeline by ID
// @Tags pipelines
// @Produce  json
// @Param id path string true "Pipeline ID"
// @Success 200 {object} models.GetPipelineDetailsResponse
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/{id} [get].
func (server *Server) GetPipeline(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetPipeline endpoint called")

	var pipeline models.GetPipelineDetails
	//Parse the pipelineID
	pipelineID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	pipeline.Pipeline, err = server.Store.GetPipeline(pipelineID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Pipeline")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	var (
		connectionMeta models.ConnectionMeta
		authResponse   models.UserDetails
	)

	//parse the airbyte connectionID
	airbyteConnectionID, err := uuid.FromString(pipeline.Pipeline.AirbyteConnectionID)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	requestBody := make(map[string]interface{})
	requestBody["connectionId"] = airbyteConnectionID
	requestBody["withRefreshedCatalog"] = false

	connectionMeta, err = server.Airbyte.GetConnectionDetails(requestBody)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	pipeline.Pipeline.AirbyteLastRun = connectionMeta.LatestSyncJobCreatedAt
	pipeline.Pipeline.AirbyteStatus = connectionMeta.LatestSyncJobStatus

	authResponse, err = server.AuthService.GetUserByID(pipeline.Pipeline.Owner)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	pipeline.Owner = authResponse.Payload.UserInfo
	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", pipeline)
	logger.Info("GetPipeline endpoint returned successfully")
}

// CreatePipelineConnection creates a pipeline connection on Airbyte
// @Summary Create pipeline on airbyte
// @Description Creates a pipeline on airbyte using the specified sources and destinations
// @Tags pipelines
// @Produce  json
// @Param pipeline body models.CreatePipelineRequest true "Pipeline Info"
// @Success 201 {object} models.CreatePipelineResponse
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/connections/ [post].
func (server *Server) CreatePipelineConnection(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("CreatePipelineConnection endpoint called")

	userID, workspaceID, airbyteWorkspaceID := utils.GetUserAndWorkspaceIDFromContext(ctx)

	var createPipelineRequest models.CreatePipelineRequest
	if err := ctx.ShouldBindJSON(&createPipelineRequest); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	connectionInfo, err := server.Store.GetSourceAndDestinationAirbyteInfo(createPipelineRequest.SourceID, createPipelineRequest.DestinationID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Source And Destination Info for Air Byte")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	workflowOptions := client.StartWorkflowOptions{
		//ID:                              wID.String(),
		TaskList:                        env.Env.TaskListName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	err = server.CadenceClient.TriggerCreateConnectionWorkflow(ctx, workflowOptions, createPipelineRequest, connectionInfo, userID, workspaceID, airbyteWorkspaceID)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, "Couldn't Trigger CreatePipelineConnection Workflow", nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusCreated, utils.SUCCESS, "", "connection created")
	logger.Info("CreatePipelineConnection endpoint returned successfully")
}

// CreatePipelineConnectionOnAirbyte creates a pipeline connection on Airbyte
// @Summary Create pipeline on airbyte
// @Description Creates a pipeline on airbyte using the specified sources and destinations
// @Tags pipelines/internal
// @Produce  json
// @Param pipeline body models.CreatePipelineRequest true "Pipeline Info"
// @Success 201 {object} models.CreatePipelineResponse
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/internal/connections/ [post].
func (server *Server) CreatePipelineConnectionOnAirbyte(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("CreatePipelineOnAirbyte endpoint called")

	userID, workspaceID, airbyteWorkspaceID := utils.GetUserAndWorkspaceIDFromContext(ctx)

	var (
		createPipelineRequest models.CreatePipelineRequest
		airbyteFrequencyUnit  int
		airbyteTimeUnit       string
	)

	if err := ctx.ShouldBindJSON(&createPipelineRequest); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	airbyteInfo, err := server.Store.GetSourceAndDestinationAirbyteInfo(createPipelineRequest.SourceID, createPipelineRequest.DestinationID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Source And Destination Info for Air Byte")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	for index := range createPipelineRequest.Operations {
		createPipelineRequest.Operations[index].WorkspaceId = airbyteWorkspaceID
	}

	createPipelineAirbyteRequest := CreatePipelineAirbyteRequestModel(airbyteInfo, createPipelineRequest)

	newAirByteConnection, err := server.Airbyte.CreateConnection(*createPipelineAirbyteRequest)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	if newAirByteConnection.Schedule != nil {
		airbyteFrequencyUnit = newAirByteConnection.Schedule.Units
		airbyteTimeUnit = newAirByteConnection.Schedule.TimeUnit
	}

	airByteConnectionInfo := models.Connection{
		ConnectionID:          airbyteInfo.ConnectionID,
		AirbyteConnectionID:   newAirByteConnection.ConnectionId,
		AirbyteStatus:         newAirByteConnection.Status,
		AirbyteFrequencyUnits: airbyteFrequencyUnit,
		AirbyteTimeUnit:       airbyteTimeUnit,
		Owner:                 userID,
		WorkspaceID:           workspaceID,
	}

	if err = server.Store.UpdateConnectionInfo(airByteConnectionInfo, airbyteInfo.DestinationID); err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Connection")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusCreated, utils.SUCCESS, "", airbyteInfo)
	logger.Info("CreatePipeline endpoint returned successfully")
}

// UpdatePipelineConnection updates a pipeline connection on AirByte
// @Summary Updates Pipeline on AirByte
// @Description Updates an pipeline on airByte
// @Tags pipelines
// @Produce  json
// @Param pipeline body models.UpdatePipelineAirByteRequest true "Pipeline Info"
// @Param id path string true "Connection ID"
// @Success 200 {object} models.CreatePipelineResponse
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/connections/{id}/ [put].
func (server *Server) UpdatePipelineConnection(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("UpdatePipelineConnection endpoint called")

	userID, workspaceID, airbyteWorkspaceID := utils.GetUserAndWorkspaceIDFromContext(ctx)

	var updatePipelineRequest models.UpdatePipelineAirByteRequest

	if err := ctx.ShouldBindJSON(&updatePipelineRequest); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	connection, err := server.Store.GetConnection(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Connection")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	workflowOptions := client.StartWorkflowOptions{
		//ID:                              wID.String(),
		TaskList:                        env.Env.TaskListName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	c, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = server.CadenceClient.TriggerUpdateConnectionWorkflow(c, workflowOptions, updatePipelineRequest, connection, userID, workspaceID, airbyteWorkspaceID)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, "Couldn't Trigger UpdatePipelineConnection Workflow", nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusCreated, utils.SUCCESS, "", "connection updated")
	logger.Info("UpdatePipelineConnection endpoint returned successfully")
}

// UpdatePipelineConnectionOnAirByte updates a pipeline connection on AirByte
// @Summary Updates Pipeline on AirByte
// @Description Updates an pipeline on airByte
// @Tags pipelines/internal
// @Produce  json
// @Param pipeline body models.UpdatePipelineAirByteRequest true "Pipeline Info"
// @Param id path string true "Connection ID"
// @Success 200 {object} models.CreatePipelineResponse
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/internal/connections/{id}/ [put].
func (server *Server) UpdatePipelineConnectionOnAirByte(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("UpdatePipelineConnectionOnAirByte endpoint called")

	_, _, airbyteWorkspaceID := utils.GetUserAndWorkspaceIDFromContext(ctx)

	connectionID := ctx.Param("id")

	connection, err := server.Store.GetConnection(connectionID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Connection")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	var (
		updatePipelineRequest models.UpdatePipelineAirByteRequest
		airbyteFrequencyUnit  int
		airbyteTimeUnit       string
	)

	if err := ctx.ShouldBindJSON(&updatePipelineRequest); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	updatePipelineRequest.ConnectionId = connection.AirbyteConnectionID

	for index := range updatePipelineRequest.Operations {
		updatePipelineRequest.Operations[index].WorkspaceId = airbyteWorkspaceID
	}

	updatedAirByteConnection, err := server.Airbyte.UpdateConnection(updatePipelineRequest)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	if updatedAirByteConnection.Schedule != nil {
		airbyteFrequencyUnit = updatedAirByteConnection.Schedule.Units
		airbyteTimeUnit = updatedAirByteConnection.Schedule.TimeUnit
	}

	airByteConnectionInfo := models.Connection{
		ConnectionID:          connectionID,
		AirbyteFrequencyUnits: airbyteFrequencyUnit,
		AirbyteTimeUnit:       airbyteTimeUnit,
		IsFirstRun:            false,
	}

	if err = server.Store.UpdateConnectionSchedule(airByteConnectionInfo); err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Connection")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", "pipeline updated successfully")
	logger.Info("UpdatePipelineConnectionOnAirByte endpoint returned successfully")
}

// RunManualSyncOnAirByte trigger the sync
// @Summary manually triggers the sync on airbyte
// @Description Triggers the sync basis the connection ID
// @Tags pipelines
// @Produce  json
// @Param connection_id path string true "Connection ID"
// @Success 200 {object} models.ManualConnectionSyncResponse
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/connections/{connection_id}/sync/ [post].
func (server *Server) RunManualSyncOnAirByte(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("RunManualSyncOnAirByte endpoint called")

	requestBody := make(map[string]interface{})
	requestBody["connectionId"] = ctx.Param("connection_id")

	manualConnectionSyncResponse, err := server.Airbyte.SyncConnectionManually(requestBody)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", manualConnectionSyncResponse)

	logger.Info("RunManualSyncOnAirByte successfully returned")
}

func CreatePipelineAirbyteRequestModel(airbyteInfo models.AirbyteSourceAndDestinations,
	createPipelineRequest models.CreatePipelineRequest) *models.CreatePipelineAirbyteRequest {
	return &models.CreatePipelineAirbyteRequest{
		DestinationId:       airbyteInfo.AirbyteDestinationID,
		SourceId:            airbyteInfo.AirbyteSourceID,
		NamespaceDefinition: utils.AIRBYTE_DEFAULT_NAMESPACE_DEFINITION,
		NamespaceFormat:     airbyteInfo.PipelineName,
		Prefix:              *createPipelineRequest.Prefix,
		Status:              utils.AIRBYTE_DEFAULT_STATUS,
		Schedule:            createPipelineRequest.Schedule,
		SyncCatalog:         createPipelineRequest.SyncCatalog,
		Operations:          createPipelineRequest.Operations,
	}
}

// GetAllConnections returns all connections
// @Summary Returns all connections
// @Description Returns all existing connections
// @Tags pipelines/internal
// @Produce  json
// @Success 200 {array} models.Connection
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/internal/connections/ [get].
func (server *Server) GetAllConnections(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetAllConnections internal endpoint called")

	connections, err := server.Store.GetAllConnections()
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Connection")
		ctx.JSON(statusCode, errMsg)

		return
	}

	ctx.JSON(http.StatusOK, connections)
	logger.Info("GetAllConnections internal endpoint successfully returned")
}

// FetchSyncHistoryFromAirByte Fetches the sync history from AirByte
// @Summary Fetch Sync History From AirByte
// @Description Fetch Sync History From AirByte for a specific source
// @Tags pipelines
// @Produce  json
// @Param connection_id path string true "Connection ID"
// @Success 200 {object} models.SyncHistoryResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /pipelines/connections/{connection_id}/sync/history/ [get].
func (server *Server) FetchSyncHistoryFromAirByte(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("FetchSyncHistoryFromAirByte endpoint called")

	airByteConnectionID := ctx.Param("connection_id")

	requestBody := models.SyncHistoryRequest{
		ConfigTypes: []string{
			utils.SYNC,
			utils.RESET_CONNECTION,
		},
		ConfigId: airByteConnectionID,
	}

	SyncHistoryResponse, err := server.Airbyte.FetchSyncHistory(requestBody)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", SyncHistoryResponse)

	logger.Info("FetchSyncHistoryFromAirByte successfully returned")
}

// GetJobLogsFromAirByte Fetches the logs of a job from AirByte
// @Summary Fetch Job logs From AirByte
// @Description Fetch logs of a connection from AirByte as per the job
// @Tags pipelines
// @Produce  json
// @Param job_id path string true "Job ID"
// @Success 200 {object} models.JobLogs
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /pipelines/connections/sync/logs/{job_id}/ [get].
func (server *Server) GetJobLogsFromAirByte(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetJobLogsFromAirByte endpoint called")

	ID := ctx.Param("job_id")
	jobID, _ := strconv.Atoi(ID)

	JobLogs, err := server.Airbyte.GetJobLogs(jobID)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", JobLogs)

	logger.Info("GetJobLogsFromFromAirByte successfully returned")
}

// GetSourceSchemaFromAirByteConnection Fetches the Source Schema from AirByte
// @Summary Fetch Source Schema from AirByte
// @Description Fetch Source Schema of a connection from AirByte
// @Tags pipelines
// @Produce  json
// @Param connection_id path string true "Connection ID"
// @Success 200 {object} models.ConnectionSourceSchema
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /pipelines/connections/{connection_id}/schema/ [get].
func (server *Server) GetSourceSchemaFromAirByteConnection(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetSourceSchemaFromAirByteConnection endpoint called")

	airByteConnectionID := ctx.Param("connection_id")

	connSourceSchema, err := server.Airbyte.GetConnectionSchema(airByteConnectionID)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", connSourceSchema)

	logger.Info("GetSourceSchemaFromAirByteConnection successfully returned")
}

// UpdateConnections update connections in Database(DB)
// @Summary update all given connections
// @Description updates all given connections in DB
// @Tags pipelines/internal
// @Produce  json
// @Param connections body []models.Connection true "Connection Instance"
// @Success 200 {object} models.Response
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/internal/connections/ [patch].
func (server *Server) UpdateConnections(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("UpdateConnections internal endpoint called")

	var connections []models.Connection
	if err := ctx.ShouldBindJSON(&connections); err != nil {
		logger.Error(err.Error())

		msg := "missing connections payload"
		utils.BuildResponse(ctx, http.StatusBadRequest, msg, utils.ERROR, nil)

		return
	}

	err := server.Store.UpdateConnections(connections)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Connection")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", nil)

	logger.Info("UpdateConnections internal endpoint successfully returned")
}

// DeletePipeline deletes the pipelines
// @Summary deletes a pipelines
// @Description deletes a pipeline by ID
// @Tags pipelines/internal
// @Produce  json
// @Param id path string true "Pipeline ID"
// @Success 204 {object} models.Response
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/internal/{id}/ [delete].
func (server *Server) DeletePipeline(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("DeletePipeline internal endpoint called")

	//Parse the pipelineID
	pipelineID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	err = server.Store.DeletePipeline(pipelineID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Pipeline")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", "Pipeline deleted successfully")
	logger.Info("DeletePipeline internal endpoint successfully returned")
}

// TriggerDeletePipeline triggers the deletion workflow
// @Summary pipeline deletion workflow is triggered
// @Description pipeline deletion workflow on temporal is triggered
// @Tags pipelines
// @Produce  json
// @Param id path string true "Pipeline ID"
// @Success 200 {object} models.Response
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/{id}/ [delete].
func (server *Server) TriggerDeletePipeline(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("TriggerDeletePipeline endpoint called")

	workflowOptions := client.StartWorkflowOptions{
		//ID:                              wID.String(),
		TaskList:                        env.Env.TaskListName,
		ExecutionStartToCloseTimeout:    time.Minute,
		DecisionTaskStartToCloseTimeout: time.Minute,
	}

	c, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err := server.CadenceClient.TriggerDeletePipelineWorkflow(c, workflowOptions, ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, "Couldn't Trigger DeletePipeline Workflow", nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", "Pipeline Deletion In-Progress")
	logger.Info("TriggerDeletePipeline successfully returned")
}

// GetPipelineSourceAndConnectionID returns a source and connection ID
// @Summary Returns a source and connection ID
// @Description Returns a source and connection ID by pipeline ID
// @Tags pipelines/internal
// @Produce  json
// @Param id path string true "Pipeline ID"
// @Success 200 {object} models.GetPipelineSourceAndConnectionID
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/internal/{id}/ [get].
func (server *Server) GetPipelineSourceAndConnectionID(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetPipelineSourceAndConnectionID internal endpoint called")

	//Parse the pipelineID
	pipelineID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	var pipeline models.PipelineSourceAndConnectionID

	pipeline, err = server.Store.GetPipelineSourceAndConnectionID(pipelineID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Pipeline")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", pipeline)
	logger.Info("GetPipelineSourceAndConnectionID internal endpoint successfully returned")
}

// UpdatePipelineStatus updates pipeline_status
// @Summary pipeline_status updated to Deletion In-Progress
// @Description updates the pipeline_status to Deletion In-Progress
// @Tags pipelines/internal
// @Produce  json
// @Param id path string true "Pipeline ID"
// @Success 200 {object} models.Response
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/internal/{id}/ [patch].
func (server *Server) UpdatePipelineStatus(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("UpdatePipelineStatus internal endpoint called")

	//Parse the pipelineID
	pipelineID, err := uuid.FromString(ctx.Param("id"))

	pipelineStatus := ctx.Query("status")

	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	err = server.Store.UpdatePipelineStatus(pipelineID, pipelineStatus)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Pipeline")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", "PipelineStatus Updated Successfully")
	logger.Info("UpdatePipelineStatus internal endpoint successfully returned")
}

// CreatePipelineSchema creates a schema for a pipeline
// @Summary Creates a schema for a pipeline
// @Description Creates a schema for a pipeline
// @Tags pipelines/internal
// @Produce  json
// @Param pipeline body models.PipelineSchemas true "Pipeline Info"
// @Success 201 {object} models.PipelineSchemas
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/internal/schema/ [post].
func (server *Server) CreatePipelineSchema(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("CreatePipelineSchema endpoint called")

	var pipelineSchema models.PipelineSchemas
	if err := ctx.ShouldBindJSON(&pipelineSchema); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	pipelineSchema, err := server.Store.CreatePipelineSchema(pipelineSchema)

	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Pipeline Schema")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusCreated, utils.SUCCESS, "", pipelineSchema)

	logger.Info("CreatePipelineSchema endpoint returned successfully")
}

// CreatePipelineAssets create pipeline in Database(DB)
// @Summary create pipeline in Database(DB)
// @Description create pipeline in Database(DB)
// @Tags pipelines/internal
// @Produce  json
// @Param pipelineAssets body []models.PipelineAssets true "Pipeline Assets"
// @Success 200 {object} models.Response
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/internal/pipeline_assets/ [post].
func (server *Server) CreatePipelineAssets(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("CreatePipelineAssets internal endpoint called")

	var pipelineAssets []models.PipelineAssets
	if err := ctx.ShouldBindJSON(&pipelineAssets); err != nil {
		logger.Error(err.Error())

		msg := "missing CreatePipelineAssets payload"
		utils.BuildResponse(ctx, http.StatusBadRequest, msg, utils.ERROR, nil)

		return
	}

	err := server.Store.CreatePipelineAssets(pipelineAssets)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Pipeline Assets")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", nil)

	logger.Info("CreatePipelineAssets internal endpoint successfully returned")
}

// DeletePipelineSchema deletes the pipeline schema and related assets
// @Summary deletes the pipeline schema and related assets
// @Description deletes the pipeline schema and related assets
// @Tags pipelines/internal
// @Produce  json
// @Param id path string true "Schema ID"
// @Success 204 {object} models.Response
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/internal/schema/{id}/ [delete].
func (server *Server) DeletePipelineSchema(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("DeletePipelineSchema internal endpoint called")

	pipelineSchemaID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	err = server.Store.DeletePipelineSchema(pipelineSchemaID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Pipeline")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", "Pipeline Schema deleted successfully")
	logger.Info("DeletePipelineSchema internal endpoint successfully returned")
}

// GetPipelineSchema returns the pipeline schema
// @Summary fetches the pipeline schema
// @Description fetches the pipeline schema from pipeline_id
// @Tags pipelines/internal
// @Produce  json
// @Param id path string true "Pipeline ID"
// @Success 204 {object} models.Response
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/internal/{id}/schema/ [get].
func (server *Server) GetPipelineSchema(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetPipelineSchema internal endpoint called")

	pipelineID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	pipelineSchema, err := server.Store.GetPipelineSchema(pipelineID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Pipeline")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", pipelineSchema)
	logger.Info("GetPipelineSchema internal endpoint successfully returned")
}

// EnablePipelineAssets enables the assets to be served
// @Summary is_enabled updated to true
// @Description updates the is_enabled to true
// @Tags pipelines/internal
// @Produce  json
// @Param connectionIDs body models.EnableAssetsInternalRequest true "connection Id's"
// @Success 200 {object} models.Response
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /pipelines/internal/assets/enable/ [patch].
func (server *Server) EnablePipelineAssets(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("ActivatePipelineAssets internal endpoint called")

	var connections models.EnableAssetsInternalRequest
	//Parse the connectionIDs
	err := ctx.ShouldBindJSON(&connections)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	err = server.Store.EnablePipelineAssets(connections.ConnectionIDs)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Pipeline")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", "Pipeline Assets Enabled successfully")
	logger.Info("ActivatePipelineAssets internal endpoint successfully returned")
}

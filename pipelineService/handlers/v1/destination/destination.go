package destination

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

// ConfigureDestinationOnAirbyte configures and creates a destination connector on Airbyte
// @Summary Create Destination on Airbyte and store related information in local database.
// @Description Configures and creates a destination connector on Airbyte
// @Tags destination
// @Param requestBody body models.CreateDestinationConnectorRequestAPI true "request body"
// @Produce  json
// @Success 201 {object} models.CreateDestinationConnectorResponseAPI
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /destinations/ [post].
func (server *Server) ConfigureDestinationOnAirbyte(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("ConfigureDestinationOnAirbyte endpoint called")

	userID, workspaceID, airbyteWorkspaceID := utils.GetUserAndWorkspaceIDFromContext(ctx)

	var configureDestinationData models.CreateDestinationConnectorRequestAPI
	if err := ctx.ShouldBindJSON(&configureDestinationData); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	requestBody := map[string]interface{}{
		"destinationDefinitionId": configureDestinationData.AirbyteDestinationDefinitionId,
		"connectionConfiguration": configureDestinationData.ConnectionConfiguration,
	}

	err := server.Airbyte.CheckDestinationConnection(requestBody)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusInternalServerError, utils.ERROR, err.Error(), nil)

		return
	}

	airbyteRequest := models.CreateDestinationConnectorRequestAirbyte{
		WorkspaceId:                       airbyteWorkspaceID,
		CreateDestinationConnectorRequest: configureDestinationData.CreateDestinationConnectorRequest,
	}

	createDestinationResponse, err := server.Airbyte.CreateDestinationConnectorOnAirByte(airbyteRequest)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusInternalServerError, utils.ERROR, err.Error(), nil)

		return
	}

	createdDestination := models.Destination{
		DestinationName:         createDestinationResponse.Name,
		AirbyteDestinationID:    createDestinationResponse.AirbyteDestinationId,
		AirbyteDestDefinitionID: createDestinationResponse.AirbyteDestinationDefinitionId,
		DestinationType:         configureDestinationData.DestinationType,
		ConfigurationDetails:    configureDestinationData.ConnectionConfiguration,
		Owner:                   userID,
		WorkspaceID:             workspaceID,
	}

	insertedDestination, err := server.Store.CreateDestination(createdDestination)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Destination")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	apiResponse := models.CreateDestinationConnectorResponseData{
		DestinationID:   insertedDestination.DestinationID,
		DestinationName: insertedDestination.DestinationName,
		DestinationType: insertedDestination.DestinationType,
	}

	utils.BuildResponse(ctx, http.StatusCreated, utils.SUCCESS, "", apiResponse)

	logger.Info("ConfigureDestinationOnAirbyte endpoint returned")
}

// GetSupportedDestinations return all the destinations supported by cdpaas
// @Summary Get All Supported Destinations
// @Description Return all the destinations supported by cdpaas
// @Tags destination
// @Produce  json
// @Success 200 {object} models.SupportedDestinationResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /destinations/ [get].
func (server *Server) GetSupportedDestinations(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetSupportedDestinations endpoint called")

	destinations, err := server.Store.GetSupportedDestinations()
	if err != nil {
		logger.Error(err.Error())

		statusCode, errMsg := utils.ParseDBError(err, "Supported Destination")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", destinations)

	logger.Info("GetSupportedDestinations endpoint returned")
}

// GetConfiguredDestinations return all the destinations that are configured on cdpaas
// @Summary Get All Configured Destinations
// @Description Return all the destinations that are configured on cdpaas
// @Tags destination
// @Produce  json
// @Success 200 {object} models.ConfiguredDestinationResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /destinations/configured/ [get].
func (server *Server) GetConfiguredDestinations(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetConfiguredDestinations endpoint called")

	_, workspaceId, _ := utils.GetUserAndWorkspaceIDFromContext(ctx)

	configuredDestinations, err := server.Store.GetConfiguredDestination(workspaceId)
	if err != nil {
		logger.Error(err.Error())

		statusCode, errMsg := utils.ParseDBError(err, "Configured Destination")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", configuredDestinations)

	logger.Info("GetConfiguredDestinations endpoint returned")
}

// GetDestinationSummary return destination's summary
// @Summary Get Destination Summary
// @Description Returns the destination's summary
// @Tags destination
// @Produce  json
// @Param id path string true "Destination ID"
// @Success 200 {object} models.DestinationSummaryResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /destinations/{id}/summary/ [get].
func (server *Server) GetDestinationSummary(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetDestinationSummary endpoint called")

	destinationID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	destinationSummary, err := server.Store.GetDestinationSummary(destinationID)
	if err != nil {
		logger.Error(err.Error())

		statusCode, errMsg := utils.ParseDBError(err, "Destination Summary")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	authResponse, err := server.AuthService.GetUserByID(destinationSummary.Owner)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	destinationSummaryResponse := models.DestinationSummaryResponse{
		DestinationName: destinationSummary.DestinationName,
		Owner:           authResponse.Payload.UserInfo,
		CreatedAt:       destinationSummary.CreatedAt,
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", destinationSummaryResponse)

	logger.Info("GetDestinationSummary endpoint returned")
}

// GetDestinationSpecification return the specification of a given destination
// @Summary Get Destination Specification
// @Description Return the specification of a given destination
// @Tags destination
// @Produce  json
// @Param destination query string true "Destination Name"
// @Success 200 {object} models.DestinationSpecification
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /destinations/specification/ [get].
func (server *Server) GetDestinationSpecification(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetDestinationSpecification endpoint called")

	destinationName := ctx.Query("destination")
	if destinationName == "" {
		err := errors.New("no destination specified")
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	logger.Info("Fetching Destination Definitions from AirByte")

	destinationDefinitions, err := server.Airbyte.GetDestinationDefinitions()
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusInternalServerError, utils.ERROR, err.Error(), nil)

		return
	}

	var destinationDefinitionID string

	for _, destinationDefinition := range destinationDefinitions.DestinationDefinitions {
		if destinationDefinition.Name == destinationName {
			destinationDefinitionID = destinationDefinition.DestinationDefinitionID

			break
		}
	}

	if destinationDefinitionID == "" {
		err = errors.New("destination provided does not exist")
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	logger.Info(fmt.Sprintf("Fetching Destination Specification from AirByte for destination: %s", destinationName))

	destinationSpecification, err := server.Airbyte.GetDestinationSpecification(destinationDefinitionID)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusInternalServerError, utils.ERROR, err.Error(), nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", destinationSpecification)

	logger.Info("GetDestinationSpecification endpoint returned")
}

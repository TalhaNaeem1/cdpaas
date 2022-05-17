package dataProduct

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"pipelineService/clients/airbyte"
	"pipelineService/clients/authService"
	"pipelineService/clients/cadenceClient"
	"pipelineService/handlers/v1/pipeline"
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

// CreateDataProduct returns newly created data product
// @Summary Create Data Product
// @Description Creates the data product
// @Tags data-products
// @Produce  json
// @Param product body models.DataProduct true "Product Info"
// @Success 201 {object} models.DataProductResponse
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /data-products/ [post].
func (server *Server) CreateDataProduct(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("CreateDataProduct endpoint called")

	userID, workspaceID, _ := utils.GetUserAndWorkspaceIDFromContext(ctx)

	var product models.DataProduct
	if err := ctx.ShouldBindJSON(&product); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	product.Owner = userID
	product.WorkspaceID = workspaceID

	newProduct, err := server.Store.CreateDataProduct(product)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Data Product")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusCreated, utils.SUCCESS, "", newProduct)
	logger.Info("CreateDataProduct endpoint returned successfully")
}

// GetDataProduct returns a data product
// @Summary Returns a data product
// @Description Returns a data product by ID
// @Tags data-products
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} models.GetDataProductViewResponse
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Security ApiKeyAuth
// @Router /data-products/{id}/ [get].
func (server *Server) GetDataProduct(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetDataProduct endpoint called")

	var (
		dataProduct  models.GetDataProductView
		authResponse models.UserDetails
	)

	//Parse the productID
	dataProductID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	dataProduct.DataProduct, err = server.Store.GetDataProduct(dataProductID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Data Product")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	authResponse, err = server.AuthService.GetUserByID(dataProduct.DataProduct.Owner)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	dataProduct.Owner = authResponse.Payload.UserInfo

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", dataProduct)
	logger.Info("GetDataProduct endpoint returned successfully")
}

// GetAllDataProducts returns a list of data products
// @Summary Returns all the data products
// @Description Returns a list of all the data products
// @Tags data-products
// @Produce  json
// @Success 200 {object} models.DataProductListResponse
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /data-products/ [get].
func (server *Server) GetAllDataProducts(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetAllDataProducts endpoint called")

	_, workspaceID, _ := utils.GetUserAndWorkspaceIDFromContext(ctx)

	products, err := server.Store.GetAllDataProducts(workspaceID)

	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Data Product")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", products)
	logger.Info("GetAllDataProducts endpoint returned successfully")
}

// AddPipeline binds pipeline to data-product
// @Summary Returns product and pipeline ID
// @Description A newly created binding is returned
// @Tags data-products
// @Produce  json
// @Param pipelines body models.InputProductsPipelines true "pipelines"
// @Param id path string true "Product ID"
// @Success 200 {object} models.ProductsPipelinesResponse
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Security ApiKeyAuth
// @Router /data-products/{id}/add-pipeline/ [POST].
func (server *Server) AddPipeline(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("Add Pipeline endpoint called")

	//Parse the productID
	dataProductID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	var pipelines models.InputProductsPipelines
	if err = ctx.ShouldBindJSON(&pipelines); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	var addPipelines []models.ProductsPipelines

	newProductsPipelines := models.ProductsPipelines{}

	for _, inputPipelineID := range pipelines.Pipelines {
		newProductsPipelines.ProductID = dataProductID
		newProductsPipelines.PipelineID = inputPipelineID
		addPipelines = append(addPipelines, newProductsPipelines)
	}

	err = server.Store.AddPipeline(dataProductID, addPipelines)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Data Product")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", "Pipelines added to product")
	logger.Info("AddPipeline endpoint returned successfully")
}

// UpdateDataProduct updates the data-product
// @Summary Returns updated data product
// @Description updates the data product
// @Tags data-products
// @Produce  json
// @Param product body models.DataProduct true "Product Info"
// @Param id path string true "Product ID"
// @Success 200 {object} models.DataProductResponse
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /data-products/{id}/ [put].
func (server *Server) UpdateDataProduct(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("UpdateDataProduct endpoint called")

	//Parse the productID
	dataProductID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	//validate input
	var inputDataProduct models.DataProduct
	if err = ctx.ShouldBindJSON(&inputDataProduct); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	inputDataProduct.ProductID = dataProductID

	updatedDataProduct, err := server.Store.UpdateDataProduct(inputDataProduct)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Data Product")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", updatedDataProduct)
	logger.Info("UpdateDataProduct endpoint returned successfully")
}

// ApplyTransformations applies transformation
// @Summary Apply transformation
// @Description Apply transformation
// @Tags data-products
// @Produce  json
// @Param configureSourceData body models.CreateSourceConnectorRequestAPI true "Source Details"
// @Param id path string true "Product ID"
// @Param destinationId query string true "Destination ID"
// @Param sourceId query string true "Source ID"
// @Success 201 {object} models.CreatePipelineRequest
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /data-products/transformations/{id} [post].
func (server *Server) ApplyTransformations(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("ApplyTransformations endpoint called")

	_, _, airbyteWorkspaceID := utils.GetUserAndWorkspaceIDFromContext(ctx)

	productID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	destinationID, err := uuid.FromString(ctx.Query("destinationId"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	sourceID := ctx.Query("sourceId")

	destination, err := server.Store.GetDestination(destinationID)

	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	source, err := server.Store.GetSource(sourceID)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	var createPipelineRequest models.CreatePipelineRequest
	if err := ctx.ShouldBindJSON(&createPipelineRequest); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	for index := range createPipelineRequest.Operations {
		createPipelineRequest.Operations[index].WorkspaceId = airbyteWorkspaceID
	}

	var airbyteInfo = models.AirbyteSourceAndDestinations{
		AirbyteDestinationID: destination.AirbyteDestinationID,
		AirbyteSourceID:      source.AirbyteSourceID,
	}

	createPipelineAirbyteRequest := pipeline.CreatePipelineAirbyteRequestModel(airbyteInfo, createPipelineRequest)

	newAirByteConnection, err := server.Airbyte.CreateConnection(*createPipelineAirbyteRequest)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	transformationPipeline := models.TransformationPipelines{
		ProductID:           productID.String(),
		SourceID:            sourceID,
		DestinationID:       destinationID.String(),
		AirbyteConnectionID: newAirByteConnection.ConnectionId,
	}

	_, err = server.Store.CreateTransformationPipeline(transformationPipeline)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusCreated, utils.SUCCESS, "", "transformations configured successfully")
	logger.Info("ApplyTransformations endpoint returned")
}

// GetTransformationDetails return DBT transformation details for a given product
// @Summary Return DBT transformation details for a given product
// @Description Return DBT transformation details for a given product
// @Tags data-products
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /data-products/transformations/{id} [get].
func (server *Server) GetTransformationDetails(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetTransformationDetails endpoint called")

	productID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	transformationPipeline, err := server.Store.GetProductConnection(productID)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	requestBody := make(map[string]interface{})
	requestBody["connectionId"] = transformationPipeline.AirbyteConnectionID
	requestBody["withRefreshedCatalog"] = false

	body, err := server.Airbyte.GetConnection(requestBody)

	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	var jsonData interface{}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, "failed to read response", nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusCreated, utils.SUCCESS, "", jsonData)
	logger.Info("GetTransformationDetails endpoint returned")
}

// SyncTransformedAssets syncs the transformed assets with the db
// @Summary syncs the transformed assets with the db
// @Description deletes the already present assets and creates the new entries
// @Tags data-products/internal
// @Produce  json
// @Success 200 {object} models.ProductAssets
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /data-products/internal/transformations/assets/ [post].
func (server *Server) SyncTransformedAssets(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("SyncTransformedAssets endpoint called")

	var productAssestDetails []models.ProductAssetDetails
	if err := ctx.ShouldBindJSON(&productAssestDetails); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	if err := server.Store.SyncTransformedAssets(productAssestDetails); err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Assets")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", nil)
	logger.Info("SyncTransformedAssets endpoint returned successfully")
}

// GetProductDetails return the name of the products
// @Summary Returns updated data product
// @Description returns the name of the data product that can be transformed
// @Tags data-products/internal
// @Produce  json
// @Param product body models.DataProduct true "Product Info"
// @Param id path string true "Product ID"
// @Success 200 {object} models.ProductAssetDetails
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /data-products/internal/ [get].
func (server *Server) GetProductDetails(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetProductNames internal endpoint called")

	productDetails, err := server.Store.GetProductDetails()
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Data Product")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", productDetails)
	logger.Info("GetProductNames internal endpoint returned successfully")
}

// UpdateTransformations updates the transformation details
// @Summary Returns updated data product
// @Description updates the data product
// @Tags data-products
// @Produce  json
// @Param transformation body models.UpdatePipelineAirByteRequest true "Transformation Info"
// @Param id path string true "Product ID"
// @Success 200 {object} models.CreatePipelineRequest
// @Failure 400	{object} models.Response
// @Failure 500	{object} models.Response
// @Router /data-products/transformations/{id}/ [put].
func (server *Server) UpdateTransformations(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("UpdateDataProduct endpoint called")

	_, _, airbyteWorkspaceID := utils.GetUserAndWorkspaceIDFromContext(ctx)

	//Parse the productID
	dataProductID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	transformationPipeline, err := server.Store.GetTransformationPipeline(dataProductID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Asset")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	var updatePipelineRequest models.UpdatePipelineAirByteRequest
	if err := ctx.ShouldBindJSON(&updatePipelineRequest); err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	updatePipelineRequest.ConnectionId = transformationPipeline.AirbyteConnectionID

	for index := range updatePipelineRequest.Operations {
		updatePipelineRequest.Operations[index].WorkspaceId = airbyteWorkspaceID
	}

	_, err = server.Airbyte.UpdateConnection(updatePipelineRequest)
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", "transformation updated successfully")
	logger.Info("UpdatePipelineConnectionOnAirByte endpoint returned successfully")
}

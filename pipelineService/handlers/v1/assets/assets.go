package assets

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/tidwall/gjson"
	"pipelineService/clients/airbyte"
	"pipelineService/clients/authService"
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

// PreviewAsset preview the data from destination
// @Summary Preview the data from destination.
// @Description Preview the data from destination
// @Tags assets
// @Produce  json
// @Param id path string true "Asset ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /assets/{id}/preview/ [get].
func (server *Server) PreviewAsset(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("PreviewAsset endpoint called")

	assetID, err := uuid.FromString(ctx.Param("id"))

	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	assetDetails, err := server.Store.GetAssetDetails(assetID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Assets")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	dbConn, err := db.GetClient(assetDetails.Host, assetDetails.UserName, assetDetails.Password, assetDetails.DbName, assetDetails.Port, "disable")

	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Assets")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	defer db.CloseConnection(dbConn)

	tableName := fmt.Sprintf("%s_%s%s", utils.AIRBYTE_DEFAULT_PREFIX, assetDetails.Prefix, assetDetails.Name)

	assetData, err := server.Store.PreviewData(dbConn, assetDetails.SchemaName, tableName)

	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, "preview data failed", nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", assetData)
	logger.Info("PreviewAsset endpoint returned")
}

// GetPipelineAssets return the assets of a given pipeline
// @Summary Return the assets of a given pipeline
// @Description Return the assets of a given pipeline
// @Tags assets
// @Produce  json
// @Param id path string true "Pipeline ID"
// @Success 200 {object} models.PipelineAssetsResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /assets/pipeline/{id} [get].
func (server *Server) GetPipelineAssets(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetPipelineAssets endpoint called")

	pipelineID, err := uuid.FromString(ctx.Param("id"))

	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	assets, err := server.Store.GetPipelineAssets(pipelineID)

	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Pipeline Assets")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", assets)
	logger.Info("GetPipelineAssets endpoint returned")
}

// GetTransformedAssets return the transformed assets of a given product
// @Summary Return the transformed assets of a given product
// @Description Return the transformed assets of a given product
// @Tags assets
// @Produce  json
// @Param id path string true "Product ID"
// @Success 200 {object} models.TransformedAssetsResponse
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /assets/products/{id}/transformed [get].
func (server *Server) GetTransformedAssets(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("GetTransformedAssets endpoint called")

	productID, err := uuid.FromString(ctx.Param("id"))
	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	transformedAssets, err := server.Store.GetTransformedAssets(productID)

	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Transformed Assets")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusCreated, utils.SUCCESS, "", transformedAssets)
	logger.Info("GetTransformedAssets endpoint returned")
}

// PreviewTransformedAsset preview the data of transformed asset from destination
// @Summary Preview the data of transformed asset from destination.
// @Description Preview the data of transformed asset from destination
// @Tags assets
// @Produce  json
// @Param id path string true "Asset ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /assets/{id}/transformed/preview/ [get].
func (server *Server) PreviewTransformedAsset(ctx *gin.Context) {
	logger := utils.GetLogger()
	logger.Info("PreviewAsset endpoint called")

	assetID, err := uuid.FromString(ctx.Param("id"))

	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, err.Error(), nil)

		return
	}

	assetDetails, err := server.Store.GetTransformedAssetDetails(assetID)
	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Assets")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	destinationConfig := assetDetails.DestinationConfiguration.String()
	host := gjson.Get(destinationConfig, "host").String()
	username := gjson.Get(destinationConfig, "username").String()
	password := gjson.Get(destinationConfig, "password").String()
	database := gjson.Get(destinationConfig, "database").String()
	port := gjson.Get(destinationConfig, "port").String()

	dbConn, err := db.GetClient(host, username, password, database, port, "disable")

	if err != nil {
		logger.Error(err.Error())
		statusCode, errMsg := utils.ParseDBError(err, "Assets")
		utils.BuildResponse(ctx, statusCode, utils.ERROR, errMsg, nil)

		return
	}

	defer db.CloseConnection(dbConn)

	assetData, err := server.Store.PreviewData(dbConn, assetDetails.ProductName, assetDetails.AssetName)

	if err != nil {
		logger.Error(err.Error())
		utils.BuildResponse(ctx, http.StatusBadRequest, utils.ERROR, "preview data failed", nil)

		return
	}

	utils.BuildResponse(ctx, http.StatusOK, utils.SUCCESS, "", assetData)
	logger.Info("PreviewAsset endpoint returned")
}

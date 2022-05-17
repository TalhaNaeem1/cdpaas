//go:generate mockgen -destination=mocks/mock_store.go -package=mock_store . Store

package db

import (
	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"pipelineService/models/v1"
)

type Store interface {
	CreateDataProduct(product models.DataProduct) (models.DataProduct, error)
	GetDataProduct(uuid uuid.UUID) (models.DataProductView, error)
	GetDataProductInfo(productID uuid.UUID) (models.DataProduct, error)
	GetProductConnection(productID uuid.UUID) (models.TransformationPipelines, error)
	GetAllDataProducts(workspaceId int) ([]models.GetAllDataProductsView, error)
	UpdateDataProduct(product models.DataProduct) (models.DataProduct, error)
	AddPipeline(productID uuid.UUID, productPipelines []models.ProductsPipelines) error
	CreatePipeline(pipeline models.Pipeline) (models.Pipeline, error)
	UpdatePipeline(pipeline models.UpdatePipeline) (models.Pipeline, error)
	GetAllPipelines(workspaceId int) ([]models.PipelinesMetaData, error)
	GetPipeline(pipelineID uuid.UUID) (models.PipelineView, error)
	DeletePipeline(pipelineID uuid.UUID) error
	GetPipelineSourceAndConnectionID(pipelineID uuid.UUID) (models.PipelineSourceAndConnectionID, error)
	EnablePipelineAssets(connectionIDs []string) error
	UpdatePipelineStatus(pipelineID uuid.UUID, pipelineStatus string) error

	GetAllConnections() ([]models.Connection, error)
	GetConnection(connectionID string) (models.Connection, error)
	UpdateConnections(connections []models.Connection) error
	UpdateConnectionSchedule(connection models.Connection) error
	GetPipelineConnection(pipelineID string) (models.PipelineConnection, error)
	CreateConnectionAndSourceAgainstAPipeline(source models.Source, connection models.Connection) (models.Source, models.Connection, error)
	GetSupportedSources() ([]models.SupportedSources, error)
	UpdateConnectionInfo(connection models.Connection, destinationId string) error

	GetSource(sourceId string) (models.Source, error)
	GetSourceAndDestinationAirbyteInfo(sourceId string, destinationId string) (models.AirbyteSourceAndDestinations, error)

	CreateDestination(source models.Destination) (models.Destination, error)
	GetSupportedDestinations() ([]models.SupportedDestinations, error)
	GetConfiguredDestination(workspaceId int) ([]models.ConfiguredDestination, error)
	GetDestinationSummary(destinationID uuid.UUID) (models.DestinationSummary, error)
	GetSourceAndConnectionDetails(sourceID string) (models.ConnectionSummary, error)
	CreatePipelineSchema(pipelineSchema models.PipelineSchemas) (models.PipelineSchemas, error)
	CreatePipelineAssets(pipelineAssets []models.PipelineAssets) error
	DeletePipelineSchema(schemaID uuid.UUID) error
	GetPipelineSchema(pipelineID uuid.UUID) (models.PipelineSchemas, error)

	PreviewData(db *gorm.DB, schema string, table string) ([]map[string]interface{}, error)
	GetAssetDetails(assetID uuid.UUID) (models.AssetDetails, error)
	GetPipelineAssets(pipelineID uuid.UUID) ([]models.PipelineAssets, error)

	GetDestination(destinationID uuid.UUID) (models.Destination, error)
	CreateTransformationPipeline(transformationPipeline models.TransformationPipelines) (models.TransformationPipelines, error)
	GetTransformedAssets(productID uuid.UUID) ([]models.ProductAssets, error)
	GetTransformedAssetDetails(assetID uuid.UUID) (models.TransformedAssetDetails, error)
	GetProductDetails() ([]models.ProductDetail, error)
	SyncTransformedAssets(productAssetDetails []models.ProductAssetDetails) error
	GetTransformationPipeline(productID uuid.UUID) (models.TransformationPipelines, error)
}

type PGStore struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) Store {
	return &PGStore{
		db: db,
	}
}

var _ Store = (*PGStore)(nil)

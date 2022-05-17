package db

import (
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"pipelineService/models/v1"
)

func (p *PGStore) CreateDataProduct(product models.DataProduct) (models.DataProduct, error) {
	createdProduct := models.DataProduct{}

	result := p.db.Create(&product).Scan(&createdProduct)

	return createdProduct, result.Error
}

func (p *PGStore) GetDataProductInfo(productID uuid.UUID) (models.DataProduct, error) {
	var product models.DataProduct

	result := p.db.Where("product_id = ?", productID).
		First(&product)

	return product, result.Error
}

func (p *PGStore) GetProductConnection(productID uuid.UUID) (models.TransformationPipelines, error) {
	var product models.TransformationPipelines

	result := p.db.Where("product_id = ?", productID).First(&product)

	return product, result.Error
}

func (p *PGStore) GetDataProduct(dataProductID uuid.UUID) (models.DataProductView, error) {
	dataProduct := models.DataProductView{}

	result := p.db.Table("data_products").Where("product_id = ?", dataProductID).First(&dataProduct)

	var pipelines []map[string]interface{}

	_ = p.db.Table("pipelines").
		Select("pipelines.*").
		Joins("join products_pipelines on pipelines.pipeline_id = products_pipelines.pipeline_id").
		Where("products_pipelines.product_id = ?", dataProductID).
		Find(&pipelines)

	dataProduct.Pipelines = pipelines

	return dataProduct, result.Error
}

func (p *PGStore) GetAllDataProducts(workspaceId int) ([]models.GetAllDataProductsView, error) {
	dataProducts := make([]models.GetAllDataProductsView, 0)

	result := p.db.Table("data_products").
		Select("data_products.*, count(products_pipelines.product_id) AS pipeline_count").
		Joins("left join products_pipelines on data_products.product_id = products_pipelines.product_id").
		Where("workspace_id = ?", workspaceId).
		Group("data_products.product_id").
		Order("data_products.created_at DESC").
		Find(&dataProducts)

	return dataProducts, result.Error
}

func (p *PGStore) UpdateDataProduct(product models.DataProduct) (models.DataProduct, error) {
	updatedDataProduct := models.DataProduct{}

	result := p.db.Updates(product).Scan(&updatedDataProduct)

	if result.RowsAffected == 0 {
		return updatedDataProduct, errors.New("Product doesn't exists")
	}

	return updatedDataProduct, result.Error
}

func (p *PGStore) AddPipeline(productID uuid.UUID, addPipelines []models.ProductsPipelines) error {
	return db.Transaction(func(tx *gorm.DB) error {
		result := p.db.Where("product_id = ? ", productID).Delete(&models.ProductsPipelines{})
		if result.Error != nil {
			return result.Error
		}

		result = p.db.Create(&addPipelines)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
}

func (p *PGStore) CreatePipeline(pipeline models.Pipeline) (models.Pipeline, error) {
	createdPipeline := models.Pipeline{}

	result := p.db.Create(&pipeline).Scan(&createdPipeline)

	return createdPipeline, result.Error
}

func (p *PGStore) UpdatePipeline(pipeline models.UpdatePipeline) (models.Pipeline, error) {
	updatedPipeline := models.Pipeline{}

	result := p.db.Table("pipelines").Updates(pipeline).Scan(&updatedPipeline)

	if result.RowsAffected == 0 {
		return updatedPipeline, errors.New("Pipeline doesn't exists")
	}

	return updatedPipeline, result.Error
}

func (p *PGStore) GetSupportedSources() ([]models.SupportedSources, error) {
	sources := make([]models.SupportedSources, 0)
	p.db.Find(&sources)
	result := p.db.Find(&sources)

	return sources, result.Error
}

func (p *PGStore) GetAllConnections() ([]models.Connection, error) {
	connections := make([]models.Connection, 0)

	result := p.db.Find(&connections)

	return connections, result.Error
}

func (p *PGStore) GetConnection(connectionID string) (models.Connection, error) {
	connection := models.Connection{}

	result := p.db.Where("connection_id = ?", connectionID).First(&connection)

	return connection, result.Error
}

func (p *PGStore) UpdateConnectionSchedule(connection models.Connection) error {
	updateConnection := make(map[string]interface{})

	if connection.AirbyteTimeUnit != "" && connection.AirbyteFrequencyUnits != 0 {
		updateConnection["airbyte_frequency_units"] = connection.AirbyteFrequencyUnits
		updateConnection["airbyte_time_unit"] = connection.AirbyteTimeUnit
		updateConnection["is_first_run"] = false
	} else {
		updateConnection["airbyte_frequency_units"] = nil
		updateConnection["airbyte_time_unit"] = "minutes"
		updateConnection["is_first_run"] = false
	}

	result := p.db.Model(&connection).Updates(&updateConnection)

	return result.Error
}

func (p *PGStore) UpdateConnections(connections []models.Connection) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, connection := range connections {
			result := p.db.Model(&connection).
				Updates(connection)

			if result.Error != nil {
				return result.Error
			}
		}

		return nil
	})
}

func (p *PGStore) GetAllPipelines(workspaceID int) ([]models.PipelinesMetaData, error) {
	results := make([]models.PipelinesMetaData, 0)

	p.db.Table("pipelines").
		Select("pipelines.pipeline_id AS pipeline_id, "+
			"pipelines.name AS pipeline_name, "+
			"pipelines.pipeline_status AS pipeline_status, "+
			"pipelines.pipeline_governance AS pipeline_governance, "+
			"sources.name AS source_name, "+
			"sources.source_id AS source_id, "+
			"destinations.name AS destination_name, "+
			"destinations.destination_id AS destination_id, "+
			"connections.airbyte_status AS airbyte_status, "+
			"connections.airbyte_connection_id AS airbyte_connection_id, "+
			"connections.connection_id AS connection_id, "+
			"connections.airbyte_last_run AS airbyte_last_run, "+
			"pipelines.owner AS owner_id").
		Joins("LEFT join connections on pipelines.pipeline_id = connections.pipeline_id").
		Joins("LEFT join sources on connections.connection_id = sources.connection_id").
		Joins("LEFT join connections_destinations on connections.connection_id = connections_destinations.connection_id").
		Joins("LEFT join destinations on connections_destinations.destination_id = destinations.destination_id").
		Where("pipelines.workspace_id = ?", workspaceID).
		Order("pipelines.created_at DESC").
		Find(&results)

	return results, nil
}

func (p *PGStore) GetPipeline(pipelineID uuid.UUID) (models.PipelineView, error) {
	pipeline := models.PipelineView{}

	result := p.db.Table("pipelines").
		Select("pipelines.pipeline_id AS pipeline_id, "+
			"pipelines.name AS pipeline_name, "+
			"pipelines.pipeline_governance AS pipeline_governance, "+
			"pipelines.created_at AS created_at, "+
			"pipelines.owner AS pipeline_owner, "+
			"sources.name AS source_name, "+
			"sources.source_id AS source_id, "+
			"destinations.name AS destination_name, "+
			"destinations.destination_id AS destination_id, "+
			"connections.airbyte_status AS airbyte_status, "+
			"connections.airbyte_connection_id AS airbyte_connection_id, "+
			"connections.connection_id AS connection_id, "+
			"connections.airbyte_last_run AS airbyte_last_run ").
		Where("pipelines.pipeline_id = ?", pipelineID).
		Joins("join connections on pipelines.pipeline_id = connections.pipeline_id").
		Joins("join sources on connections.connection_id = sources.connection_id").
		Joins("join connections_destinations on connections.connection_id = connections_destinations.connection_id").
		Joins("join destinations on connections_destinations.destination_id = destinations.destination_id").
		Find(&pipeline)

	var dataProducts []map[string]interface{}

	_ = p.db.Table("data_products").
		Select("data_products.*").
		Joins("join products_pipelines on data_products.product_id = products_pipelines.product_id").
		Where("products_pipelines.pipeline_id = ?", pipelineID).
		Find(&dataProducts)

	pipeline.Product = dataProducts

	return pipeline, result.Error
}

func (p *PGStore) DeletePipeline(pipelineID uuid.UUID) error {
	result := p.db.Where("pipeline_id = ?", pipelineID).Delete(&models.Pipeline{})

	return result.Error
}

func (p *PGStore) GetPipelineSourceAndConnectionID(pipelineID uuid.UUID) (models.PipelineSourceAndConnectionID, error) {
	pipeline := models.PipelineSourceAndConnectionID{}

	result := p.db.Table("pipelines").
		Select("pipelines.pipeline_id AS pipeline_id, "+
			"connections.connection_id, "+
			"sources.airbyte_source_id AS airbyte_source_id, "+
			"connections.airbyte_connection_id AS airbyte_connection_id").
		Joins("LEFT join connections on pipelines.pipeline_id = connections.pipeline_id").
		Joins("LEFT join sources on connections.connection_id = sources.connection_id").
		Where("pipelines.pipeline_id = ?", pipelineID).
		Find(&pipeline)

	return pipeline, result.Error
}

func (p *PGStore) GetSourceAndConnectionDetails(sourceID string) (models.ConnectionSummary, error) {
	connectionSummary := models.ConnectionSummary{}

	result := p.db.Table("sources").
		Select("sources.name,"+
			" connections.owner,"+
			" connections.created_at,"+
			" connections.airbyte_connection_id").
		Joins("join connections on sources.connection_id = connections.connection_id").
		Where("sources.source_id = ?", sourceID).
		Scan(&connectionSummary)

	return connectionSummary, result.Error
}

func (p *PGStore) UpdatePipelineStatus(pipelineID uuid.UUID, pipelineStatus string) error {
	result := p.db.Table("pipelines").Where("pipeline_id = ?", pipelineID).Update("pipeline_status", pipelineStatus)

	if result.RowsAffected == 0 {
		return errors.New("Pipeline Doesnt exist")
	}

	return result.Error
}

func (p *PGStore) GetSource(sourceID string) (models.Source, error) {
	source := models.Source{}

	result := p.db.Where("source_id = ?", sourceID).First(&source)

	return source, result.Error
}

func (p *PGStore) GetPipelineConnection(pipelineID string) (models.PipelineConnection, error) {
	var pipelineConnections models.PipelineConnection

	result := p.db.Table("connections").
		Select("connections.connection_id AS connection_id, "+
			"connections.pipeline_id AS pipeline_id, "+
			"sources.source_id AS source_id, "+
			"sources.name AS source_name ").
		Where("connections.pipeline_id = ?", pipelineID).
		Joins("join sources on sources.connection_id = connections.connection_id").
		Scan(&pipelineConnections)

	return pipelineConnections, result.Error
}

func (p *PGStore) CreateConnectionAndSourceAgainstAPipeline(source models.Source, connection models.Connection) (models.Source, models.Connection, error) {
	var insertedSource models.Source

	var insertedConnection models.Connection

	tx := p.db.Begin()

	result := tx.Select("pipeline_id", "connection_id").Create(&connection).Scan(&insertedConnection)
	if result.Error != nil {
		err := tx.Rollback()
		if err.Error != nil {
			return insertedSource, insertedConnection, err.Error
		}

		return insertedSource, insertedConnection, result.Error
	}

	source.ConnectionID = insertedConnection.ConnectionID

	result = tx.Create(&source).Scan(&insertedSource)
	if result.Error != nil {
		err := tx.Rollback()
		if err.Error != nil {
			return insertedSource, insertedConnection, err.Error
		}

		return insertedSource, insertedConnection, result.Error
	}

	if err := tx.Commit(); err.Error != nil {
		return insertedSource, insertedConnection, err.Error
	}

	return insertedSource, insertedConnection, nil
}

func (p *PGStore) GetSourceAndDestinationAirbyteInfo(sourceId, destinationId string) (models.AirbyteSourceAndDestinations, error) {
	var (
		sourceAndDestination models.AirbyteSourceAndDestinations
		destination          models.Destination
		source               models.SourceAndPipelineName
	)

	result := p.db.Select("airbyte_destination_id", "destination_id", "destination_type", "configuration_details").
		Where("destination_id = ?", destinationId).
		First(&destination)
	if result.Error != nil {
		return sourceAndDestination, result.Error
	}

	if result.RowsAffected != 1 {
		err := errors.New("no such destination exists")

		return sourceAndDestination, err
	}

	result = p.db.Table("sources").
		Select("sources.airbyte_source_id, "+
			"sources.connection_id, "+
			"sources.source_id, "+
			"pipelines.name AS pipeline_name, "+
			"pipelines.pipeline_id AS pipeline_id ").
		Joins("join connections on sources.connection_id = connections.connection_id").
		Joins("join pipelines on connections.pipeline_id = pipelines.pipeline_id").
		Where("source_id = ?", sourceId).
		First(&source)
	if result.Error != nil {
		return sourceAndDestination, result.Error
	}

	if result.RowsAffected != 1 {
		err := errors.New("no such source exists")

		return sourceAndDestination, err
	}

	sourceAndDestination = models.AirbyteSourceAndDestinations{
		DestinationID:                   destination.DestinationID,
		DestinationType:                 destination.DestinationType,
		DestinationConfigurationDetails: destination.ConfigurationDetails,
		AirbyteDestinationID:            destination.AirbyteDestinationID,
		SourceID:                        source.SourceID,
		AirbyteSourceID:                 source.AirbyteSourceID,
		ConnectionID:                    source.ConnectionID,
		PipelineName:                    source.PipelineName,
		PipelineID:                      source.PipelineID,
	}

	return sourceAndDestination, nil
}

func (p *PGStore) UpdateConnectionInfo(connection models.Connection, destinationId string) error {
	tx := p.db.Begin()

	result := tx.Clauses(clause.Returning{Columns: []clause.Column{{Name: "connection_id"}}}).Updates(&connection)
	if result.Error != nil {
		err := tx.Rollback()
		if err.Error != nil {
			return err.Error
		}

		return result.Error
	}

	connectionDestination := models.ConnectionsDestinations{
		ConnectionID:  connection.ConnectionID,
		DestinationID: destinationId,
	}

	result = tx.Create(&connectionDestination)
	if result.Error != nil {
		err := tx.Rollback()
		if err.Error != nil {
			return err.Error
		}

		return result.Error
	}

	if err := tx.Commit(); err.Error != nil {
		return err.Error
	}

	return nil
}

func (p *PGStore) CreateDestination(destination models.Destination) (models.Destination, error) {
	var insertedDestination models.Destination

	result := p.db.Create(&destination).Scan(&insertedDestination)
	if result.Error != nil {
		return insertedDestination, result.Error
	}

	return insertedDestination, nil
}

func (p *PGStore) GetSupportedDestinations() ([]models.SupportedDestinations, error) {
	destinations := make([]models.SupportedDestinations, 0)

	result := p.db.Find(&destinations)
	if result.Error != nil {
		return destinations, result.Error
	}

	return destinations, result.Error
}

func (p *PGStore) GetConfiguredDestination(workspaceId int) ([]models.ConfiguredDestination, error) {
	configuredDestinations := make([]models.ConfiguredDestination, 0)

	result := p.db.Table("destinations").
		Select("destination_id,name,destination_type,airbyte_destination_id,configuration_details::json->>'host' as host").
		Where("destinations.workspace_id = ?", workspaceId).
		Order("destinations.created_at DESC").
		Scan(&configuredDestinations)

	return configuredDestinations, result.Error
}

func (p *PGStore) GetDestinationSummary(destinationID uuid.UUID) (models.DestinationSummary, error) {
	var destinationSummary models.DestinationSummary

	result := p.db.Table("destinations").
		Select("name,"+
			"owner,"+
			"created_at").
		Where("destinations.destination_id = ?", destinationID).
		First(&destinationSummary)

	return destinationSummary, result.Error
}

func (p *PGStore) CreatePipelineSchema(pipelineSchema models.PipelineSchemas) (models.PipelineSchemas, error) {
	createdPipelineSchema := models.PipelineSchemas{}

	result := p.db.Where("name =?", pipelineSchema.Name).FirstOrCreate(&pipelineSchema).Scan(&createdPipelineSchema)

	return createdPipelineSchema, result.Error
}

func (p *PGStore) CreatePipelineAssets(pipelineAssets []models.PipelineAssets) error {
	return db.Transaction(func(tx *gorm.DB) error {
		pipelineID, _ := uuid.FromString(pipelineAssets[0].PipelineID)

		result := p.db.Where("pipeline_id = ? ", pipelineID).Delete(&models.PipelineAssets{})
		if result.Error != nil {
			return result.Error
		}

		for _, pipelineAsset := range pipelineAssets {
			result := p.db.Model(&pipelineAsset).Create(&pipelineAsset)

			if result.Error != nil {
				return result.Error
			}
		}

		return nil
	})
}

func (p *PGStore) DeletePipelineSchema(schemaID uuid.UUID) error {
	result := p.db.Where("schema_id = ?", schemaID).Delete(&models.PipelineSchemas{})

	return result.Error
}

func (p *PGStore) GetPipelineSchema(pipelineID uuid.UUID) (models.PipelineSchemas, error) {
	var schema models.PipelineSchemas
	result := p.db.Where("pipeline_id = ?", pipelineID).First(&schema)

	return schema, result.Error
}

func (p *PGStore) EnablePipelineAssets(connectionIDs []string) error {
	var pipelineIDs []string

	p.db.Table("connections").
		Select("pipeline_id").
		Where("connection_id IN ? ", connectionIDs).
		Find(&pipelineIDs)

	result := p.db.Table("pipeline_assets").
		Where("pipeline_assets.pipeline_id IN ?", pipelineIDs).
		Update("is_enabled", true)

	if result.RowsAffected == 0 {
		return errors.New("No assets exists")
	}

	return result.Error
}

func (p *PGStore) GetDestination(destinationID uuid.UUID) (models.Destination, error) {
	var destination models.Destination

	result := p.db.Table("destinations").Where("destinations.destination_id = ?", destinationID).
		First(&destination)

	return destination, result.Error
}

func (p *PGStore) CreateTransformationPipeline(transformationPipeline models.TransformationPipelines) (models.TransformationPipelines, error) {
	createdTransformationPipeline := models.TransformationPipelines{}

	result := p.db.Create(&transformationPipeline).Scan(&createdTransformationPipeline)

	return createdTransformationPipeline, result.Error
}

func (p *PGStore) GetProductDetails() ([]models.ProductDetail, error) {
	var products []models.ProductDetail

	result := p.db.Table("transformation_pipelines").
		Select("data_products.name, " +
			"data_products.product_id, " +
			"data_products.owner as owner, " +
			"data_products.workspace_id as workspace_id, " +
			"data_products.product_id, " +
			"destinations.configuration_details").
		Joins("join data_products on transformation_pipelines.product_id = data_products.product_id").
		Joins("join destinations on transformation_pipelines.destination_id = destinations.destination_id").
		Find(&products)

	return products, result.Error
}

func (p *PGStore) GetTransformationPipeline(productID uuid.UUID) (models.TransformationPipelines, error) {
	var transformationPipeline models.TransformationPipelines
	result := p.db.Where("product_id = ?", productID).First(&transformationPipeline)

	return transformationPipeline, result.Error
}

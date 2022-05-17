package db

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
	"pipelineService/models/v1"
	"pipelineService/utils"
)

func (p *PGStore) PreviewData(db *gorm.DB, schema string, table string) ([]map[string]interface{}, error) {
	schema = strings.ToLower(schema)

	tableName := fmt.Sprintf("%s.%s", schema, table)

	records, err := db.Table(tableName).Select("*").Limit(utils.PREVIEW_DATA_LIMIT).Rows()

	data := make([]map[string]interface{}, 0)

	if err != nil {
		return data, err
	}

	for records.Next() {
		var record map[string]interface{}
		err := db.ScanRows(records, &record)

		if err != nil {
			return data, err
		}

		data = append(data, record)
	}

	return data, err
}

func (p *PGStore) GetAssetDetails(assetID uuid.UUID) (models.AssetDetails, error) {
	var asset models.AssetDetails

	result := p.db.Table("pipeline_assets").
		Select("pipeline_assets.name AS name, "+
			"pipeline_schemas.name AS pipeline_schema_name, "+
			"pipeline_schemas.prefix AS prefix, "+
			"configuration_details::json->>'host' as host,"+
			"configuration_details::json->>'port' as port,"+
			"configuration_details::json->>'username' as username,"+
			"configuration_details::json->>'password' as password,"+
			"configuration_details::json->>'database' as database").
		Where("pipeline_assets.asset_id = ?", assetID).
		Where("pipeline_assets.is_enabled = ?", true).
		Joins("join pipelines on pipelines.pipeline_id = pipeline_assets.pipeline_id").
		Joins("join pipeline_schemas on pipelines.pipeline_id = pipeline_schemas.pipeline_id").
		Joins("join connections on pipelines.pipeline_id = connections.pipeline_id").
		Joins("join connections_destinations on connections.connection_id = connections_destinations.connection_id").
		Joins("join destinations on connections_destinations.destination_id = destinations.destination_id").
		Find(&asset)

	return asset, result.Error
}

func (p *PGStore) GetPipelineAssets(pipelineID uuid.UUID) ([]models.PipelineAssets, error) {
	var assets []models.PipelineAssets

	result := p.db.Table("pipeline_assets").
		Select("*").
		Where("pipeline_assets.pipeline_id = ?", pipelineID).
		Where("pipeline_assets.is_enabled = ?", true).
		Find(&assets)

	return assets, result.Error
}

func (p *PGStore) GetTransformedAssets(productID uuid.UUID) ([]models.ProductAssets, error) {
	var transformedAssets []models.ProductAssets

	result := p.db.Table("product_assets").Where("product_id = ?", productID).Find(&transformedAssets)

	return transformedAssets, result.Error
}

func (p *PGStore) GetTransformedAssetDetails(assetID uuid.UUID) (models.TransformedAssetDetails, error) {
	var transformedAsset models.TransformedAssetDetails

	result := p.db.Table("product_assets").
		Select("product_assets.name AS asset_name, "+
			"product_assets.asset_id AS asset_id, "+
			"data_products.name AS product_name, "+
			"data_products.product_id AS product_id, "+
			"destinations.configuration_details as destination_configuration").
		Joins("join data_products on data_products.product_id = product_assets.product_id").
		Joins("join transformation_pipelines on transformation_pipelines.product_id = product_assets.product_id").
		Joins("join destinations on destinations.destination_id = transformation_pipelines.destination_id").
		Where("product_assets.asset_id = ?", assetID).
		Where("product_assets.is_enabled = ?", true).
		Find(&transformedAsset)

	return transformedAsset, result.Error
}

func (p *PGStore) SyncTransformedAssets(productAssetDetails []models.ProductAssetDetails) error {
	var createProductAssets []models.ProductAssets

	return db.Transaction(func(tx *gorm.DB) error {
		for _, productAsset := range productAssetDetails {
			result := p.db.Where("product_id = ? ", productAsset.ProductID).Delete(&models.ProductAssets{})
			if result.Error != nil {
				return result.Error
			}

			if len(productAsset.Table) > 0 {
				for _, tableName := range productAsset.Table {
					newProductAsset := models.ProductAssets{
						ProductID:   productAsset.ProductID,
						Name:        tableName,
						IsEnabled:   true,
						Columns:     nil,
						Owner:       productAsset.Owner,
						WorkspaceID: productAsset.WorkspaceID,
					}

					createProductAssets = append(createProductAssets, newProductAsset)
				}

				result = p.db.Table("product_assets").Create(&createProductAssets)
				if result.Error != nil {
					return result.Error
				}
			}
		}

		return nil
	})
}

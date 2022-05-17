package models

import "gorm.io/datatypes"

type PipelineAssets struct {
	AssetID     string         `gorm:"type:uuid;primaryKey;default:(-)" json:"assetID" example:"b251379e-01a1-11ec-82d6-a312edcd9c7b"`
	SchemaID    string         `gorm:"column:pipeline_schemas_id" json:"pipelineSchemaID" example:"b251379e-01a1-11ec-82d6-a312edcd9c7b"`
	PipelineID  string         `gorm:"column:pipeline_id" json:"pipelineID" example:"b251379e-01a1-11ec-82d6-a312edcd9c7b"`
	Name        string         `gorm:"column:name" json:"name" example:"schema"`
	IsEnabled   bool           `gorm:"column:is_enabled" json:"isEnabled" example:"false"`
	Columns     datatypes.JSON `json:"columns" gorm:"column:columns;type:json" example:"{}"`
	Owner       int            `json:"owner" gorm:"type:int" example:"1"`
	WorkspaceID int            `json:"workspaceId" gorm:"type:int" example:"1"`
}

type AssetDetails struct {
	Name       string `json:"name" gorm:"column:name;type:string;size:50"`
	SchemaName string `json:"schemaName" gorm:"column:pipeline_schema_name;type:string;size:50"`
	Prefix     string `json:"prefix" gorm:"column:prefix;type:string;size:50"`
	Host       string `json:"host" gorm:"column:host;type:string;size:50"`
	Port       string `json:"port" gorm:"column:port;type:string;size:50"`
	DbName     string `json:"database" gorm:"column:database;type:string;size:50"`
	UserName   string `json:"username" gorm:"column:username;type:string;size:50"`
	Password   string `json:"password" gorm:"column:password;type:string;size:50"`
}

type PipelineAssetsResponse struct {
	Status string           `json:"status" example:"success"`
	Errors string           `json:"errors" example:""`
	Data   []PipelineAssets `json:"data"`
}

type ProductAssets struct {
	AssetID     string         `gorm:"type:uuid;primaryKey;default:(-)" json:"assetID" example:"b251379e-01a1-11ec-82d6-a312edcd9c7b"`
	ProductID   string         `gorm:"column:product_id" json:"productID" example:"b251379e-01a1-11ec-82d6-a312edcd9c7b"`
	Name        string         `gorm:"column:name" json:"name" example:"schema"`
	IsEnabled   bool           `gorm:"column:is_enabled" json:"isEnabled" example:"false"`
	Columns     datatypes.JSON `json:"columns" gorm:"column:columns;type:json" example:"{}"`
	Owner       int            `json:"owner" gorm:"type:int" example:"1"`
	WorkspaceID int            `json:"workspaceId" gorm:"type:int" example:"1"`
}

type TransformedAssetsResponse struct {
	Status string          `json:"status" example:"success"`
	Errors string          `json:"errors" example:""`
	Data   []ProductAssets `json:"data"`
}

type TransformedAssetDetails struct {
	AssetID                  string         `gorm:"type:uuid;primaryKey;default:(-)" json:"assetID" example:"b251379e-01a1-11ec-82d6-a312edcd9c7b"`
	AssetName                string         `gorm:"column:asset_name" json:"assetName" example:"asset"`
	ProductID                string         `gorm:"column:product_id" json:"productID" example:"b251379e-01a1-11ec-82d6-a312edcd9c7b"`
	ProductName              string         `gorm:"column:product_name" json:"productName" example:"product"`
	DestinationConfiguration datatypes.JSON `gorm:"column:destination_configuration" json:"destinationConfiguration" example:"{}"`
}

type ProductAssetDetails struct {
	Schema      string   `json:"schemaName"`
	ProductID   string   `json:"productID"`
	Owner       int      `json:"owner"`
	WorkspaceID int      `json:"workspaceId"`
	Table       []string `json:"tableNames" gorm:"column:table_name" example:""`
}

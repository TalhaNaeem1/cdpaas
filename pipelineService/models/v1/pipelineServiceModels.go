package models

import (
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type Response struct {
	Status string      `json:"status"`
	Errors string      `json:"errors"`
	Data   interface{} `json:"data"`
}

type DataProduct struct {
	ProductID             uuid.UUID      `json:"productID" gorm:"type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	Name                  *string        `json:"name" gorm:"type:string;size:50" example:"data-product-1"`
	DataProductGovernance pq.StringArray `json:"dataProductGovernance" gorm:"type:string[]" example:"[sales, marketing]"`
	DataDomain            *string        `json:"dataDomain" gorm:"type:string;size:100" example:"data-product-1"`
	Description           *string        `json:"description" gorm:"type:text" example:"This data product will help you in creating better visualizations"`
	DataProductStatus     *string        `json:"dataProductStatus" gorm:"type:string;size:50" example:"draft"`
	LastUpdated           int64          `json:"lastUpdated" gorm:"autoUpdateTime:milli"`
	Owner                 int            `json:"owner" gorm:"type:int" example:"1"`
	WorkspaceID           int            `json:"workspaceId" gorm:"type:int" example:"1"`
}

type GetAllDataProductsView struct {
	ProductID             uuid.UUID      `json:"productID" gorm:"type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	Name                  string         `json:"name" gorm:"type:string;size:50" example:"data-product-1"`
	DataProductGovernance pq.StringArray `json:"dataProductGovernance" gorm:"type:string[]" example:"[sales, marketing]"`
	DataDomain            string         `json:"dataDomain" gorm:"type:string;size:100" example:"data-product-1"`
	Description           string         `json:"description" gorm:"type:text" example:"This data product will help you in creating better visualizations"`
	DataProductStatus     string         `json:"dataProductStatus" gorm:"type:string;size:50" example:"draft"`
	LastUpdated           int64          `json:"lastUpdated" gorm:"autoUpdateTime:milli"`
	Owner                 int            `json:"owner" gorm:"type:int" example:"1"`
	WorkspaceID           int            `json:"workspaceId" gorm:"type:int" example:"1"`
	PipelineCount         int            `json:"pipelineCount" gorm:"column:pipeline_count"  example:"0"`
}

type DataProductResponse struct {
	Status string      `json:"status" example:"success"`
	Errors string      `json:"errors" example:""`
	Data   DataProduct `json:"data"`
}

type DataProductListResponse struct {
	Status string                   `json:"status" example:"success"`
	Errors string                   `json:"errors" example:""`
	Data   []GetAllDataProductsView `json:"data"`
}

type DataProductView struct {
	ProductID             uuid.UUID                `json:"productID" gorm:"type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	Name                  string                   `json:"name" gorm:"type:string;size:50" example:"data-product-1"`
	DataProductGovernance pq.StringArray           `json:"dataProductGovernance" gorm:"type:string[]" example:"[sales, marketing]"`
	DataDomain            string                   `json:"dataDomain" gorm:"type:string;size:100" example:"data-product-1"`
	Description           string                   `json:"description" gorm:"type:text" example:"This data product will help you in creating better visualizations"`
	DataProductStatus     string                   `json:"dataProductStatus" gorm:"type:string;size:50" example:"completed"`
	LastUpdated           int64                    `json:"lastUpdated" gorm:"autoUpdateTime:milli"`
	Owner                 int                      `json:"-" gorm:"type:int" example:"1"`
	WorkspaceID           int                      `json:"workspaceId" gorm:"type:int" example:"1"`
	Pipelines             []map[string]interface{} `json:"pipelines"`
}

type GetDataProductView struct {
	DataProduct DataProductView `json:"dataProduct"`
	Owner       interface{}     `json:"owner"`
}

type GetDataProductViewResponse struct {
	Status string             `json:"status" example:"success"`
	Errors string             `json:"errors" example:""`
	Data   GetDataProductView `json:"data"`
}

type InputProductsPipelines struct {
	Pipelines []string `json:"pipelines" example:"[\"a152379e-01a1-11ec-82d6-a312edcd9c7b\", \"a152379e-01a1-11ec-82d6-a312edcd9c7b\"]"`
}

type ProductsPipelines struct {
	ProductsPipelinesID uuid.UUID `json:"productPipelineID" gorm:"type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	ProductID           uuid.UUID `json:"productID" gorm:"type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	PipelineID          string    `json:"pipelineID" gorm:"column:pipeline_id; type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
}

type ProductsPipelinesResponse struct {
	Status string `json:"status" example:"success"`
	Errors string `json:"errors" example:""`
	Data   string `json:"data"`
}

type ProductDetail struct {
	ProductID            uuid.UUID      `json:"productID" gorm:"type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	Name                 string         `json:"name" gorm:"column:name" example:"data-product-1"`
	ConfigurationDetails datatypes.JSON `json:"configurationDetails" gorm:"column:configuration_details; type:string; default:(-)" example:"example_destination"`
	Owner                int            `json:"owner" gorm:"column:owner;type:int" example:"2"`
	WorkspaceID          int            `json:"workspaceId" gorm:"column:workspace_id;type:int" example:"1"`
}

type ConfigurationDetails struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	DbName   string `json:"database"`
	UserName string `json:"username"`
	Password string `json:"password"`
}

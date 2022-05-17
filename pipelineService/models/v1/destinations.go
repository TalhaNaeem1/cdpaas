package models

import (
	"github.com/gofrs/uuid"
	"gorm.io/datatypes"
)

type Destination struct {
	DestinationID           string         `json:"destinationId" binding:"required" gorm:"column:destination_id; type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	DestinationName         string         `json:"destinationName" binding:"required" gorm:"column:name; type:string; default:(-)" example:"example_destination"`
	AirbyteDestinationID    string         `json:"airbyteDestinationId" gorm:"column:airbyte_destination_id; type:uuid; default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	AirbyteDestDefinitionID string         `json:"airbyteDestinationDefinitionId" gorm:"column:airbyte_destination_definition_id; type:uuid ; default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	DestinationType         string         `json:"destinationType" binding:"required" gorm:"column:destination_type; type:string; default:(-)" example:"example_destination"`
	ConfigurationDetails    datatypes.JSON `json:"configurationDetails" binding:"required" gorm:"column:configuration_details; type:string; default:(-)" example:"example_destination"`
	Owner                   int            `json:"owner" gorm:"type:int" example:"1"`
	WorkspaceID             int            `json:"workspaceId" gorm:"type:int" example:"1"`
	CreatedAt               int64          `json:"createdAt" gorm:"default:(extract(epoch from now()) * 1000)"`
}

type ConfiguredDestination struct {
	DestinationID           string `json:"destinationId" binding:"required" gorm:"column:destination_id; type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	DestinationName         string `json:"destinationName" binding:"required" gorm:"column:name; type:string; default:(-)" example:"example_destination"`
	AirbyteDestinationID    string `json:"airbyteDestinationId" gorm:"column:airbyte_destination_id; type:uuid; default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	AirbyteDestDefinitionID string `json:"airbyteDestinationDefinitionId" gorm:"column:airbyte_destination_definition_id; type:uuid ; default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	DestinationType         string `json:"destinationType" binding:"required" gorm:"column:destination_type; type:string; default:(-)" example:"example_destination"`
	Host                    string `json:"host" gorm:"column:host; type:string; default:(-)"`
}

type DestinationSummary struct {
	DestinationName string `json:"destinationName" binding:"required" gorm:"column:name; type:string; default:(-)" example:"example_destination"`
	Owner           int    `json:"owner" gorm:"type:int" example:"1"`
	CreatedAt       int64  `json:"createdAt" gorm:"column:created_at;"`
}

type DestinationSummaryResponse struct {
	DestinationName string      `json:"destinationName" binding:"required" example:"example_destination"`
	Owner           interface{} `json:"owner" binding:"required"`
	CreatedAt       int64       `json:"createdAt" binding:"required"`
}

type SupportedDestinations struct {
	ID   uuid.UUID `json:"id" binding:"required" gorm:"type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	Name string    `json:"name" binding:"required" gorm:"type:string;size:50" example:"destination name"`
	Type string    `json:"type" binding:"required" gorm:"type:string;size:50" example:"destination type"`
}

type SupportedDestinationResponse struct {
	Status string                  `json:"status" example:"success"`
	Errors string                  `json:"errors" example:""`
	Data   []SupportedDestinations `json:"data"`
}

type ConfiguredDestinationResponse struct {
	Status string                  `json:"status" example:"success"`
	Errors string                  `json:"errors" example:""`
	Data   []ConfiguredDestination `json:"data"`
}

type CreateDestinationConnectorRequestAPI struct {
	CreateDestinationConnectorRequest
	DestinationType string `json:"destinationName" binding:"required" example:"Postgres"`
}

type CreateDestinationConnectorResponseData struct {
	DestinationID   string `json:"destinationId" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	DestinationName string `json:"name" example:"example_destination"`
	DestinationType string `json:"destinationName" example:"Postgres"`
}

type CreateDestinationConnectorResponseAPI struct {
	Status string                                 `json:"status" example:"success"`
	Errors string                                 `json:"errors" example:""`
	Data   CreateDestinationConnectorResponseData `json:"data"`
}

type DestinationDefinitionSpecification struct {
	DestinationDefinitionId       string                  `json:"destinationDefinitionId"`
	DocumentationUrl              string                  `json:"documentationUrl"`
	ConnectionSpecification       ConnectionSpecification `json:"connectionSpecification"`
	AuthSpecification             AuthSpecification       `json:"authSpecification"`
	AdvancedAuth                  AdvancedAuth            `json:"advancedAuth"`
	JobInfo                       JobInfo                 `json:"jobInfo"`
	SupportedDestinationSyncModes []string                `json:"supportedDestinationSyncModes"`
	SupportsDbt                   bool                    `json:"supportsDbt"`
	SupportsNormalization         bool                    `json:"supportsNormalization"`
}

type CreateDestinationConnectorRequest struct {
	AirbyteDestinationDefinitionId string         `json:"destinationDefinitionId" binding:"required" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	ConnectionConfiguration        datatypes.JSON `json:"connectionConfiguration" binding:"required" example:""`
	Name                           string         `json:"name" binding:"required" example:"example_destination_name"`
}

type CreateDestinationConnectorRequestAirbyte struct {
	CreateDestinationConnectorRequest
	WorkspaceId string `json:"workspaceId" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
}

type CreateDestinationConnectorResponseAirbyte struct {
	AirbyteDestinationId string `json:"destinationId"`
	DestinationName      string `json:"destinationName"`
	CreateDestinationConnectorRequestAirbyte
}

type DestinationDefinition struct {
	DestinationDefinitionID string `json:"destinationDefinitionId"`
	Name                    string `json:"name"`
	DockerRepository        string `json:"dockerRepository"`
	DockerImageTag          string `json:"dockerImageTag"`
	DocumentationURL        string `json:"documentationUrl"`
	Icon                    string `json:"icon"`
}

type DestinationDefinitions struct {
	DestinationDefinitions []DestinationDefinition `json:"DestinationDefinitions"`
}

type DestinationSpecification struct {
	DestinationDefinitionID       string      `json:"destinationDefinitionId"`
	DocumentationURL              string      `json:"documentationUrl"`
	ConnectionSpecification       interface{} `json:"connectionSpecification"`
	AuthSpecification             interface{} `json:"authSpecification"`
	AdvancedAuth                  interface{} `json:"advancedAuth"`
	JobInfo                       interface{} `json:"jobInfo"`
	SupportedDestinationSyncModes []string    `json:"supportedDestinationSyncModes"`
	SupportsDbt                   bool        `json:"supportsDbt"`
	SupportsNormalization         bool        `json:"supportsNormalization"`
}

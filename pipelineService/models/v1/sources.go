package models

import "github.com/gofrs/uuid"

type Source struct {
	SourceID                  string `json:"sourceId" binding:"required" gorm:"column:source_id; type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	SourceName                string `json:"sourceName" binding:"required" gorm:"column:name; type:string; default:(-)" example:"example_source"`
	AirbyteSourceID           string `json:"airbyteSourceId" gorm:"column:airbyte_source_id; type:uuid; default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	AirbyteSourceDefinitionID string `json:"airbyteSourceDefinitionId" gorm:"column:airbyte_source_definition_id; type:uuid ; default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	ConnectionID              string `json:"connectionId" binding:"required" gorm:"column:connection_id; type:uuid;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	Owner                     int    `json:"owner" gorm:"type:int" example:"1"`
	WorkspaceID               int    `json:"workspaceId" gorm:"type:int" example:"1"`
}

type CreateSourceConnectorRequestAPI struct {
	CreateSourceConnectorRequest
	Pipeline string `json:"pipelineId" binding:"required" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
}

type CreateSourceConnectorResposneData struct {
	SourceID     string `json:"sourceId" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	SourceName   string `json:"sourceName" example:"example_source"`
	ConnectionID string `json:"connectionId"  example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	PipelineID   string `json:"pipelineId" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
}

type CreateSourceConnectorResponseAPI struct {
	Status string                            `json:"status" example:"success"`
	Errors string                            `json:"errors" example:""`
	Data   CreateSourceConnectorResposneData `json:"data"`
}

type ConnectionSpecification struct {
	User User `json:"user"`
}

type Oauth2Specification struct {
	RootObject                []interface{} `json:"rootObject"`
	OauthFlowInitParameters   [][]string    `json:"oauthFlowInitParameters"`
	OauthFlowOutputParameters [][]string    `json:"oauthFlowOutputParameters"`
}

type AuthSpecification struct {
	AuthType            string              `json:"auth_type"`
	Oauth2Specification Oauth2Specification `json:"oauth2Specification"`
}

type OauthConfigSpecification struct {
	OauthUserInputFromConnectorConfigSpecification string `json:"oauthUserInputFromConnectorConfigSpecification"`
	CompleteOAuthOutputSpecification               string `json:"completeOAuthOutputSpecification"`
	CompleteOAuthServerInputSpecification          string `json:"completeOAuthServerInputSpecification"`
	CompleteOAuthServerOutputSpecification         string `json:"completeOAuthServerOutputSpecification"`
}

type AdvancedAuth struct {
	AuthFlowType             string                   `json:"authFlowType"`
	PredicateKey             []string                 `json:"predicateKey"`
	PredicateValue           string                   `json:"predicateValue"`
	OauthConfigSpecification OauthConfigSpecification `json:"oauthConfigSpecification"`
}
type Logs struct {
	LogLines []string `json:"logLines"`
}

type JobInfo struct {
	Id         string `json:"id"`
	ConfigType string `json:"configType"`
	ConfigId   string `json:"configId"`
	CreatedAt  int    `json:"createdAt"`
	EndedAt    int    `json:"endedAt"`
	Succeeded  bool   `json:"succeeded"`
	Logs       Logs   `json:"logs"`
}

type SourceDefinitionSpecification struct {
	SourceDefinitionId      string                  `json:"sourceDefinitionId"`
	DocumentationUrl        string                  `json:"documentationUrl"`
	ConnectionSpecification ConnectionSpecification `json:"connectionSpecification"`
	AuthSpecification       AuthSpecification       `json:"authSpecification"`
	AdvancedAuth            AdvancedAuth            `json:"advancedAuth"`
	JobInfo                 JobInfo                 `json:"jobInfo"`
}

type CreateSourceConnectorRequest struct {
	AirbyteSourceDefinitionId string      `json:"sourceDefinitionId" binding:"required" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	ConnectionConfiguration   interface{} `json:"connectionConfiguration" binding:"required" example:""`
	Name                      string      `json:"name" binding:"required" example:"example_source_name"`
}

type EditSourceConnectorRequest struct {
	ConnectionConfiguration interface{} `json:"connectionConfiguration" binding:"required" example:""`
}

type EditSourceConnectorRequestAirByte struct {
	AirByteSourceID         string      `json:"sourceId" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	ConnectionConfiguration interface{} `json:"connectionConfiguration" binding:"required" example:""`
	Name                    string      `json:"name" binding:"required" example:"example_source_name"`
}

type CreateSourceConnectorRequestAirbyte struct {
	CreateSourceConnectorRequest
	WorkspaceId string `json:"workspaceId" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
}

type CreateSourceConnectorResponseAirbyte struct {
	AirbyteSourceId string `json:"sourceId"`
	SourceName      string `json:"sourceName"`
	CreateSourceConnectorRequest
}

type SupportedSources struct {
	ID       uuid.UUID `json:"id" binding:"required" gorm:"type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	Name     string    `json:"name" binding:"required" gorm:"type:string;size:50" example:"source name"`
	Type     string    `json:"type" binding:"required" gorm:"type:string;size:50" example:"source type"`
	IsActive *bool     `json:"isActive" binding:"required" gorm:"default:(-)" example:"true"`
}

type SupportedSourcesResponse struct {
	Status string             `json:"status" example:"success"`
	Errors string             `json:"errors" example:""`
	Data   []SupportedSources `json:"data"`
}

type SourceDefinition struct {
	SourceDefinitionID string `json:"sourceDefinitionId"`
	Name               string `json:"name"`
	DockerRepository   string `json:"dockerRepository"`
	DockerImageTag     string `json:"dockerImageTag"`
	DocumentationURL   string `json:"documentationUrl"`
	Icon               string `json:"icon"`
}

type SourceDefinitions struct {
	SourceDefinitions []SourceDefinition `json:"sourceDefinitions"`
}

type SourceSpecification struct {
	SourceDefinitionID      string      `json:"sourceDefinitionId"`
	DocumentationURL        string      `json:"documentationUrl"`
	ConnectionSpecification interface{} `json:"connectionSpecification"`
	AuthSpecification       interface{} `json:"authSpecification"`
	AdvancedAuth            interface{} `json:"advancedAuth"`
	JobInfo                 interface{} `json:"jobInfo"`
}

type SourceSchema struct {
	Catalog struct {
		Streams []struct {
			Stream struct {
				Name                    string      `json:"name"`
				JSONSchema              interface{} `json:"jsonSchema"`
				SupportedSyncModes      []string    `json:"supportedSyncModes"`
				SourceDefinedCursor     bool        `json:"sourceDefinedCursor"`
				DefaultCursorField      []string    `json:"defaultCursorField"`
				SourceDefinedPrimaryKey [][]string  `json:"sourceDefinedPrimaryKey"`
				Namespace               interface{} `json:"namespace"`
			} `json:"stream"`
			Config struct {
				SyncMode            string     `json:"syncMode"`
				CursorField         []string   `json:"cursorField"`
				DestinationSyncMode string     `json:"destinationSyncMode"`
				PrimaryKey          [][]string `json:"primaryKey"`
				AliasName           string     `json:"aliasName"`
				Selected            bool       `json:"selected"`
			} `json:"config"`
		} `json:"streams"`
	} `json:"catalog"`
	JobInfo struct {
		ID         string `json:"id"`
		ConfigType string `json:"configType"`
		ConfigID   string `json:"configId"`
		CreatedAt  int    `json:"createdAt"`
		EndedAt    int    `json:"endedAt"`
		Succeeded  bool   `json:"succeeded"`
		Logs       struct {
			LogLines []string `json:"logLines"`
		} `json:"logs"`
	} `json:"jobInfo"`
}

type SourceSchemaResponse struct {
	Status string       `json:"status" example:"success"`
	Errors string       `json:"errors" example:""`
	Data   SourceSchema `json:"data"`
}

type ConfiguredSource struct {
	SourceDefinitionId      string      `json:"sourceDefinitionId"`
	SourceId                string      `json:"sourceId"`
	WorkspaceId             string      `json:"workspaceId"`
	ConnectionConfiguration interface{} `json:"connectionConfiguration"`
	Name                    string      `json:"name"`
	SourceName              string      `json:"sourceName"`
}

type ConfiguredSourceResponse struct {
	Status string           `json:"status" example:"success"`
	Errors string           `json:"errors" example:""`
	Data   ConfiguredSource `json:"data"`
}

package models

import (
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"gorm.io/datatypes"
)

type Pipeline struct {
	PipelineID         uuid.UUID      `json:"pipelineID" gorm:"type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	Name               string         `json:"name" binding:"required" gorm:"type:string;size:50" example:"pipeline-1"`
	PipelineGovernance pq.StringArray `json:"pipelineGovernance" binding:"required" gorm:"type:string[]" example:"[sales, marketing]"`
	CreatedAt          int64          `json:"createdAt" gorm:"default:(extract(epoch from now()) * 1000)"`
	Owner              int            `json:"owner" gorm:"type:int" example:"1"`
	WorkspaceID        int            `json:"workspaceId" gorm:"type:int" example:"1"`
}

type PipelineView struct {
	PipelineID          uuid.UUID                `json:"pipelineID" gorm:"column:pipeline_id;type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	Name                string                   `json:"name" binding:"required" gorm:"column:pipeline_name;type:string;size:50" example:"pipeline-1"`
	PipelineGovernance  pq.StringArray           `json:"pipelineGovernance" binding:"required" gorm:"type:string[]" example:"[sales, marketing]"`
	CreatedAt           int64                    `json:"createdAt" gorm:"default:(extract(epoch from now()) * 1000)"`
	ProductID           uuid.UUID                `json:"-" gorm:"column:product_id;type:uuid" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	Owner               int                      `json:"-" gorm:"column:pipeline_owner;type:int" example:"1"`
	SourceName          string                   `json:"sourceName" gorm:"column:source_name" example:"postgres"`
	SourceID            uuid.UUID                `json:"sourceID" gorm:"column:source_id" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	DestinationName     string                   `json:"destinationName" gorm:"column:destination_name"  example:"redshift"`
	DestinationID       string                   `json:"destinationId" binding:"required" gorm:"column:destination_id; type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	AirbyteStatus       string                   `json:"airbyteStatus" gorm:"column:airbyte_status"  example:"Inactive"`
	AirbyteLastRun      int                      `json:"airbyteLastRun" gorm:"column:airbyte_last_run" example:"1645517210"`
	AirbyteConnectionID string                   `json:"airbyteConnectionId" gorm:"column:airbyte_connection_id" example:""`
	ConnectionID        string                   `json:"connectionId" gorm:"column:connection_id" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	Product             []map[string]interface{} `json:"dataProduct"`
}

type UpdatePipeline struct {
	PipelineID         uuid.UUID      `json:"pipelineID" gorm:"type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	Name               string         `json:"name" binding:"required" gorm:"type:string;size:50" example:"pipeline-1"`
	PipelineGovernance pq.StringArray `json:"pipelineGovernance" binding:"required" gorm:"type:string[]" example:"[sales, marketing]"`
}

type GetPipelineDetails struct {
	Pipeline PipelineView `json:"pipeline"`
	Owner    interface{}  `json:"owner"`
}

type GetPipelineDetailsResponse struct {
	Status string             `json:"status" example:"success"`
	Errors string             `json:"errors" example:""`
	Data   GetPipelineDetails `json:"data"`
}

type GetAllConnectionsResponse struct {
	Status string       `json:"status" example:"success"`
	Errors string       `json:"errors" example:""`
	Data   []Connection `json:"data"`
}

type PipelineResponse struct {
	Status string   `json:"status" example:"success"`
	Errors string   `json:"errors" example:""`
	Data   Pipeline `json:"data"`
}

type PipelinesMetaData struct {
	PipelineID          string         `gorm:"column:pipeline_id" json:"pipelineID" example:"b251379e-01a1-11ec-82d6-a312edcd9c7b"`
	PipelineName        string         `gorm:"column:pipeline_name" json:"pipelineName" example:"pipeline-1"`
	PipelineGovernance  pq.StringArray `gorm:"column:pipeline_governance;type:varchar[]" json:"pipelineGovernance" example:"[sales, IoT]"`
	PipelineStatus      string         `gorm:"column:pipeline_status" json:"pipelineStatus" example:"Inactive"`
	SourceName          string         `gorm:"column:source_name" json:"sourceName" example:"postgres"`
	SourceID            uuid.UUID      `json:"sourceID" gorm:"column:source_id" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	DestinationName     string         `gorm:"column:destination_name" json:"destinationName" example:"redshift"`
	DestinationID       string         `json:"destinationId" binding:"required" gorm:"column:destination_id; type:uuid;primaryKey;default:(-)" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	AirbyteStatus       string         `gorm:"column:airbyte_status" json:"airbyteStatus" example:"Inactive"`
	AirbyteLastRun      int            `gorm:"column:airbyte_last_run" json:"airbyteLastRun" example:"1645517210"`
	AirbyteConnectionID string         `gorm:"column:airbyte_connection_id" json:"airbyteConnectionId" example:""`
	ConnectionID        string         `json:"connectionId" gorm:"column:connection_id" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	OwnerID             int            `gorm:"column:owner_id;type:int" json:"-" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	Owner               interface{}    `gorm:"-" json:"owner" example:"{}"`
}

type PipelineMetaDataResponse struct {
	Status string              `json:"status" example:"success"`
	Errors string              `json:"errors" example:""`
	Data   []PipelinesMetaData `json:"data"`
}

type PipelineSourceAndConnectionID struct {
	PipelineID          string `gorm:"column:pipeline_id" json:"pipelineID" example:"b251379e-01a1-11ec-82d6-a312edcd9c7b"`
	SourceID            string `gorm:"column:airbyte_source_id" json:"sourceID" example:"b251379e-01a1-11ec-82d6-a312edcd9c7b"`
	AirByteConnectionID string `gorm:"column:airbyte_connection_id" json:"airbyteConnectionId" example:""`
}

type GetPipelineSourceAndConnectionID struct {
	Status string                        `json:"status"`
	Errors string                        `json:"errors"`
	Data   PipelineSourceAndConnectionID `json:"data"`
}

type CreatePipelineRequest struct {
	SourceID      string       `json:"sourceId" binding:"required" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	DestinationID string       `json:"destinationId" binding:"required" example:"a152379e-01a1-11ec-82d6-a312edcd9c7b"`
	Schedule      *Schedule    `json:"schedule" example:""`
	SyncCatalog   SyncCatalog  `json:"syncCatalog" binding:"required" example:""`
	Prefix        *string      `json:"prefix" binding:"required" example:"t1"`
	Operations    []Operations `json:"operations,omitempty"`
}

type CreatePipelineResponse struct {
	Status string `json:"status" example:"success"`
	Errors string `json:"errors" example:""`
	Data   string `json:"data"`
}

type Stream struct {
	Name                    string      `json:"name"`
	JsonSchema              interface{} `json:"jsonSchema"`
	SupportedSyncModes      []string    `json:"supportedSyncModes"`
	SourceDefinedCursor     bool        `json:"sourceDefinedCursor"`
	DefaultCursorField      []string    `json:"defaultCursorField"`
	SourceDefinedPrimaryKey [][]string  `json:"sourceDefinedPrimaryKey"`
	Namespace               interface{} `json:"namespace"`
}

type Config struct {
	SyncMode            string     `json:"syncMode"`
	CursorField         []string   `json:"cursorField"`
	DestinationSyncMode string     `json:"destinationSyncMode"`
	PrimaryKey          [][]string `json:"primaryKey"`
	AliasName           string     `json:"aliasName"`
	Selected            bool       `json:"selected"`
}

type Streams struct {
	Stream Stream `json:"stream"`
	Config Config `json:"config"`
}

type SyncCatalog struct {
	Streams []Streams `json:"streams"`
}

type Schedule struct {
	Units    int    `json:"units" example:"1"`
	TimeUnit string `json:"timeUnit" example:"hours"`
}
type ResourceRequirements struct {
	CpuRequest    string `json:"cpu_request"`
	CpuLimit      string `json:"cpu_limit"`
	MemoryRequest string `json:"memory_request"`
	MemoryLimit   string `json:"memory_limit"`
}

type Normalization struct {
	Option string `json:"option"`
}

type Dbt struct {
	GitRepoUrl    string `json:"gitRepoUrl"`
	GitRepoBranch string `json:"gitRepoBranch"`
	DockerImage   string `json:"dockerImage"`
	DbtArguments  string `json:"dbtArguments"`
}

type OperatorConfiguration struct {
	OperatorType  string         `json:"operatorType"`
	Normalization *Normalization `json:"normalization"`
	Dbt           *Dbt           `json:"dbt"`
}

type Operations struct {
	WorkspaceId           string                `json:"workspaceId"`
	Name                  string                `json:"name"`
	OperatorConfiguration OperatorConfiguration `json:"operatorConfiguration"`
}

type CreatePipelineAirbyteRequest struct {
	//Name                 string               `json:"name"`
	NamespaceDefinition string `json:"namespaceDefinition"`
	NamespaceFormat     string `json:"namespaceFormat"`
	Prefix              string `json:"prefix"`
	SourceId            string `json:"sourceId"`
	DestinationId       string `json:"destinationId"`
	//OperationIds         []string             `json:"operationIds"`
	SyncCatalog SyncCatalog `json:"syncCatalog"`
	Schedule    *Schedule   `json:"schedule"`
	Status      string      `json:"status"`
	//ResourceRequirements ResourceRequirements `json:"resourceRequirements"`
	Operations []Operations `json:"operations,omitempty"`
}

type UpdatePipelineAirByteRequest struct {
	ConnectionId string       `json:"connectionId"`
	Prefix       *string      `json:"prefix" binding:"required"`
	SyncCatalog  SyncCatalog  `json:"syncCatalog" binding:"required" example:""`
	Schedule     *Schedule    `json:"schedule"  example:""`
	Status       string       `json:"status" binding:"required" example:"active"`
	Operations   []Operations `json:"operations,omitempty"`
}

type AirbyteSourceAndDestinations struct {
	ConnectionID                    string         `json:"connectionId"`
	SourceID                        string         `json:"sourceId"`
	DestinationID                   string         `json:"destinationId"`
	DestinationType                 string         `json:"destinationType"`
	DestinationConfigurationDetails datatypes.JSON `json:"configurationDetails"`
	AirbyteSourceID                 string         `json:"airbyteSourceId"`
	AirbyteDestinationID            string         `json:"airbyteDestinationId"`
	PipelineName                    string         `json:"pipelineName" gorm:"column:pipeline_name"`
	PipelineID                      string         `json:"pipelineId" gorm:"column:pipeline_id"`
}

type CreatePipelineAirbyteResponse struct {
	ConnectionId           string                                    `json:"connectionId"`
	Name                   string                                    `json:"name"`
	NamespaceDefinition    string                                    `json:"namespaceDefinition"`
	NamespaceFormat        string                                    `json:"namespaceFormat"`
	Prefix                 string                                    `json:"prefix"`
	SourceId               string                                    `json:"sourceId"`
	DestinationId          string                                    `json:"destinationId"`
	SyncCatalog            SyncCatalog                               `json:"syncCatalog"`
	Schedule               *Schedule                                 `json:"schedule"`
	Status                 string                                    `json:"status"`
	OperationIds           []interface{}                             `json:"operationIds"`
	Source                 CreateSourceConnectorResponseAirbyte      `json:"source"`
	Destination            CreateDestinationConnectorResponseAirbyte `json:"destination"`
	Operations             []Operations                              `json:"operations,omitempty"`
	LatestSyncJobCreatedAt interface{}                               `json:"latestSyncJobCreatedAt"`
	LatestSyncJobStatus    interface{}                               `json:"latestSyncJobStatus"`
	IsSyncing              bool                                      `json:"isSyncing"`
	ResourceRequirements   ResourceRequirements                      `json:"resourceRequirements"`
}

type PipelineSchemas struct {
	SchemaID   string `gorm:"type:uuid;primaryKey;default:(-)" json:"schemaID" example:"b251379e-01a1-11ec-82d6-a312edcd9c7b"`
	PipelineID string `gorm:"column:pipeline_id" json:"pipelineID" example:"b251379e-01a1-11ec-82d6-a312edcd9c7b"`
	Name       string `gorm:"column:name" json:"name" example:"schema"`
	Prefix     string `gorm:"column:prefix" json:"prefix" example:"schema"`
}

type SourceAndPipelineName struct {
	ConnectionID    string `json:"connectionId"`
	SourceID        string `json:"sourceId"`
	AirbyteSourceID string `json:"airbyteSourceId"`
	PipelineName    string `json:"pipelineName"`
	PipelineID      string `json:"pipelineID"`
}

type TransformationPipelines struct {
	TransformationPipelineId string `json:"transformationPipelineId" gorm:"column:transformation_pipeline_id; type:uuid;primaryKey;default:(-)"`
	ProductID                string `json:"productId" gorm:"column:product_id; type:uuid;default:(-)"`
	SourceID                 string `json:"sourceId" gorm:"column:source_id; type:uuid;default:(-)"`
	DestinationID            string `json:"destinationId" gorm:"column:destination_id; type:uuid;default:(-)"`
	AirbyteConnectionID      string `json:"airbyteConnectionId" gorm:"column:airbyte_connection_id; type:uuid;"`
}

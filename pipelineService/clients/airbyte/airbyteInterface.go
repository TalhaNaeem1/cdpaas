//go:generate mockgen -destination=mocks/mock_airbyte.go -package=mock_airbyte . AirByteQuerier
package airbyte

import (
	"pipelineService/models/v1"
)

type AirByteQuerier interface {
	CreateWorkspace(workspace models.WorkspaceRequest) (models.WorkspaceAPIResponse, error)
	GetWorkspaceID() (string, error)
	CreateSourceConnectorOnAirByte(airbyte models.CreateSourceConnectorRequestAirbyte) (models.CreateSourceConnectorResponseAirbyte, error)
	EditSourceConnectorOnAirByte(requestBody models.EditSourceConnectorRequestAirByte) (models.CreateSourceConnectorResponseAirbyte, error)
	GetSourceDefinitions() (models.SourceDefinitions, error)
	GetConfiguredSource(sourceId string) (models.ConfiguredSource, error)
	GetSourceSpecification(sourceDefinitionID string) (models.SourceSpecification, error)
	CreateConnection(request models.CreatePipelineAirbyteRequest) (models.CreatePipelineAirbyteResponse, error)
	UpdateConnection(request models.UpdatePipelineAirByteRequest) (models.CreatePipelineAirbyteResponse, error)
	DiscoverSourceSchema(sourceId string) (models.SourceSchema, error)
	CreateDestinationConnectorOnAirByte(airbyte models.CreateDestinationConnectorRequestAirbyte) (models.CreateDestinationConnectorResponseAirbyte, error)
	GetDestinationDefinitions() (models.DestinationDefinitions, error)
	GetDestinationSpecification(destinationDefinitionID string) (models.DestinationSpecification, error)
	GetConnectionDetails(connection map[string]interface{}) (models.ConnectionMeta, error)
	SyncConnectionManually(requestBody map[string]interface{}) (models.ManualConnectionSyncResponse, error)
	FetchSyncHistory(request models.SyncHistoryRequest) (models.SyncHistoryResponse, error)
	GetJobLogs(jobID int) (models.JobLogs, error)
	GetConnectionSchema(connectionID string) (models.ConnectionSourceSchema, error)
	GetConnectionSummary(connectionID string) (models.ConnectionSummaryAirByte, error)
	CheckDestinationConnection(requestBody map[string]interface{}) error
	CheckSourceConnection(requestBody map[string]interface{}) error
	GetConnection(requestBody map[string]interface{}) ([]byte, error)
}

var _ AirByteQuerier = (*RequestMaker)(nil)

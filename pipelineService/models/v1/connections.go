package models

type Connection struct {
	ConnectionID          string `json:"connectionId" gorm:"column:connection_id; type:uuid;primaryKey;default:(-)"`
	PipelineID            string `json:"pipelineId" gorm:"column:pipeline_id; type:uuid;default:(-)"`
	CreatedAt             int64  `json:"createdAt" gorm:"default:(extract(epoch from now()) * 1000)"`
	AirbyteStatus         string `json:"status" gorm:"column:airbyte_status; type:uuid; default:(-)"`
	AirbyteLastRun        int    `json:"airbyteLastRun" gorm:"column:airbyte_last_run;"`
	AirbyteConnectionID   string `json:"airbyteConnectionId" gorm:"column:airbyte_connection_id; type:uuid;"`
	AirbyteFrequencyUnits int    `json:"airbyteFrequencyUnits" gorm:"column:airbyte_frequency_units;"`
	AirbyteTimeUnit       string `json:"airbyteTimeUnit" gorm:"column:airbyte_time_unit;"`
	IsFirstRun            bool   `json:"isFirstRun" gorm:"type:bool" example:"false"`
	FirstSucceededRunAt   int    `json:"firstSucceededRunAt" gorm:"column:first_succeeded_run_at;"`
	Owner                 int    `json:"owner" gorm:"type:int" example:"1"`
	WorkspaceID           int    `json:"workspaceId" gorm:"type:int" example:"1"`
}

type PipelineConnection struct {
	ConnectionID string `json:"connectionId" gorm:"column:connection_id; type:uuid;primaryKey;default:(-)"`
	PipelineID   string `json:"pipelineId" gorm:"column:pipeline_id; type:uuid;default:(-)"`
	SourceID     string `json:"sourceId" gorm:"column:source_id; type:uuid;default:(-)"`
	SourceName   string `json:"sourceName" gorm:"column:source_name; type:string;default:(-)"`
}

type UpdateConnection struct {
	Connection    Connection
	DestinationID string `json:"destinationId"`
}

type ConnectionsDestinations struct {
	ConnectionID  string `json:"connectionId" gorm:"column:connection_id; type:uuid;default:(-)"`
	DestinationID string `json:"destinationId" gorm:"column:destination_id; type:uuid;default:(-)" `
}

type ConnectionMeta struct {
	LatestSyncJobCreatedAt int    `json:"latestSyncJobCreatedAt"`
	LatestSyncJobStatus    string `json:"latestSyncJobStatus"`
}

type SyncHistoryRequest struct {
	ConfigTypes []string `json:"configTypes"`
	ConfigId    string   `json:"configId"`
}

type Job struct {
	ID         int    `json:"id"`
	ConfigType string `json:"configType"`
	ConfigID   string `json:"configId"`
	CreatedAt  int    `json:"createdAt"`
	UpdatedAt  int    `json:"updatedAt"`
	Status     string `json:"status"`
}

type ManualConnectionSyncResponse struct {
	Job      Job           `json:"job"`
	Attempts []interface{} `json:"attempts"`
}

type SyncHistoryResponse struct {
	Jobs []struct {
		Job      Job `json:"job"`
		Attempts []struct {
			Id            int    `json:"id"`
			Status        string `json:"status"`
			CreatedAt     int    `json:"createdAt"`
			UpdatedAt     int    `json:"updatedAt"`
			EndedAt       int    `json:"endedAt"`
			BytesSynced   int    `json:"bytesSynced"`
			RecordsSynced int    `json:"recordsSynced"`
			TotalStats    struct {
				RecordsEmitted       int `json:"recordsEmitted"`
				BytesEmitted         int `json:"bytesEmitted"`
				StateMessagesEmitted int `json:"stateMessagesEmitted"`
				RecordsCommitted     int `json:"recordsCommitted"`
			} `json:"totalStats"`
			StreamStats []struct {
				StreamName string `json:"streamName"`
				Stats      struct {
					RecordsEmitted       int         `json:"recordsEmitted"`
					BytesEmitted         int         `json:"bytesEmitted"`
					StateMessagesEmitted interface{} `json:"stateMessagesEmitted"`
					RecordsCommitted     int         `json:"recordsCommitted"`
				} `json:"stats"`
			} `json:"streamStats"`
			FailureSummary interface{} `json:"failureSummary"`
		} `json:"attempts"`
	} `json:"jobs"`
}

type JobLogs struct {
	Attempts []struct {
		Logs struct {
			LogLines []string `json:"logLines"`
		} `json:"logs"`
	} `json:"attempts"`
}

type ConnectionSourceSchema struct {
	ConnectionId        string      `json:"connectionId"`
	Name                string      `json:"name"`
	NamespaceDefinition string      `json:"namespaceDefinition"`
	NamespaceFormat     string      `json:"namespaceFormat"`
	Prefix              string      `json:"prefix"`
	SourceId            string      `json:"sourceId"`
	DestinationId       string      `json:"destinationId"`
	SyncCatalog         SyncCatalog `json:"syncCatalog"`
	Schedule            Schedule    `json:"schedule"`
	Status              string      `json:"status"`
	OperationIds        []string    `json:"operationIds"`
	Source              struct {
		SourceDefinitionId      string `json:"sourceDefinitionId"`
		SourceId                string `json:"sourceId"`
		WorkspaceId             string `json:"workspaceId"`
		ConnectionConfiguration struct {
			User string `json:"user"`
		} `json:"connectionConfiguration"`
		Name       string `json:"name"`
		SourceName string `json:"sourceName"`
	} `json:"source"`
	Destination struct {
		DestinationDefinitionId string `json:"destinationDefinitionId"`
		DestinationId           string `json:"destinationId"`
		WorkspaceId             string `json:"workspaceId"`
		ConnectionConfiguration struct {
			User string `json:"user"`
		} `json:"connectionConfiguration"`
		Name            string `json:"name"`
		DestinationName string `json:"destinationName"`
	} `json:"destination"`
	Operations             []Operations         `json:"operations"`
	LatestSyncJobCreatedAt int                  `json:"latestSyncJobCreatedAt"`
	LatestSyncJobStatus    string               `json:"latestSyncJobStatus"`
	IsSyncing              bool                 `json:"isSyncing"`
	ResourceRequirements   ResourceRequirements `json:"resourceRequirements"`
}

type ConnectionSummaryResponse struct {
	AirByteSummary    ConnectionSummaryAirByte
	SourceName        string      `json:"sourceName"`
	Owner             interface{} `json:"configuredBy"`
	ConfigurationDate int64       `json:"configurationDate"`
}

type ConnectionSummaryAirByte struct {
	Schedule struct {
		Units    int    `json:"units"`
		TimeUnit string `json:"timeUnit"`
	} `json:"schedule"`
	ConnectionStatus string `json:"status"`
}

type ConnectionSummary struct {
	AirbyteConnectionID string `json:"airbyteConnectionId" gorm:"column:airbyte_connection_id; type:uuid;"`
	SourceName          string `json:"sourceName" binding:"required" gorm:"column:name; type:string; default:(-)" example:"example_source"`
	Owner               int    `json:"owner" gorm:"type:int" example:"1"`
	CreatedAt           int64  `json:"createdAt" gorm:"column:created_at;"`
}

type CheckConnection struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type EnableAssetsInternalRequest struct {
	ConnectionIDs []string `json:"connectionIDs"`
}

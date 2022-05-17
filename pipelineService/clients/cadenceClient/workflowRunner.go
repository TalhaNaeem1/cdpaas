package cadenceClient

import (
	"context"

	"go.uber.org/cadence/client"
	"pipelineService/models/v1"
)

type WorkflowRunner interface {
	TriggerDeletePipelineWorkflow(ctx context.Context, workFlowOptions client.StartWorkflowOptions, pipelineID string) error
	TriggerSendEmailWorkflow(ctx context.Context, emailTemplate models.EmailTemplate, wfOptions client.StartWorkflowOptions) error
	TriggerCreateConnectionWorkflow(ctx context.Context, workFlowOptions client.StartWorkflowOptions,
		createPipelineRequest models.CreatePipelineRequest, connectionInfo models.AirbyteSourceAndDestinations,
		userID int, workspaceID int, airbyteWorkspaceID string) error
	TriggerUpdateConnectionWorkflow(ctx context.Context, workFlowOptions client.StartWorkflowOptions, updatePipelineRequest models.UpdatePipelineAirByteRequest,
		connection models.Connection, userID int, workspaceID int, airbyteWorkspaceID string) error
}

var _ WorkflowRunner = (*Workflows)(nil)

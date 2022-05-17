package cadenceClient

import (
	"context"

	"go.uber.org/cadence/client"
	"pipelineService/env"
	"pipelineService/models/v1"
	"pipelineService/utils"
)

var (
	SendEmailWorkflow        = env.Env.CadenceWorkerServiceName + "/v1/workflows/authService.SendEmailWorkflow"
	DeletePipelineWorkflow   = env.Env.CadenceWorkerServiceName + "/v1/workflows/pipeline.DeletePipelineWorkflow"
	CreateConnectionWorkflow = env.Env.CadenceWorkerServiceName + "/v1/workflows/pipeline.CreateConnectionWorkflow"
	UpdateConnectionWorkflow = env.Env.CadenceWorkerServiceName + "/v1/workflows/pipeline.UpdateConnectionWorkflow"
)

func (wr *Workflows) TriggerSendEmailWorkflow(ctx context.Context, emailTemplate models.EmailTemplate, wfOptions client.StartWorkflowOptions) error {
	_, err := wr.client.ExecuteWorkflow(ctx, wfOptions, SendEmailWorkflow, emailTemplate)
	if err != nil {
		return err
	}

	return nil
}

func (wr *Workflows) TriggerDeletePipelineWorkflow(ctx context.Context, workFlowOptions client.StartWorkflowOptions, pipelineID string) error {
	logger := utils.GetLogger()
	logger.Info("TriggerDeletePipelineWorkflow endpoint called")

	_, err := wr.client.ExecuteWorkflow(ctx, workFlowOptions, DeletePipelineWorkflow, pipelineID)
	if err != nil {
		return err
	}

	return err
}

func (wr *Workflows) TriggerCreateConnectionWorkflow(ctx context.Context, workFlowOptions client.StartWorkflowOptions,
	createPipelineRequest models.CreatePipelineRequest, connectionInfo models.AirbyteSourceAndDestinations,
	userID int, workspaceID int, airbyteWorkspaceID string) error {
	logger := utils.GetLogger()
	logger.Info("TriggerCreateConnectionWorkflow endpoint called")

	workflowRun, err := wr.client.ExecuteWorkflow(ctx, workFlowOptions, CreateConnectionWorkflow, createPipelineRequest, connectionInfo, userID, workspaceID, airbyteWorkspaceID)
	if err != nil {
		return err
	}

	var emptyInterface interface{}
	err = workflowRun.Get(ctx, &emptyInterface)

	return err
}

func (wr *Workflows) TriggerUpdateConnectionWorkflow(ctx context.Context, workFlowOptions client.StartWorkflowOptions,
	updatePipelineRequest models.UpdatePipelineAirByteRequest, connection models.Connection,
	userID int, workspaceID int, airbyteWorkspaceID string) error {
	logger := utils.GetLogger()
	logger.Info("TriggerUpdateConnectionWorkflow endpoint called")

	workflowRun, err := wr.client.ExecuteWorkflow(ctx, workFlowOptions, UpdateConnectionWorkflow, updatePipelineRequest, connection, userID, workspaceID, airbyteWorkspaceID)
	if err != nil {
		return err
	}

	var emptyInterface interface{}
	err = workflowRun.Get(ctx, &emptyInterface)

	return err
}

package cadenceClient

import (
	"context"

	"go.uber.org/cadence/client"
)

type CadClient interface {
	ExecuteWorkflow(ctx context.Context, options client.StartWorkflowOptions, workflow interface{}, args ...interface{}) (client.WorkflowRun, error)
}

type Workflows struct {
	client CadClient
}

func NewRunner(wfRunner CadClient) *Workflows {
	return &Workflows{client: wfRunner}
}

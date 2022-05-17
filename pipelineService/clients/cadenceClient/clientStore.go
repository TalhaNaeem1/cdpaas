package cadenceClient

import "go.uber.org/cadence/client"

type CadStore interface {
	WorkflowRunner
}

type StoreForCadenceClient struct {
	cadClient *client.Client
	*Workflows
}

func NewStore(cc *client.Client)  CadStore{
	return &StoreForCadenceClient{
		cadClient: cc,
		Workflows: NewRunner(*cc),
	}
}

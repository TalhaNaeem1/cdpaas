package cadenceClient

import (
	"go.uber.org/cadence/.gen/go/cadence/workflowserviceclient"
	"go.uber.org/cadence/client"
	"go.uber.org/yarpc"
	"go.uber.org/yarpc/transport/tchannel"
	"pipelineService/env"
)

var HostPort = env.Env.CadenceService
var Domain = env.Env.CadenceDomain
var ClientName = env.Env.ClientName
var CadenceService = env.Env.CadenceServiceName

func buildCadenceServiceClient() (workflowserviceclient.Interface, error) {
	ch, err := tchannel.NewChannelTransport(tchannel.ServiceName(ClientName))
	if err != nil {
		return nil, err
	}

	dispatcher := yarpc.NewDispatcher(yarpc.Config{
		Name: ClientName,
		Outbounds: yarpc.Outbounds{
			CadenceService: {Unary: ch.NewSingleOutbound(HostPort)},
		},
	})

	if err := dispatcher.Start(); err != nil {
		return nil, err
	}

	return workflowserviceclient.New(dispatcher.ClientConfig(CadenceService)), nil
}

func GetNewCadenceClient() (client.Client, error) {
	service, err := buildCadenceServiceClient()
	if err != nil {
		return nil, err
	}

	return client.NewClient(service, Domain, &client.Options{}), nil
}

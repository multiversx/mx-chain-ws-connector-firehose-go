package factory

import (
	"github.com/multiversx/mx-chain-communication-go/websocket/data"
	factoryHost "github.com/multiversx/mx-chain-communication-go/websocket/factory"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-core-go/marshal/factory"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
)

var log = logger.GetOrCreate("ws-connector")

// CreateWSConnector will create a ws connector able to receive and process incoming data
// from a multiversx node
func CreateWSConnector(cfg config.WebSocketConfig) (process.WSConnector, error) {
	marshaller, err := factory.NewMarshalizer(cfg.MarshallerType)
	if err != nil {
		return nil, err
	}

	dataProcessor, err := process.NewLogDataProcessor(marshaller, log)
	if err != nil {
		return nil, err
	}

	wsHost, err := createWsHost(marshaller, cfg)
	if err != nil {
		return nil, err
	}

	err = wsHost.SetPayloadHandler(dataProcessor)
	if err != nil {
		return nil, err
	}

	return wsHost, nil
}

func createWsHost(wsMarshaller marshal.Marshalizer, cfg config.WebSocketConfig) (factoryHost.FullDuplexHost, error) {
	return factoryHost.CreateWebSocketHost(factoryHost.ArgsWebSocketHost{
		WebSocketConfig: data.WebSocketConfig{
			URL:                        cfg.Url,
			WithAcknowledge:            cfg.WithAcknowledge,
			Mode:                       cfg.Mode,
			RetryDurationInSec:         int(cfg.RetryDuration),
			BlockingAckOnError:         cfg.BlockingAckOnError,
			DropMessagesIfNoConnection: false,
		},
		Marshaller: wsMarshaller,
		Log:        log,
	})
}

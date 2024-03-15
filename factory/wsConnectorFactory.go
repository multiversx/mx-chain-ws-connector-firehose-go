package factory

import (
	"os"

	"github.com/multiversx/mx-chain-communication-go/websocket/data"
	factoryHost "github.com/multiversx/mx-chain-communication-go/websocket/factory"
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-core-go/marshal/factory"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-storage-go/storageUnit"

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

	blockContainer, err := createBlockContainer()
	if err != nil {
		return nil, err
	}

	protoMarshaller := &marshal.GogoProtoMarshalizer{}

	firehosePublisher, err := process.NewFirehosePublisher(
		os.Stdout, // DO NOT CHANGE
		blockContainer,
		protoMarshaller,
	)
	if err != nil {
		return nil, err
	}

	// TODO: move cache to config
	cacheConfig := storageUnit.CacheConfig{
		Type:        storageUnit.SizeLRUCache,
		SizeInBytes: 209715200, // 200MB
		Capacity:    100,
	}

	cacher, err := storageUnit.NewCache(cacheConfig)
	if err != nil {
		return nil, err
	}

	blocksPool, err := process.NewBlocksPool(cacher, blockContainer, protoMarshaller)
	if err != nil {
		return nil, err
	}

	dataAggregator, err := process.NewDataAggregator(blocksPool)
	if err != nil {
		return nil, err
	}

	// TODO: move to separate factory
	dataProcessor, err := process.NewDataProcessor(firehosePublisher, protoMarshaller, blocksPool, dataAggregator)
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

func createBlockContainer() (process.BlockContainerHandler, error) {
	container := block.NewEmptyBlockCreatorsContainer()

	err := container.Add(core.ShardHeaderV1, block.NewEmptyHeaderCreator())
	if err != nil {
		return nil, err
	}
	err = container.Add(core.ShardHeaderV2, block.NewEmptyHeaderV2Creator())
	if err != nil {
		return nil, err
	}
	err = container.Add(core.MetaHeader, block.NewEmptyMetaBlockCreator())
	if err != nil {
		return nil, err
	}

	return container, nil
}

func createWsHost(wsMarshaller marshal.Marshalizer, cfg config.WebSocketConfig) (factoryHost.FullDuplexHost, error) {
	return factoryHost.CreateWebSocketHost(factoryHost.ArgsWebSocketHost{
		WebSocketConfig: data.WebSocketConfig{
			URL:                        cfg.Url,
			WithAcknowledge:            cfg.WithAcknowledge,
			Mode:                       cfg.Mode,
			RetryDurationInSec:         int(cfg.RetryDuration),
			BlockingAckOnError:         cfg.BlockingAckOnError,
			DropMessagesIfNoConnection: cfg.DropMessagesIfNoConnection,
			AcknowledgeTimeoutInSec:    cfg.AcknowledgeTimeoutInSec,
			Version:                    cfg.Version,
		},
		Marshaller: wsMarshaller,
		Log:        log,
	})
}

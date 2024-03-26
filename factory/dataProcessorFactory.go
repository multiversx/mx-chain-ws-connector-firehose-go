package factory

import (
	"os"

	"github.com/multiversx/mx-chain-communication-go/websocket"
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-storage-go/storageUnit"
	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
)

// CreateDataProcessor will create a new instance of data processor
func CreateDataProcessor(cfg config.Config, storer process.PruningStorer) (websocket.PayloadHandler, error) {
	protoMarshaller := &marshal.GogoProtoMarshalizer{}

	blockContainer, err := createBlockContainer()
	if err != nil {
		return nil, err
	}

	firehosePublisher, err := process.NewFirehosePublisher(
		os.Stdout,
		blockContainer,
		protoMarshaller,
	)
	if err != nil {
		return nil, err
	}

	blocksPool, err := process.NewBlocksPool(storer, protoMarshaller, cfg.DataPool.NumberOfShards, cfg.DataPool.MaxDelta, cfg.DataPool.PruningWindow)
	if err != nil {
		return nil, err
	}

	dataAggregator, err := process.NewDataAggregator(blocksPool)
	if err != nil {
		return nil, err
	}

	return process.NewDataProcessor(firehosePublisher, protoMarshaller, blocksPool, dataAggregator, blockContainer)
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

// CreateStorer will create a new pruning storer instace
func CreateStorer(cfg config.Config, importDBMode bool) (process.PruningStorer, error) {
	cacheConfig := storageUnit.CacheConfig{
		Type:        storageUnit.CacheType(cfg.OutportBlocksStorage.Cache.Type),
		SizeInBytes: cfg.OutportBlocksStorage.Cache.SizeInBytes,
		Capacity:    cfg.OutportBlocksStorage.Cache.Capacity,
	}

	cacher, err := storageUnit.NewCache(cacheConfig)
	if err != nil {
		return nil, err
	}

	if importDBMode {
		return process.NewImportDBStorer(cacher)
	}

	return process.NewPruningStorer(cfg.OutportBlocksStorage.DB, cacher, cfg.DataPool.NumPersistersToKeep)
}

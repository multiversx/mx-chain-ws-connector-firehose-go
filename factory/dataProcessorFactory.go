package factory

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-storage-go/storageUnit"

	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process/dataPool"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process/dataPool/disabled"
)

// CreateBlockContainer will create a new block container component
func CreateBlockContainer() (process.BlockContainerHandler, error) {
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

// CreateBlocksPool will create a new blocks pool component
func CreateBlocksPool(
	cfg config.Config,
	importDBMode bool,
	marshaller marshal.Marshalizer,
) (process.DataPool, error) {
	blocksStorer, err := CreateStorer(cfg, importDBMode)
	if err != nil {
		return nil, err
	}

	return dataPool.NewBlocksPool(blocksStorer, marshaller, cfg.DataPool.NumberOfShards, cfg.DataPool.MaxDelta, cfg.DataPool.PruningWindow)
}

// CreateHyperBlocksPool will create a new hyper blocks pool component
func CreateHyperBlocksPool(
	grpcServerMode bool,
	cfg config.Config,
	importDBMode bool,
	marshaller marshal.Marshalizer,
) (process.HyperOutportBlocksPool, error) {
	if !grpcServerMode {
		return disabled.NewDisabledHyperOutportBlocksPool(), nil
	}

	hyperOutportBlockDataPool, err := CreateBlocksPool(cfg, importDBMode, marshaller)
	if err != nil {
		return nil, err
	}

	return dataPool.NewHyperOutportBlocksPool(hyperOutportBlockDataPool, marshaller)
}

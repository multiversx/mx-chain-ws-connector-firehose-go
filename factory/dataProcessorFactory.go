package factory

import (
	"errors"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-storage-go/storageUnit"
	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process/dataPool"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process/dataPool/disabled"
)

// ErrNotSupportedDBMode signals that an invalid db mode was provided
var ErrNotSupportedDBMode = errors.New("not supported db mode")

const (
	FullPersisterDBMode      = "full-persister"
	ImportDBMode             = "import-db"
	OptimizedPersisterDBMode = "optimized-persister"
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
func CreateStorer(cfg config.Config, dbMode string) (process.PruningStorer, error) {
	cacheConfig := storageUnit.CacheConfig{
		Type:        storageUnit.CacheType(cfg.OutportBlocksStorage.Cache.Type),
		SizeInBytes: cfg.OutportBlocksStorage.Cache.SizeInBytes,
		Capacity:    cfg.OutportBlocksStorage.Cache.Capacity,
	}

	cacher, err := storageUnit.NewCache(cacheConfig)
	if err != nil {
		return nil, err
	}

	switch dbMode {
	case FullPersisterDBMode:
		return process.NewPruningStorer(cfg.OutportBlocksStorage.DB, cacher, cfg.DataPool.NumPersistersToKeep, true)
	case OptimizedPersisterDBMode:
		return process.NewPruningStorer(cfg.OutportBlocksStorage.DB, cacher, cfg.DataPool.NumPersistersToKeep, false)
	case ImportDBMode:
		return process.NewImportDBStorer(cacher)
	default:
		return nil, ErrNotSupportedDBMode
	}
}

// CreateBlocksPool will create a new blocks pool component
func CreateBlocksPool(
	cfg config.Config,
	dbMode string,
	marshaller marshal.Marshalizer,
) (process.DataPool, error) {
	blocksStorer, err := CreateStorer(cfg, dbMode)
	if err != nil {
		return nil, err
	}

	return dataPool.NewBlocksPool(blocksStorer, marshaller, cfg.DataPool.NumberOfShards, cfg.DataPool.MaxDelta, cfg.DataPool.PruningWindow)
}

// CreateHyperBlocksPool will create a new hyper blocks pool component
func CreateHyperBlocksPool(
	grpcServerMode bool,
	cfg config.Config,
	dbMode string,
	marshaller marshal.Marshalizer,
) (process.HyperOutportBlocksPool, error) {
	if !grpcServerMode {
		return disabled.NewDisabledHyperOutportBlocksPool(), nil
	}

	hyperOutportBlockDataPool, err := CreateBlocksPool(cfg, dbMode, marshaller)
	if err != nil {
		return nil, err
	}

	return dataPool.NewHyperOutportBlocksPool(hyperOutportBlockDataPool, marshaller)
}

package factory

import (
	"errors"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	storageCommon "github.com/multiversx/mx-chain-storage-go/common"
	"github.com/multiversx/mx-chain-storage-go/factory"

	"github.com/multiversx/mx-chain-ws-connector-firehose-go/common"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/config"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
)

// ErrNotSupportedDBMode signals that an invalid db mode was provided
var ErrNotSupportedDBMode = errors.New("not supported db mode")

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
func CreateStorer(cfg config.Config, dbMode common.DBMode) (process.PruningStorer, error) {
	cacheConfig := storageCommon.CacheConfig{
		Type:        storageCommon.CacheType(cfg.OutportBlocksStorage.Cache.Type),
		SizeInBytes: cfg.OutportBlocksStorage.Cache.SizeInBytes,
		Capacity:    cfg.OutportBlocksStorage.Cache.Capacity,
	}

	cacher, err := factory.NewCache(cacheConfig)
	if err != nil {
		return nil, err
	}

	switch dbMode {
	case common.FullPersisterDBMode:
		return process.NewPruningStorer(cfg.OutportBlocksStorage.DB, cacher, cfg.DataPool.NumPersistersToKeep, true)
	case common.OptimizedPersisterDBMode:
		return process.NewPruningStorer(cfg.OutportBlocksStorage.DB, cacher, cfg.DataPool.NumPersistersToKeep, false)
	case common.ImportDBMode:
		return process.NewImportDBStorer(cacher)
	default:
		return nil, ErrNotSupportedDBMode
	}
}

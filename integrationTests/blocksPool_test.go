package integrationtests

import (
	"fmt"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-template-go/common"
	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
	"github.com/multiversx/mx-chain-ws-connector-template-go/factory"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/stretchr/testify/require"
)

func getDefaultConfig() config.Config {
	maxDelta := uint64(10)
	pruningWindow := uint64(25)

	cfg := config.Config{
		DataPool: config.DataPoolConfig{
			MaxDelta:             maxDelta,
			PruningWindow:        pruningWindow,
			NumPersistersToKeep:  2,
			FirstCommitableBlock: 0,
		},
		OutportBlocksStorage: config.StorageConfig{
			Cache: config.CacheConfig{
				Name:        "OutportBlocksStorage",
				Type:        "SizeLRU",
				Capacity:    25,
				SizeInBytes: 20971520, // 20MB
			},
			DB: config.DBConfig{
				FilePath:          "OutportBlocks",
				Type:              "LvlDBSerial",
				BatchDelaySeconds: 2,
				MaxBatchSize:      100,
				MaxOpenFiles:      10,
			},
		},
	}

	return cfg
}

func TestBlocksPool_FullPersisterMode(t *testing.T) {
	dbMode := common.FullPersisterDBMode

	marshaller := &marshal.GogoProtoMarshalizer{}

	cfg := getDefaultConfig()
	cfg.OutportBlocksStorage.DB.FilePath = t.TempDir()

	blocksStorer, err := factory.CreateStorer(cfg, dbMode)
	defer func() {
		_ = blocksStorer.Destroy()
	}()
	require.Nil(t, err)

	blocksPool, err := process.NewBlocksPool(blocksStorer, marshaller, cfg.DataPool.MaxDelta, cfg.DataPool.PruningWindow, cfg.DataPool.FirstCommitableBlock)
	require.Nil(t, err)

	shardID := uint32(2)
	maxDelta := cfg.DataPool.MaxDelta
	pruningWindow := cfg.DataPool.PruningWindow

	outportData := &outport.OutportBlock{ShardID: shardID}
	metaOutportData := &outport.OutportBlock{ShardID: core.MetachainShardId}

	// should fail after maxDelta attempts
	for i := uint64(0); i < maxDelta; i++ {
		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, outportData, i)
		require.Nil(t, err)
	}

	oldHash := []byte("hash_" + fmt.Sprintf("%d", maxDelta-1))
	hash := []byte("hash_" + fmt.Sprintf("%d", maxDelta))

	err = blocksPool.PutBlock(hash, outportData, maxDelta)
	require.Nil(t, err)

	// should find in storer
	retData, err := blocksStorer.Get(hash)
	require.Nil(t, err)
	require.NotNil(t, retData)

	hash = []byte("hash_" + fmt.Sprintf("%d", maxDelta+1))
	err = blocksPool.PutBlock(hash, outportData, maxDelta+1)
	require.Equal(t, process.ErrFailedToPutBlockDataToPool, err)

	if pruningWindow <= maxDelta {
		require.Fail(t, "prunning window should be bigger than delta")
	}

	for i := maxDelta + 1; i < pruningWindow; i++ {
		checkpoint := &data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				core.MetachainShardId: i,
			},
		}
		blocksPool.UpdateMetaState(checkpoint)

		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, metaOutportData, i)
		require.Nil(t, err)

		err = blocksPool.PutBlock(hash, outportData, i)
		require.Nil(t, err)
	}

	// should still find in storer
	retData, err = blocksStorer.Get(oldHash)
	require.Nil(t, err)
	require.NotNil(t, retData)

	for i := pruningWindow; i < pruningWindow*2; i++ {
		checkpoint := &data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				core.MetachainShardId: i,
			},
		}
		blocksPool.UpdateMetaState(checkpoint)

		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, metaOutportData, i)
		require.Nil(t, err)

		err = blocksPool.PutBlock(hash, outportData, i)
		require.Nil(t, err)
	}

	// should still find in storer
	retData, err = blocksStorer.Get(oldHash)
	require.Nil(t, err)
	require.NotNil(t, retData)

	for i := pruningWindow * 2; i < pruningWindow*3; i++ {
		checkpoint := &data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				core.MetachainShardId: i,
			},
		}
		blocksPool.UpdateMetaState(checkpoint)

		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, metaOutportData, i)
		require.Nil(t, err)

		err = blocksPool.PutBlock(hash, outportData, i)
		require.Nil(t, err)
	}

	checkpoint := &data.BlockCheckpoint{
		LastRounds: map[uint32]uint64{
			core.MetachainShardId: pruningWindow * 3,
		},
	}
	blocksPool.UpdateMetaState(checkpoint)

	// should not find in storer anymore
	retData, err = blocksStorer.Get(oldHash)
	require.Nil(t, retData)
	require.Error(t, err)
}

func TestBlocksPool_OptimizedPersisterMode(t *testing.T) {
	dbMode := common.OptimizedPersisterDBMode

	marshaller := &marshal.GogoProtoMarshalizer{}

	cfg := getDefaultConfig()
	cfg.OutportBlocksStorage.DB.FilePath = t.TempDir()
	cfg.DataPool.PruningWindow = 25
	cfg.OutportBlocksStorage.Cache.Capacity = 20

	blocksStorer, err := factory.CreateStorer(cfg, dbMode)
	defer func() {
		_ = blocksStorer.Destroy()
	}()
	require.Nil(t, err)

	blocksPool, err := process.NewBlocksPool(blocksStorer, marshaller, cfg.DataPool.MaxDelta, cfg.DataPool.PruningWindow, cfg.DataPool.FirstCommitableBlock)
	require.Nil(t, err)

	shardID := uint32(2)
	maxDelta := cfg.DataPool.MaxDelta
	pruningWindow := cfg.DataPool.PruningWindow

	outportData := &outport.OutportBlock{ShardID: shardID}
	metaOutportData := &outport.OutportBlock{ShardID: core.MetachainShardId}

	// should fail after maxDelta attempts
	for i := uint64(0); i < maxDelta; i++ {
		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, outportData, i)
		require.Nil(t, err)
	}

	hash := []byte("hash_" + fmt.Sprintf("%d", maxDelta))

	err = blocksPool.PutBlock(hash, outportData, maxDelta)
	require.Nil(t, err)

	// should find in storer
	retData, err := blocksStorer.Get(hash)
	require.Nil(t, err)
	require.NotNil(t, retData)

	hash = []byte("hash_" + fmt.Sprintf("%d", maxDelta+1))
	err = blocksPool.PutBlock(hash, outportData, maxDelta+1)
	require.Equal(t, process.ErrFailedToPutBlockDataToPool, err)

	if pruningWindow <= maxDelta {
		require.Fail(t, "prunning window should be bigger than delta")
	}

	for i := maxDelta + 1; i < pruningWindow; i++ {
		checkpoint := &data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				core.MetachainShardId: i,
			},
		}
		blocksPool.UpdateMetaState(checkpoint)

		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, metaOutportData, i)
		require.Nil(t, err)

		err = blocksPool.PutBlock(hash, outportData, i)
		require.Nil(t, err)
	}

	for i := pruningWindow; i < pruningWindow*2; i++ {
		checkpoint := &data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				core.MetachainShardId: i,
			},
		}
		blocksPool.UpdateMetaState(checkpoint)

		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, metaOutportData, i)
		require.Nil(t, err)

		err = blocksPool.PutBlock(hash, outportData, i)
		require.Nil(t, err)
	}

	checkpoint := &data.BlockCheckpoint{
		LastRounds: map[uint32]uint64{
			core.MetachainShardId: pruningWindow * 2,
		},
	}
	blocksPool.UpdateMetaState(checkpoint)

	recentHash := []byte("hash_" + fmt.Sprintf("%d", pruningWindow*2-1))

	// should find in cache
	retData, err = blocksPool.Get(recentHash)
	require.Nil(t, err)
	require.NotNil(t, retData)

	// close will dump cached data to persister
	err = blocksPool.Close()
	require.Nil(t, err)

	// should find in newly created storer
	blocksStorer2, err := factory.CreateStorer(cfg, dbMode)
	defer func() {
		_ = blocksStorer2.Destroy()
	}()
	require.Nil(t, err)

	retData, err = blocksStorer2.Get(recentHash)
	require.Nil(t, err)
	require.NotNil(t, retData)
}

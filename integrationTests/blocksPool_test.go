package integrationtests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/common"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/config"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/factory"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
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

	marshaller := &process.ProtoMarshaller{}

	cfg := getDefaultConfig()
	cfg.OutportBlocksStorage.DB.FilePath = t.TempDir()

	blocksStorer, err := factory.CreateStorer(cfg, dbMode)
	defer func() {
		_ = blocksStorer.Destroy()
	}()
	require.Nil(t, err)

	argsBlocksPool := process.BlocksPoolArgs{
		Storer:               blocksStorer,
		Marshaller:           marshaller,
		MaxDelta:             cfg.DataPool.MaxDelta,
		CleanupInterval:      cfg.DataPool.PruningWindow,
		FirstCommitableBlock: cfg.DataPool.FirstCommitableBlock,
	}
	blocksPool, err := process.NewBlocksPool(argsBlocksPool)
	require.Nil(t, err)

	shardID := uint32(2)
	maxDelta := cfg.DataPool.MaxDelta
	pruningWindow := cfg.DataPool.PruningWindow

	// should fail after maxDelta attempts
	for i := uint64(0); i < maxDelta; i++ {

		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, []byte("data1"), i, shardID)
		require.Nil(t, err)
	}

	checkpoint := &data.BlockCheckpoint{
		LastNonces: map[uint32]uint64{
			core.MetachainShardId: 1,
			shardID:               1,
		},
	}
	err = blocksPool.UpdateMetaState(checkpoint)
	require.Nil(t, err)

	oldHash := []byte("hash_" + fmt.Sprintf("%d", maxDelta-1))
	hash := []byte("hash_" + fmt.Sprintf("%d", maxDelta))

	err = blocksPool.PutBlock(hash, []byte("data1"), maxDelta, shardID)
	require.Nil(t, err)

	// should find in storer
	retData, err := blocksStorer.Get(hash)
	require.Nil(t, err)
	require.NotNil(t, retData)

	hash = []byte("hash_" + fmt.Sprintf("%d", maxDelta+1))
	err = blocksPool.PutBlock(hash, []byte("data1"), maxDelta+1, shardID)
	require.True(t, errors.Is(err, process.ErrFailedToPutBlockDataToPool))

	require.True(t, pruningWindow > maxDelta, "prunning window should be bigger than delta")

	for i := maxDelta + 1; i < pruningWindow; i++ {
		checkpoint := &data.BlockCheckpoint{
			LastNonces: map[uint32]uint64{
				core.MetachainShardId: i,
				shardID:               i,
			},
		}
		err = blocksPool.UpdateMetaState(checkpoint)
		require.Nil(t, err)

		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, []byte("data1"), i, core.MetachainShardId)
		require.Nil(t, err)

		err = blocksPool.PutBlock(hash, []byte("data1"), i, shardID)
		require.Nil(t, err)
	}

	// should still find in storer
	retData, err = blocksStorer.Get(oldHash)
	require.Nil(t, err)
	require.NotNil(t, retData)

	for i := pruningWindow; i < pruningWindow*2; i++ {
		checkpoint := &data.BlockCheckpoint{
			LastNonces: map[uint32]uint64{
				core.MetachainShardId: i,
				shardID:               i,
			},
		}
		err = blocksPool.UpdateMetaState(checkpoint)
		require.Nil(t, err)

		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, []byte("data1"), i, core.MetachainShardId)
		require.Nil(t, err)

		err = blocksPool.PutBlock(hash, []byte("data1"), i, shardID)
		require.Nil(t, err)
	}

	// should still find in storer
	retData, err = blocksStorer.Get(oldHash)
	require.Nil(t, err)
	require.NotNil(t, retData)

	for i := pruningWindow * 2; i < pruningWindow*3; i++ {
		checkpoint := &data.BlockCheckpoint{
			LastNonces: map[uint32]uint64{
				core.MetachainShardId: i,
				shardID:               i,
			},
		}
		err = blocksPool.UpdateMetaState(checkpoint)
		require.Nil(t, err)

		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, []byte("data1"), i, core.MetachainShardId)
		require.Nil(t, err)

		err = blocksPool.PutBlock(hash, []byte("data1"), i, shardID)
		require.Nil(t, err)
	}

	checkpoint = &data.BlockCheckpoint{
		LastNonces: map[uint32]uint64{
			core.MetachainShardId: pruningWindow * 3,
		},
	}
	err = blocksPool.UpdateMetaState(checkpoint)
	require.Nil(t, err)

	// should not find in storer anymore
	retData, err = blocksStorer.Get(oldHash)
	require.Nil(t, retData)
	require.Error(t, err)
}

func TestBlocksPool_OptimizedPersisterMode(t *testing.T) {
	dbMode := common.OptimizedPersisterDBMode

	marshaller := &process.ProtoMarshaller{}

	cfg := getDefaultConfig()
	cfg.OutportBlocksStorage.DB.FilePath = t.TempDir()
	cfg.DataPool.PruningWindow = 25
	cfg.OutportBlocksStorage.Cache.Capacity = 20

	blocksStorer, err := factory.CreateStorer(cfg, dbMode)
	defer func() {
		_ = blocksStorer.Destroy()
	}()
	require.Nil(t, err)

	argsBlocksPool := process.BlocksPoolArgs{
		Storer:               blocksStorer,
		Marshaller:           marshaller,
		MaxDelta:             cfg.DataPool.MaxDelta,
		CleanupInterval:      cfg.DataPool.PruningWindow,
		FirstCommitableBlock: cfg.DataPool.FirstCommitableBlock,
	}
	blocksPool, err := process.NewBlocksPool(argsBlocksPool)
	require.Nil(t, err)

	shardID := uint32(2)
	maxDelta := cfg.DataPool.MaxDelta
	pruningWindow := cfg.DataPool.PruningWindow

	// should fail after maxDelta attempts
	for i := uint64(0); i < maxDelta; i++ {
		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, []byte("data1"), i, shardID)
		require.Nil(t, err)
	}

	checkpoint := &data.BlockCheckpoint{
		LastNonces: map[uint32]uint64{
			core.MetachainShardId: 1,
			shardID:               1,
		},
	}
	err = blocksPool.UpdateMetaState(checkpoint)
	require.Nil(t, err)

	hash := []byte("hash_" + fmt.Sprintf("%d", maxDelta))

	err = blocksPool.PutBlock(hash, []byte("data1"), maxDelta, shardID)
	require.Nil(t, err)

	// should find in storer
	retData, err := blocksStorer.Get(hash)
	require.Nil(t, err)
	require.NotNil(t, retData)

	hash = []byte("hash_" + fmt.Sprintf("%d", maxDelta+1))
	err = blocksPool.PutBlock(hash, []byte("data1"), maxDelta+1, shardID)
	require.True(t, errors.Is(err, process.ErrFailedToPutBlockDataToPool))

	require.True(t, pruningWindow > maxDelta, "prunning window should be bigger than delta")

	for i := maxDelta + 1; i < pruningWindow; i++ {
		checkpoint := &data.BlockCheckpoint{
			LastNonces: map[uint32]uint64{
				core.MetachainShardId: i,
				shardID:               i,
			},
		}
		err = blocksPool.UpdateMetaState(checkpoint)
		require.Nil(t, err)

		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, []byte("data1"), i, core.MetachainShardId)
		require.Nil(t, err)

		err = blocksPool.PutBlock(hash, []byte("data1"), i, shardID)
		require.Nil(t, err)
	}

	for i := pruningWindow; i < pruningWindow*2; i++ {
		checkpoint := &data.BlockCheckpoint{
			LastNonces: map[uint32]uint64{
				core.MetachainShardId: i,
				shardID:               i,
			},
		}
		err = blocksPool.UpdateMetaState(checkpoint)
		require.Nil(t, err)

		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, []byte("data1"), i, core.MetachainShardId)
		require.Nil(t, err)

		err = blocksPool.PutBlock(hash, []byte("data1"), i, shardID)
		require.Nil(t, err)
	}

	checkpoint = &data.BlockCheckpoint{
		LastNonces: map[uint32]uint64{
			core.MetachainShardId: pruningWindow * 2,
			shardID:               pruningWindow + 2,
		},
	}
	err = blocksPool.UpdateMetaState(checkpoint)
	require.Nil(t, err)

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

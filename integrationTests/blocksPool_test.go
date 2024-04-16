package integrationtests

import (
	"fmt"
	"testing"

	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/factory"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process/dataPool"
	"github.com/stretchr/testify/require"
)

func getDefaultConfig() config.Config {
	maxDelta := uint64(10)
	pruningWindow := uint64(25)

	cfg := config.Config{
		DataPool: config.DataPoolConfig{
			NumberOfShards:      3,
			MaxDelta:            maxDelta,
			PruningWindow:       pruningWindow,
			NumPersistersToKeep: 2,
		},
		OutportBlocksStorage: config.StorageConfig{
			Cache: config.CacheConfig{
				Name:        "OutportBlocksStorage",
				Type:        "SizeLRU",
				Capacity:    100,
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

func TestBlocksPool(t *testing.T) {
	dbMode := "full-persister"

	marshaller := &marshal.GogoProtoMarshalizer{}

	cfg := getDefaultConfig()

	blocksStorer, err := factory.CreateStorer(cfg, dbMode)
	defer blocksStorer.Destroy()
	require.Nil(t, err)

	blocksPool, err := dataPool.NewBlocksPool(blocksStorer, marshaller, cfg.DataPool.NumberOfShards, cfg.DataPool.MaxDelta, cfg.DataPool.PruningWindow)
	require.Nil(t, err)

	shardID := uint32(2)
	maxDelta := cfg.DataPool.MaxDelta
	pruningWindow := cfg.DataPool.PruningWindow

	// should fail after maxDelta attempts
	for i := uint64(0); i < maxDelta; i++ {
		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, []byte("value1"), i, shardID)
		require.Nil(t, err)
	}

	oldHash := []byte("hash_" + fmt.Sprintf("%d", maxDelta-1))
	hash := []byte("hash_" + fmt.Sprintf("%d", maxDelta))

	err = blocksPool.PutBlock(hash, []byte("value1"), maxDelta, shardID)
	require.Nil(t, err)

	// should find in storer
	retData, err := blocksStorer.Get(hash)
	require.Nil(t, err)
	require.NotNil(t, retData)

	hash = []byte("hash_" + fmt.Sprintf("%d", maxDelta+1))
	err = blocksPool.PutBlock(hash, []byte("value1"), maxDelta+1, shardID)
	require.Equal(t, dataPool.ErrFailedToPutBlockDataToPool, err)

	if pruningWindow <= maxDelta {
		require.Fail(t, "prunning window should be bigger than delta")
	}

	for i := maxDelta; i < pruningWindow; i++ {
		blocksPool.UpdateMetaState(i)

		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, []byte("value1"), i, shardID)
		require.Nil(t, err)
	}

	// should still find in storer
	retData, err = blocksStorer.Get(oldHash)
	require.Nil(t, err)
	require.NotNil(t, retData)

	for i := pruningWindow; i < pruningWindow*2; i++ {
		blocksPool.UpdateMetaState(i)

		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, []byte("value1"), i, shardID)
		require.Nil(t, err)
	}

	// should still find in storer
	retData, err = blocksStorer.Get(oldHash)
	require.Nil(t, err)
	require.NotNil(t, retData)

	for i := pruningWindow * 2; i < pruningWindow*3; i++ {
		blocksPool.UpdateMetaState(i)

		hash := []byte("hash_" + fmt.Sprintf("%d", i))

		err = blocksPool.PutBlock(hash, []byte("value1"), i, shardID)
		require.Nil(t, err)
	}

	blocksPool.UpdateMetaState(pruningWindow * 3)

	// should not find in storer anymore
	retData, err = blocksStorer.Get(oldHash)
	require.Nil(t, retData)
	require.Error(t, err)
}

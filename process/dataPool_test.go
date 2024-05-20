package process_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/testscommon"
)

func createDefaultBlocksPoolArgs() process.DataPoolArgs {
	return process.DataPoolArgs{
		Storer:          &testscommon.PruningStorerMock{},
		Marshaller:      gogoProtoMarshaller,
		MaxDelta:        10,
		CleanupInterval: 100,
		FirstCommitableBlocks: map[uint32]uint64{
			core.MetachainShardId: 0,
			0:                     0,
			1:                     0,
			2:                     0,
		},
	}
}

func TestNewDataPool(t *testing.T) {
	t.Parallel()

	t.Run("nil pruning storer", func(t *testing.T) {
		t.Parallel()

		args := createDefaultBlocksPoolArgs()
		args.Storer = nil
		bp, err := process.NewDataPool(args)
		require.Nil(t, bp)
		require.Equal(t, process.ErrNilPruningStorer, err)
	})

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		args := createDefaultBlocksPoolArgs()
		args.Marshaller = nil
		bp, err := process.NewDataPool(args)
		require.Nil(t, bp)
		require.Equal(t, process.ErrNilMarshaller, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		args := createDefaultBlocksPoolArgs()
		bp, err := process.NewDataPool(args)
		require.Nil(t, err)
		require.False(t, bp.IsInterfaceNil())
	})
}

func TestDataPool_GetBlock(t *testing.T) {
	t.Parallel()

	t.Run("failed to get data from storer", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")

		args := createDefaultBlocksPoolArgs()
		args.Storer = &testscommon.PruningStorerMock{
			GetCalled: func(key []byte) ([]byte, error) {
				return nil, expectedErr
			},
		}
		bp, _ := process.NewDataPool(args)

		ret, err := bp.Get([]byte("hash1"))
		require.Nil(t, ret)
		require.Equal(t, expectedErr, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		outportBlock := &outport.OutportBlock{
			ShardID:                1,
			NotarizedHeadersHashes: []string{"hash1", "hash2"},
			NumberOfShards:         3,
		}
		outportBlockBytes, _ := gogoProtoMarshaller.Marshal(outportBlock)

		args := createDefaultBlocksPoolArgs()
		args.Storer = &testscommon.PruningStorerMock{
			GetCalled: func(key []byte) ([]byte, error) {
				return outportBlockBytes, nil
			},
		}
		bp, _ := process.NewDataPool(args)

		ret, err := bp.Get([]byte("hash1"))
		require.Nil(t, err)
		require.Equal(t, outportBlockBytes, ret)
	})
}

func TestDataPool_UpdateMetaState(t *testing.T) {
	t.Parallel()

	cleanupInterval := uint64(100)

	t.Run("should not set checkpoint if index is not commitable", func(t *testing.T) {
		t.Parallel()

		firstCommitableBlock := uint64(10)
		firstCommitableBlocks := map[uint32]uint64{
			core.MetachainShardId: firstCommitableBlock,
		}

		args := createDefaultBlocksPoolArgs()
		args.CleanupInterval = cleanupInterval
		args.FirstCommitableBlocks = firstCommitableBlocks
		args.Marshaller = protoMarshaller
		args.Storer = &testscommon.PruningStorerMock{
			PutCalled: func(key, data []byte) error {
				assert.Fail(t, "should have not been called")

				return nil
			},
			PruneCalled: func(index uint64) error {
				assert.Fail(t, "should have not been called")

				return nil
			},
		}
		bp, _ := process.NewDataPool(args)

		checkpoint := &data.BlockCheckpoint{
			LastNonces: map[uint32]uint64{
				core.MetachainShardId: firstCommitableBlock - 1,
			},
		}

		err := bp.UpdateMetaState(checkpoint)
		require.Nil(t, err)
	})

	t.Run("should set checkpoint if index if commitable", func(t *testing.T) {
		t.Parallel()

		firstCommitableBlock := uint64(10)
		firstCommitableBlocks := map[uint32]uint64{
			core.MetachainShardId: firstCommitableBlock,
		}

		putCalled := false
		args := createDefaultBlocksPoolArgs()
		args.CleanupInterval = cleanupInterval
		args.FirstCommitableBlocks = firstCommitableBlocks
		args.Marshaller = protoMarshaller
		args.Storer = &testscommon.PruningStorerMock{
			PutCalled: func(key, data []byte) error {
				putCalled = true

				return nil
			},
			PruneCalled: func(index uint64) error {
				assert.Fail(t, "should have not been called")

				return nil
			},
		}
		bp, _ := process.NewDataPool(args)

		checkpoint := &data.BlockCheckpoint{
			LastNonces: map[uint32]uint64{
				core.MetachainShardId: firstCommitableBlock,
			},
		}

		err := bp.UpdateMetaState(checkpoint)
		require.Nil(t, err)

		require.True(t, putCalled)
	})

	t.Run("should not trigger prune if not cleanup interval", func(t *testing.T) {
		t.Parallel()

		args := createDefaultBlocksPoolArgs()
		args.Marshaller = protoMarshaller
		args.Storer = &testscommon.PruningStorerMock{
			PruneCalled: func(index uint64) error {
				assert.Fail(t, "should have not been called")

				return nil
			},
		}
		bp, _ := process.NewDataPool(args)

		checkpoint := &data.BlockCheckpoint{
			LastNonces: map[uint32]uint64{
				core.MetachainShardId: 2,
			},
		}

		err := bp.UpdateMetaState(checkpoint)
		require.Nil(t, err)
	})

	t.Run("should trigger prune if cleanup interval", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		args := createDefaultBlocksPoolArgs()
		args.CleanupInterval = cleanupInterval
		args.Marshaller = protoMarshaller
		args.Storer = &testscommon.PruningStorerMock{
			PruneCalled: func(index uint64) error {
				wasCalled = true

				return nil
			},
		}
		bp, _ := process.NewDataPool(args)

		checkpoint := &data.BlockCheckpoint{
			LastNonces: map[uint32]uint64{
				core.MetachainShardId: 100,
			},
		}

		err := bp.UpdateMetaState(checkpoint)
		require.Nil(t, err)

		require.True(t, wasCalled)
	})
}

func TestDataPool_PutBlock(t *testing.T) {
	t.Parallel()

	shardID := uint32(2)

	t.Run("should work on init, index 0", func(t *testing.T) {
		t.Parallel()

		maxDelta := uint64(10)

		wasCalled := false
		args := createDefaultBlocksPoolArgs()
		args.MaxDelta = maxDelta
		args.Storer = &testscommon.PruningStorerMock{
			PutCalled: func(key, data []byte) error {
				wasCalled = true

				return nil
			},
		}
		bp, _ := process.NewDataPool(args)

		startIndex := uint64(0)
		err := bp.PutBlock([]byte("hash1"), []byte("data1"), startIndex, shardID)
		require.Nil(t, err)

		require.True(t, wasCalled)
	})

	t.Run("should work on init with any index if not previous checkpoint", func(t *testing.T) {
		t.Parallel()

		maxDelta := uint64(10)

		wasCalled := false

		args := createDefaultBlocksPoolArgs()
		args.MaxDelta = maxDelta
		args.Marshaller = protoMarshaller
		args.Storer = &testscommon.PruningStorerMock{
			GetCalled: func(key []byte) ([]byte, error) {
				if string(key) == process.MetaCheckpointKey {
					wasCalled = true
					return nil, fmt.Errorf("no checkpoint key found")
				}

				return []byte{}, nil
			},
		}
		bp, _ := process.NewDataPool(args)

		startIndex := uint64(123)
		err := bp.PutBlock([]byte("hash1"), []byte("data1"), startIndex, shardID)
		require.Nil(t, err)

		require.True(t, wasCalled)
	})

	t.Run("should work succesively from init if there is previous checkpoint", func(t *testing.T) {
		t.Parallel()

		maxDelta := uint64(10)

		shardID := uint32(1)
		startIndex := uint64(123)

		lastCheckpointData, err := protoMarshaller.Marshal(&data.BlockCheckpoint{
			LastNonces: map[uint32]uint64{
				shardID:               startIndex,
				core.MetachainShardId: startIndex - 2,
			},
		})
		require.Nil(t, err)

		putCalled := false

		args := createDefaultBlocksPoolArgs()
		args.MaxDelta = maxDelta
		args.Storer = &testscommon.PruningStorerMock{
			GetCalled: func(key []byte) ([]byte, error) {
				if string(key) == process.MetaCheckpointKey {
					return lastCheckpointData, nil
				}

				return []byte{}, nil
			},
			PutCalled: func(key, data []byte) error {
				putCalled = true

				return nil
			},
		}
		bp, _ := process.NewDataPool(args)

		err = bp.PutBlock([]byte("hash1"), []byte("data1"), startIndex+1, shardID)
		require.Nil(t, err)

		require.True(t, putCalled)
	})

	t.Run("should fail if max delta is reached", func(t *testing.T) {
		t.Parallel()

		maxDelta := uint64(10)

		wasCalled := false

		checkpoint := &data.BlockCheckpoint{
			LastNonces: map[uint32]uint64{
				core.MetachainShardId: 2,
				shardID:               2,
			},
		}
		checkpointBytes, _ := protoMarshaller.Marshal(checkpoint)

		args := createDefaultBlocksPoolArgs()
		args.MaxDelta = maxDelta
		args.Marshaller = protoMarshaller
		args.Storer = &testscommon.PruningStorerMock{
			PutCalled: func(key, data []byte) error {
				wasCalled = true

				return nil
			},
			GetCalled: func(key []byte) ([]byte, error) {
				return checkpointBytes, nil
			},
		}
		bp, _ := process.NewDataPool(args)

		startIndex := uint64(2)
		err := bp.PutBlock([]byte("hash1"), []byte("data1"), startIndex, shardID)
		require.Nil(t, err)

		require.True(t, wasCalled)

		err = bp.UpdateMetaState(checkpoint)
		require.Nil(t, err)

		err = bp.PutBlock([]byte("hash2"), []byte("data1"), startIndex, shardID)
		require.Nil(t, err)

		for i := uint64(1); i < maxDelta; i++ {
			err = bp.PutBlock([]byte("hash2"), []byte("data1"), startIndex+i, shardID)
			require.Nil(t, err)
		}

		err = bp.PutBlock([]byte("hash3"), []byte("data1"), startIndex+maxDelta, shardID)
		require.Error(t, err)
	})
}

func TestDataPool_Close(t *testing.T) {
	t.Parallel()

	wasCalled := false
	args := createDefaultBlocksPoolArgs()
	args.Storer = &testscommon.PruningStorerMock{
		CloseCalled: func() error {
			wasCalled = true

			return nil
		},
	}
	bp, _ := process.NewDataPool(args)

	err := bp.Close()
	require.Nil(t, err)

	require.True(t, wasCalled)
}

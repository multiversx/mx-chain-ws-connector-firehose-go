package process_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/testscommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBlocksPool(t *testing.T) {
	t.Parallel()

	t.Run("nil pruning storer", func(t *testing.T) {
		t.Parallel()

		bp, err := process.NewBlocksPool(
			nil,
			&testscommon.MarshallerMock{},
			10,
			100,
			0,
		)
		require.Nil(t, bp)
		require.Equal(t, process.ErrNilPruningStorer, err)
	})

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		bp, err := process.NewBlocksPool(
			&testscommon.PruningStorerStub{},
			nil,
			10,
			100,
			0,
		)
		require.Nil(t, bp)
		require.Equal(t, process.ErrNilMarshaller, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		bp, err := process.NewBlocksPool(
			&testscommon.PruningStorerStub{},
			&testscommon.MarshallerMock{},
			10,
			100,
			0,
		)
		require.Nil(t, err)
		require.False(t, bp.IsInterfaceNil())
	})
}

func TestBlocksPool_GetBlock(t *testing.T) {
	t.Parallel()

	t.Run("failed to get data from storer", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		bp, _ := process.NewBlocksPool(
			&testscommon.PruningStorerStub{
				GetCalled: func(key []byte) ([]byte, error) {
					return nil, expectedErr
				},
			},
			protoMarshaller,
			10,
			100,
			0,
		)

		ret, err := bp.GetBlock([]byte("hash1"))
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

		bp, _ := process.NewBlocksPool(
			&testscommon.PruningStorerStub{
				GetCalled: func(key []byte) ([]byte, error) {
					return outportBlockBytes, nil
				},
			},
			gogoProtoMarshaller,
			10,
			100,
			0,
		)

		ret, err := bp.GetBlock([]byte("hash1"))
		require.Nil(t, err)
		require.Equal(t, outportBlock, ret)
	})
}

func TestBlocksPool_UpdateMetaState(t *testing.T) {
	t.Parallel()

	cleanupInterval := uint64(100)

	t.Run("should not set checkpoint if index is not commitable", func(t *testing.T) {
		t.Parallel()

		firstCommitableBlock := uint64(10)

		bp, _ := process.NewBlocksPool(
			&testscommon.PruningStorerStub{
				PutCalled: func(key, data []byte) error {
					assert.Fail(t, "should have not been called")

					return nil
				},
				PruneCalled: func(index uint64) error {
					assert.Fail(t, "should have not been called")

					return nil
				},
			},
			protoMarshaller,
			100,
			cleanupInterval,
			firstCommitableBlock,
		)

		checkpoint := &data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				core.MetachainShardId: firstCommitableBlock - 1,
			},
		}

		bp.UpdateMetaState(checkpoint)
	})

	t.Run("should set checkpoint if index if commitable", func(t *testing.T) {
		t.Parallel()

		firstCommitableBlock := uint64(10)

		putCalled := false
		bp, _ := process.NewBlocksPool(
			&testscommon.PruningStorerStub{
				PutCalled: func(key, data []byte) error {
					putCalled = true

					return nil
				},
				PruneCalled: func(index uint64) error {
					assert.Fail(t, "should have not been called")

					return nil
				},
			},
			gogoProtoMarshaller,
			100,
			cleanupInterval,
			firstCommitableBlock,
		)

		checkpoint := &data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				core.MetachainShardId: firstCommitableBlock,
			},
		}

		bp.UpdateMetaState(checkpoint)

		require.True(t, putCalled)
	})

	t.Run("should not trigger prune if not cleanup interval", func(t *testing.T) {
		t.Parallel()

		bp, _ := process.NewBlocksPool(
			&testscommon.PruningStorerStub{
				PruneCalled: func(index uint64) error {
					assert.Fail(t, "should have not been called")

					return nil
				},
			},
			protoMarshaller,
			100,
			cleanupInterval,
			0,
		)

		checkpoint := &data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				core.MetachainShardId: 2,
			},
		}

		bp.UpdateMetaState(checkpoint)
	})

	t.Run("should trigger prune if cleanup interval", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		bp, _ := process.NewBlocksPool(
			&testscommon.PruningStorerStub{
				PruneCalled: func(index uint64) error {
					wasCalled = true

					return nil
				},
			},
			protoMarshaller,
			100,
			cleanupInterval,
			0,
		)

		checkpoint := &data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				core.MetachainShardId: 100,
			},
		}

		bp.UpdateMetaState(checkpoint)

		require.True(t, wasCalled)
	})
}

func TestBlocksPool_PutBlock(t *testing.T) {
	t.Parallel()

	t.Run("should work on init, index 0", func(t *testing.T) {
		t.Parallel()

		maxDelta := uint64(10)

		wasCalled := false
		bp, _ := process.NewBlocksPool(
			&testscommon.PruningStorerStub{
				PutCalled: func(key, data []byte) error {
					wasCalled = true

					return nil
				},
			},
			gogoProtoMarshaller,
			maxDelta,
			100,
			0,
		)

		startIndex := uint64(0)
		err := bp.PutBlock([]byte("hash1"), &outport.OutportBlock{}, startIndex)
		require.Nil(t, err)

		require.True(t, wasCalled)
	})

	t.Run("should work on init with any index if not previous checkpoint", func(t *testing.T) {
		t.Parallel()

		maxDelta := uint64(10)

		wasCalled := false
		bp, _ := process.NewBlocksPool(
			&testscommon.PruningStorerStub{
				GetCalled: func(key []byte) ([]byte, error) {
					if string(key) == process.MetaCheckpointKey {
						wasCalled = true
						return nil, fmt.Errorf("no checkpoint key found")
					}

					return []byte{}, nil
				},
			},
			gogoProtoMarshaller,
			maxDelta,
			100,
			0,
		)

		startIndex := uint64(123)
		err := bp.PutBlock([]byte("hash1"), &outport.OutportBlock{}, startIndex)
		require.Nil(t, err)

		require.True(t, wasCalled)
	})

	t.Run("should work succesively from init if there is previous checkpoint", func(t *testing.T) {
		t.Parallel()

		maxDelta := uint64(10)

		shardID := uint32(1)
		startIndex := uint64(123)

		lastCheckpointData, err := gogoProtoMarshaller.Marshal(&data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				shardID:               startIndex,
				core.MetachainShardId: startIndex - 2,
			},
		})
		require.Nil(t, err)

		putCalled := false
		bp, _ := process.NewBlocksPool(
			&testscommon.PruningStorerStub{
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
			},
			gogoProtoMarshaller,
			maxDelta,
			100,
			0,
		)

		err = bp.PutBlock([]byte("hash1"), &outport.OutportBlock{ShardID: shardID}, startIndex+1)
		require.Nil(t, err)

		require.True(t, putCalled)
	})

	t.Run("should fail if no succesive index", func(t *testing.T) {
		t.Parallel()

		maxDelta := uint64(10)

		shardID := uint32(1)
		startIndex := uint64(123)

		lastCheckpointData, err := gogoProtoMarshaller.Marshal(&data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				shardID:               startIndex,
				core.MetachainShardId: startIndex - 2,
			},
		})
		require.Nil(t, err)

		putCalled := false
		bp, _ := process.NewBlocksPool(
			&testscommon.PruningStorerStub{
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
			},
			gogoProtoMarshaller,
			maxDelta,
			100,
			0,
		)

		err = bp.PutBlock([]byte("hash1"), &outport.OutportBlock{ShardID: shardID}, startIndex)
		require.True(t, errors.Is(err, process.ErrFailedToPutBlockDataToPool))

		require.False(t, putCalled)
	})

	t.Run("should fail if max delta is reached", func(t *testing.T) {
		t.Parallel()

		maxDelta := uint64(10)

		wasCalled := false
		bp, _ := process.NewBlocksPool(
			&testscommon.PruningStorerStub{
				PutCalled: func(key, data []byte) error {
					wasCalled = true

					return nil
				},
			},
			gogoProtoMarshaller,
			maxDelta,
			100,
			0,
		)

		outportData := &outport.OutportBlock{ShardID: uint32(1)}

		startIndex := uint64(2)
		err := bp.PutBlock([]byte("hash1"), outportData, startIndex)
		require.Nil(t, err)

		require.True(t, wasCalled)

		checkpoint := &data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				core.MetachainShardId: 2,
			},
		}

		bp.UpdateMetaState(checkpoint)
		require.Nil(t, err)

		metaOutportData := &outport.OutportBlock{ShardID: core.MetachainShardId}
		err = bp.PutBlock([]byte("hash2"), metaOutportData, startIndex)
		require.Nil(t, err)

		for i := uint64(1); i <= maxDelta; i++ {
			err = bp.PutBlock([]byte("hash2"), outportData, startIndex+i)
			require.Nil(t, err)
		}

		err = bp.PutBlock([]byte("hash3"), outportData, startIndex+maxDelta+1)
		require.Error(t, err)
	})
}

func TestBlocksPool_Close(t *testing.T) {
	t.Parallel()

	wasCalled := false
	bp, _ := process.NewBlocksPool(
		&testscommon.PruningStorerStub{
			CloseCalled: func() error {
				wasCalled = true

				return nil
			},
		},
		protoMarshaller,
		10,
		100,
		0,
	)

	err := bp.Close()
	require.Nil(t, err)

	require.True(t, wasCalled)
}

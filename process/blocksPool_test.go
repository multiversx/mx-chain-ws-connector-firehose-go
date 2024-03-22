package process_test

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/testscommon"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testStruct struct {
	s    string
	a, b int
}

func TestNewBlocksPool(t *testing.T) {
	t.Parallel()

	t.Run("nil pruning storer", func(t *testing.T) {
		t.Parallel()

		bp, err := process.NewBlocksPool(
			nil,
			&testscommon.MarshallerMock{},
			3,
			10,
			100,
		)
		require.Nil(t, bp)
		require.Equal(t, process.ErrNilPruningStorer, err)
	})

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		bp, err := process.NewBlocksPool(
			&testscommon.PruningStorerStub{},
			nil,
			3,
			10,
			100,
		)
		require.Nil(t, bp)
		require.Equal(t, process.ErrNilMarshaller, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		bp, err := process.NewBlocksPool(
			&testscommon.PruningStorerStub{},
			&testscommon.MarshallerMock{},
			3,
			10,
			100,
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
			3,
			10,
			100,
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
		outportBlockBytes, _ := protoMarshaller.Marshal(outportBlock)

		bp, _ := process.NewBlocksPool(
			&testscommon.PruningStorerStub{
				GetCalled: func(key []byte) ([]byte, error) {
					return outportBlockBytes, nil
				},
			},
			protoMarshaller,
			3,
			10,
			100,
		)

		ret, err := bp.GetBlock([]byte("hash1"))
		require.Nil(t, err)
		require.Equal(t, outportBlock, ret)
	})
}

func TestBlocksPool_UpdateMetaState(t *testing.T) {
	t.Parallel()

	t.Run("should not trigger prune if not cleanup interval", func(t *testing.T) {
		t.Parallel()

		cleanupInterval := uint64(100)

		bp, _ := process.NewBlocksPool(
			&testscommon.PruningStorerStub{
				PruneCalled: func(index uint64) error {
					assert.Fail(t, "should have not been called")

					return nil
				},
			},
			protoMarshaller,
			3,
			100,
			cleanupInterval,
		)

		bp.UpdateMetaState(2)
	})

	t.Run("should trigger prune if cleanup interval", func(t *testing.T) {
		t.Parallel()

		cleanupInterval := uint64(100)

		wasCalled := false
		bp, _ := process.NewBlocksPool(
			&testscommon.PruningStorerStub{
				PruneCalled: func(index uint64) error {
					wasCalled = true

					return nil
				},
			},
			protoMarshaller,
			3,
			100,
			cleanupInterval,
		)

		bp.UpdateMetaState(100)

		require.True(t, wasCalled)
	})
}

func TestBlocksPool_PutBlock(t *testing.T) {
	t.Parallel()

	t.Run("first put, should put directly", func(t *testing.T) {
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
			protoMarshaller,
			3,
			maxDelta,
			100,
		)

		err := bp.PutBlock([]byte("hash1"), &outport.OutportBlock{}, 2)
		require.Nil(t, err)

		require.True(t, wasCalled)

		bp.UpdateMetaState(2)
		require.Nil(t, err)

		err = bp.PutBlock([]byte("hash2"), &outport.OutportBlock{}, 2+maxDelta+1)
		require.Nil(t, err)

		err = bp.PutBlock([]byte("hash2"), &outport.OutportBlock{}, 2+maxDelta+2)
		require.Error(t, err)
	})
}

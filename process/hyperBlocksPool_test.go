package process_test

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/testscommon"
	"github.com/stretchr/testify/require"
)

func TestNewHyperBlocksPool(t *testing.T) {
	t.Parallel()

	t.Run("nil data pool", func(t *testing.T) {
		t.Parallel()

		hbp, err := process.NewHyperOutportBlocksPool(
			nil,
			&testscommon.MarshallerMock{},
		)
		require.Nil(t, hbp)
		require.Equal(t, process.ErrNilDataPool, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		hbp, err := process.NewHyperOutportBlocksPool(
			&testscommon.BlocksPoolStub{},
			nil,
		)
		require.Nil(t, hbp)
		require.Equal(t, process.ErrNilMarshaller, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		hbp, err := process.NewHyperOutportBlocksPool(
			&testscommon.BlocksPoolStub{},
			&testscommon.MarshallerMock{},
		)
		require.Nil(t, err)
		require.False(t, hbp.IsInterfaceNil())
	})
}

func TestHyperOutportBlocksPool_PutMetaBlock(t *testing.T) {
	t.Parallel()

	t.Run("should fail if no succesive index", func(t *testing.T) {
		t.Parallel()

		hbp, err := process.NewHyperOutportBlocksPool(
			&testscommon.BlocksPoolStub{},
			&testscommon.MarshallerMock{},
		)
		require.Nil(t, err)

		currentIndex := uint64(10)

		outportBlock := &hyperOutportBlocks.MetaOutportBlock{
			ShardID: 1,
			BlockData: &hyperOutportBlocks.MetaBlockData{
				Header: &hyperOutportBlocks.MetaHeader{
					Round: currentIndex,
				},
			},
			HighestFinalBlockNonce: currentIndex + 1,
		}

		err = hbp.PutMetaBlock([]byte("hash1"), outportBlock)
		require.True(t, errors.Is(err, process.ErrFailedToPutBlockDataToPool))
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		hbp, err := process.NewHyperOutportBlocksPool(
			&testscommon.BlocksPoolStub{},
			&testscommon.MarshallerMock{},
		)
		require.Nil(t, err)

		currentIndex := uint64(10)

		outportBlock := &hyperOutportBlocks.MetaOutportBlock{
			ShardID: 1,
			BlockData: &hyperOutportBlocks.MetaBlockData{
				Header: &hyperOutportBlocks.MetaHeader{
					Round: currentIndex,
				},
			},
			HighestFinalBlockNonce: currentIndex - 1,
		}

		err = hbp.PutMetaBlock([]byte("hash1"), outportBlock)
		require.Nil(t, err)
	})
}

func TestHyperOutportBlocksPool_PutShardBlock(t *testing.T) {
	t.Parallel()

	t.Run("should fail if no succesive index", func(t *testing.T) {
		t.Parallel()

		hbp, err := process.NewHyperOutportBlocksPool(
			&testscommon.BlocksPoolStub{},
			&testscommon.MarshallerMock{},
		)
		require.Nil(t, err)

		currentIndex := uint64(10)

		outportBlock := &hyperOutportBlocks.ShardOutportBlock{
			ShardID: 1,
			BlockData: &hyperOutportBlocks.BlockData{
				Header: &hyperOutportBlocks.Header{
					Round: currentIndex,
				},
			},
			HighestFinalBlockNonce: currentIndex + 1,
		}

		err = hbp.PutShardBlock([]byte("hash1"), outportBlock)
		require.True(t, errors.Is(err, process.ErrFailedToPutBlockDataToPool))
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		hbp, err := process.NewHyperOutportBlocksPool(
			&testscommon.BlocksPoolStub{},
			&testscommon.MarshallerMock{},
		)
		require.Nil(t, err)

		currentIndex := uint64(10)

		outportBlock := &hyperOutportBlocks.ShardOutportBlock{
			ShardID: 1,
			BlockData: &hyperOutportBlocks.BlockData{
				Header: &hyperOutportBlocks.Header{
					Round: currentIndex,
				},
			},
			HighestFinalBlockNonce: currentIndex - 1,
		}

		err = hbp.PutShardBlock([]byte("hash1"), outportBlock)
		require.Nil(t, err)
	})
}

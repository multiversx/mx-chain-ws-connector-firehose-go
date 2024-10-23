package process_test

import (
	"errors"
	"testing"

	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/testscommon"
	"github.com/stretchr/testify/require"
)

func TestBlocksPool(t *testing.T) {
	t.Parallel()

	t.Run("nil data pool", func(t *testing.T) {
		t.Parallel()

		hbp, err := process.NewBlocksPool(
			nil,
			&testscommon.MarshallerMock{},
		)
		require.Nil(t, hbp)
		require.Equal(t, process.ErrNilDataPool, err)
	})

	t.Run("nil marshaller", func(t *testing.T) {
		t.Parallel()

		hbp, err := process.NewBlocksPool(
			&testscommon.BlocksPoolMock{},
			nil,
		)
		require.Nil(t, hbp)
		require.Equal(t, process.ErrNilMarshaller, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		hbp, err := process.NewBlocksPool(
			&testscommon.BlocksPoolMock{},
			&testscommon.MarshallerMock{},
		)
		require.Nil(t, err)
		require.False(t, hbp.IsInterfaceNil())
	})
}

func TestBlocksPool_PutBlock(t *testing.T) {
	t.Parallel()

	t.Run("should fail if no succesive index", func(t *testing.T) {
		t.Parallel()

		hbp, err := process.NewBlocksPool(
			&testscommon.BlocksPoolMock{},
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

		err = hbp.PutBlock([]byte("hash1"), outportBlock)
		require.True(t, errors.Is(err, process.ErrFailedToPutBlockDataToPool))
	})

	t.Run("should work for meta block", func(t *testing.T) {
		t.Parallel()

		hbp, err := process.NewBlocksPool(
			&testscommon.BlocksPoolMock{},
			&testscommon.MarshallerMock{},
		)
		require.Nil(t, err)

		currentIndex := uint64(10)

		outportBlock := &hyperOutportBlocks.MetaOutportBlock{
			ShardID: 1,
			BlockData: &hyperOutportBlocks.MetaBlockData{
				Header: &hyperOutportBlocks.MetaHeader{
					Nonce: currentIndex,
				},
			},
			HighestFinalBlockNonce: currentIndex - 1,
		}

		err = hbp.PutBlock([]byte("hash1"), outportBlock)
		require.Nil(t, err)
	})

	t.Run("should work for shard block", func(t *testing.T) {
		t.Parallel()

		hbp, err := process.NewBlocksPool(
			&testscommon.BlocksPoolMock{},
			&testscommon.MarshallerMock{},
		)
		require.Nil(t, err)

		currentIndex := uint64(10)

		outportBlock := &hyperOutportBlocks.ShardOutportBlock{
			ShardID: 1,
			BlockData: &hyperOutportBlocks.BlockData{
				Header: &hyperOutportBlocks.Header{
					Nonce: currentIndex,
				},
			},
			HighestFinalBlockNonce: currentIndex - 1,
		}

		err = hbp.PutBlock([]byte("hash1"), outportBlock)
		require.Nil(t, err)
	})
}

func TestBlocksPool_GetMetaBlock(t *testing.T) {
	t.Parallel()

	marshaller := &testscommon.MarshallerMock{}

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
	outportBlockBytes, _ := marshaller.Marshal(outportBlock)

	hbp, err := process.NewBlocksPool(
		&testscommon.BlocksPoolMock{
			GetCalled: func(hash []byte) ([]byte, error) {
				return outportBlockBytes, nil
			},
		},
		marshaller,
	)
	require.Nil(t, err)

	ret, err := hbp.GetMetaBlock([]byte("hash1"))
	require.Nil(t, err)
	require.Equal(t, outportBlock, ret)
}

func TestBlocksPool_GetShardBlock(t *testing.T) {
	t.Parallel()

	marshaller := &testscommon.MarshallerMock{}

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
	outportBlockBytes, _ := marshaller.Marshal(outportBlock)

	hbp, err := process.NewBlocksPool(
		&testscommon.BlocksPoolMock{
			GetCalled: func(hash []byte) ([]byte, error) {
				return outportBlockBytes, nil
			},
		},
		marshaller,
	)
	require.Nil(t, err)

	ret, err := hbp.GetShardBlock([]byte("hash1"))
	require.Nil(t, err)
	require.Equal(t, outportBlock, ret)
}

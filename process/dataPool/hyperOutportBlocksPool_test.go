package dataPool_test

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process/dataPool"
	"github.com/multiversx/mx-chain-ws-connector-template-go/testscommon"
	"github.com/stretchr/testify/require"
)

func TestNewHyperOutportBlocksPool(t *testing.T) {
	t.Parallel()

	t.Run("nil marshaller, should fail", func(t *testing.T) {
		t.Parallel()

		hbp, err := dataPool.NewHyperOutportBlocksPool(&testscommon.BlocksPoolStub{}, nil)
		require.Nil(t, hbp)
		require.Equal(t, process.ErrNilMarshaller, err)
	})

	t.Run("nil data pool, should fail", func(t *testing.T) {
		t.Parallel()

		hbp, err := dataPool.NewHyperOutportBlocksPool(nil, &testscommon.MarshallerMock{})
		require.Nil(t, hbp)
		require.Equal(t, dataPool.ErrNilDataPool, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		hbp, err := dataPool.NewHyperOutportBlocksPool(&testscommon.BlocksPoolStub{}, &testscommon.MarshallerMock{})
		require.Nil(t, err)
		require.False(t, hbp.IsInterfaceNil())
	})
}

func TestPutBlock(t *testing.T) {
	t.Parallel()

	t.Run("nil meta outport block, should fail", func(t *testing.T) {
		t.Parallel()

		hbp, err := dataPool.NewHyperOutportBlocksPool(&testscommon.BlocksPoolStub{}, &testscommon.MarshallerMock{})
		require.Nil(t, err)

		err = hbp.PutBlock([]byte("hash1"), &data.HyperOutportBlock{MetaOutportBlock: nil}, 2)
		require.Equal(t, dataPool.ErrNilMetaOutportBlock, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		wasCalled := false
		hbp, err := dataPool.NewHyperOutportBlocksPool(&testscommon.BlocksPoolStub{
			PutBlockCalled: func(hash, data []byte, round uint64, shardID uint32) error {
				wasCalled = true
				return nil
			},
		}, &testscommon.MarshallerMock{})
		require.Nil(t, err)

		hyperOutportBlock := &data.HyperOutportBlock{
			MetaOutportBlock: &outport.OutportBlock{
				ShardID: 2,
				BlockData: &outport.BlockData{
					ShardID:    2,
					HeaderHash: []byte("headerHash1"),
				},
			},
		}

		err = hbp.PutBlock([]byte("hash1"), hyperOutportBlock, 2)
		require.Nil(t, err)

		require.True(t, wasCalled)
	})
}

func TestGetBlock(t *testing.T) {
	t.Parallel()

	marshaller := &testscommon.MarshallerMock{}

	hyperOutportBlock := &data.HyperOutportBlock{
		MetaOutportBlock: &outport.OutportBlock{
			ShardID: 2,
			BlockData: &outport.BlockData{
				ShardID:    2,
				HeaderHash: []byte("headerHash1"),
			},
		},
	}
	hyperOutportBlockBytes, _ := marshaller.Marshal(hyperOutportBlock)

	hash1 := []byte("hash1")

	hbp, err := dataPool.NewHyperOutportBlocksPool(&testscommon.BlocksPoolStub{
		GetBlockCalled: func(hash []byte) ([]byte, error) {
			require.Equal(t, hash1, hash)
			return hyperOutportBlockBytes, nil
		},
	}, marshaller)
	require.Nil(t, err)

	ret, err := hbp.GetBlock(hash1)
	require.Nil(t, err)
	require.Equal(t, hyperOutportBlock, ret)
}

func TestNewHyperOutportBlocksPool_Close(t *testing.T) {
	t.Parallel()

	wasCalled := false
	hbp, err := dataPool.NewHyperOutportBlocksPool(&testscommon.BlocksPoolStub{
		CloseCalled: func() error {
			wasCalled = true
			return nil
		},
	}, &testscommon.MarshallerMock{})
	require.Nil(t, err)

	require.Nil(t, hbp.Close())

	require.True(t, wasCalled)
}

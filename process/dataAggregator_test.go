package process_test

import (
	"encoding/hex"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/stretchr/testify/require"

	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/testscommon"
)

func TestNewDataAggregator(t *testing.T) {
	t.Parallel()

	t.Run("nil blocks pool", func(t *testing.T) {
		t.Parallel()

		da, err := process.NewDataAggregator(nil, process.NewOutportBlockConverter())
		require.Nil(t, da)
		require.Equal(t, process.ErrNilBlocksPool, err)
	})

	t.Run("nil outport block converter", func(t *testing.T) {
		t.Parallel()

		da, err := process.NewDataAggregator(&testscommon.BlocksPoolStub{}, nil)
		require.Nil(t, da)
		require.Equal(t, process.ErrNilOutportBlocksConverter, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		da, err := process.NewDataAggregator(&testscommon.BlocksPoolStub{}, process.NewOutportBlockConverter())
		require.Nil(t, err)
		require.False(t, da.IsInterfaceNil())
	})
}

func TestDataAggregator_ProcessHyperBlock(t *testing.T) {
	t.Parallel()

	t.Run("invalid outport block provided", func(t *testing.T) {
		t.Parallel()

		blocksPoolStub := &testscommon.BlocksPoolStub{}

		da, err := process.NewDataAggregator(blocksPoolStub, process.NewOutportBlockConverter())
		require.Nil(t, err)

		shardOutportBlock := createOutportBlock()

		hyperOutportBlock, err := da.ProcessHyperBlock(shardOutportBlock)
		require.Nil(t, hyperOutportBlock)
		require.Equal(t, process.ErrInvalidOutportBlock, err)
	})

	headerHash := []byte("headerHash1")

	shardOutportBlock := createOutportBlock()
	shardOutportBlock.BlockData.HeaderHash = headerHash

	blocksPoolStub := &testscommon.BlocksPoolStub{
		GetBlockCalled: func(hash []byte) (*outport.OutportBlock, error) {

			return shardOutportBlock, nil
		},
	}

	converter := process.NewOutportBlockConverter()
	expectedResult, err := converter.HandleShardOutportBlock(shardOutportBlock)
	require.NoError(t, err)

	da, err := process.NewDataAggregator(blocksPoolStub, process.NewOutportBlockConverter())
	require.Nil(t, err)

	outportBlock := createMetaOutportBlock()
	outportBlock.NotarizedHeadersHashes = []string{hex.EncodeToString(headerHash)}

	hyperOutportBlock, err := da.ProcessHyperBlock(outportBlock)
	require.Nil(t, err)
	require.Equal(t, expectedResult, hyperOutportBlock.NotarizedHeadersOutportData[0].OutportBlock)
}

package process_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/testscommon"
)

func TestNewDataAggregator(t *testing.T) {
	t.Parallel()

	t.Run("nil blocks pool", func(t *testing.T) {
		t.Parallel()

		da, err := process.NewDataAggregator(nil)
		require.Nil(t, da)
		require.Equal(t, process.ErrNilHyperBlocksPool, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		da, err := process.NewDataAggregator(&testscommon.HyperBlocksPoolStub{})
		require.Nil(t, err)
		require.False(t, da.IsInterfaceNil())
	})
}

func TestDataAggregator_ProcessHyperBlock(t *testing.T) {
	t.Parallel()

	t.Run("invalid outport block provided", func(t *testing.T) {
		t.Parallel()

		da, err := process.NewDataAggregator(&testscommon.HyperBlocksPoolStub{})
		require.Nil(t, err)

		invalidOutportBlock := createMetaOutportBlock()
		invalidOutportBlock.ShardID = 2

		hyperOutportBlock, err := da.ProcessHyperBlock(invalidOutportBlock)
		require.Nil(t, hyperOutportBlock)
		require.Equal(t, process.ErrInvalidOutportBlock, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		headerHash := []byte("headerHash1")

		shardOutportBlock := createShardOutportBlock()
		shardOutportBlock.BlockData.HeaderHash = headerHash

		blocksPoolStub := &testscommon.HyperBlocksPoolStub{
			GetShardBlockCalled: func(hash []byte) (*hyperOutportBlocks.ShardOutportBlock, error) {
				return shardOutportBlock, nil
			},
		}

		da, err := process.NewDataAggregator(blocksPoolStub)
		require.Nil(t, err)

		outportBlock := createMetaOutportBlock()
		outportBlock.NotarizedHeadersHashes = []string{hex.EncodeToString(headerHash)}

		expectedHyperOutportBlock := &hyperOutportBlocks.HyperOutportBlock{}
		expectedHyperOutportBlock.MetaOutportBlock = outportBlock
		expectedHyperOutportBlock.NotarizedHeadersOutportData = append(expectedHyperOutportBlock.NotarizedHeadersOutportData, &hyperOutportBlocks.NotarizedHeaderOutportData{
			NotarizedHeaderHash: hex.EncodeToString(headerHash),
			OutportBlock:        shardOutportBlock,
		})

		hyperOutportBlock, err := da.ProcessHyperBlock(outportBlock)
		require.Nil(t, err)
		require.Equal(t, expectedHyperOutportBlock, hyperOutportBlock)
	})
}

package process_test

import (
	"encoding/hex"
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/testscommon"
	"github.com/stretchr/testify/require"
)

func TestNewDataAggregator(t *testing.T) {
	t.Parallel()

	t.Run("nil blocks pool", func(t *testing.T) {
		t.Parallel()

		da, err := process.NewDataAggregator(nil)
		require.Nil(t, da)
		require.Equal(t, process.ErrNilBlocksPool, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		da, err := process.NewDataAggregator(&testscommon.BlocksPoolStub{})
		require.Nil(t, err)
		require.False(t, da.IsInterfaceNil())
	})
}

func TestDataAggregator_ProcessHyperBlock(t *testing.T) {
	t.Parallel()

	headerHash := []byte("headerHash1")

	shardOutportBlock := createOutportBlock()
	shardOutportBlock.BlockData.HeaderHash = headerHash

	blocksPoolStub := &testscommon.BlocksPoolStub{
		GetBlockCalled: func(hash []byte) (*outport.OutportBlock, error) {

			return shardOutportBlock, nil
		},
	}

	da, err := process.NewDataAggregator(blocksPoolStub)
	require.Nil(t, err)

	outportBlock := createMetaOutportBlock()
	outportBlock.NotarizedHeadersHashes = []string{hex.EncodeToString(headerHash)}

	hyperOutportBlock, err := da.ProcessHyperBlock(outportBlock)
	require.Nil(t, err)
	require.Equal(t, shardOutportBlock, hyperOutportBlock.NotarizedHeadersOutportData[0].OutportBlock)
}

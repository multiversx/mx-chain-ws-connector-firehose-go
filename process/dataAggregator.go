package process

import (
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
)

type dataAggregator struct {
	blocksPool BlocksPool
}

// NewDataAggregator will create a new data aggregator instance
func NewDataAggregator(
	blocksPool BlocksPool,
) (*dataAggregator, error) {
	if check.IfNil(blocksPool) {
		return nil, ErrNilBlocksPool
	}

	return &dataAggregator{
		blocksPool: blocksPool,
	}, nil
}

// ProcessHyperBlock will process meta outport block. It will try to fetch and aggregate
// notarized shards data
func (da *dataAggregator) ProcessHyperBlock(outportBlock *outport.OutportBlock) (*data.HyperOutportBlock, error) {
	if outportBlock.ShardID != core.MetachainShardId {
		return nil, ErrInvalidOutportBlock
	}

	hyperOutportBlock := &data.HyperOutportBlock{}
	hyperOutportBlock.MetaOutportBlock = outportBlock

	notarizedShardOutportBlocks := make([]*data.NotarizedHeaderOutportData, 0)
	for _, notarizedHash := range outportBlock.NotarizedHeadersHashes {
		hash, err := hex.DecodeString(notarizedHash)
		if err != nil {
			return nil, err
		}

		outportBlockShard, err := da.blocksPool.GetBlock(hash)
		if err != nil {
			return nil, err
		}

		notarizedShardOutportBlock := &data.NotarizedHeaderOutportData{
			NotarizedHeaderHash: notarizedHash,
			OutportBlock:        outportBlockShard,
		}

		notarizedShardOutportBlocks = append(notarizedShardOutportBlocks, notarizedShardOutportBlock)
	}

	hyperOutportBlock.NotarizedHeadersOutportData = notarizedShardOutportBlocks

	return hyperOutportBlock, nil
}

// IsInterfaceNil returns true if there is no value under interface
func (da *dataAggregator) IsInterfaceNil() bool {
	return da == nil
}

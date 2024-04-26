package process

import (
	"encoding/hex"
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"

	"github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

type dataAggregator struct {
	blocksPool HyperBlocksPool
}

// NewDataAggregator will create a new data aggregator instance
func NewDataAggregator(
	blocksPool HyperBlocksPool,
) (*dataAggregator, error) {
	if check.IfNil(blocksPool) {
		return nil, ErrNilHyperBlocksPool
	}

	return &dataAggregator{
		blocksPool: blocksPool,
	}, nil
}

// ProcessHyperBlock will process meta outport block. It will try to fetch and aggregate
// notarized shards data
func (da *dataAggregator) ProcessHyperBlock(metaOutportBlock *hyperOutportBlocks.MetaOutportBlock) (*data.HyperOutportBlock, error) {
	if metaOutportBlock.ShardID != core.MetachainShardId {
		return nil, ErrInvalidOutportBlock
	}

	hyperOutportBlock := &data.HyperOutportBlock{}
	hyperOutportBlock.MetaOutportBlock = metaOutportBlock

	notarizedShardOutportBlocks := make([]*data.NotarizedHeaderOutportData, 0)

	log.Info("dataAggregator: notarizedHashes", "block hash", metaOutportBlock.BlockData.HeaderHash,
		"num notarizedHashes", len(metaOutportBlock.NotarizedHeadersHashes))

	for _, notarizedHash := range metaOutportBlock.NotarizedHeadersHashes {
		hash, err := hex.DecodeString(notarizedHash)
		if err != nil {
			return nil, fmt.Errorf("failed to decode notarized hash string: %w", err)
		}

		shardOutportBlock, err := da.blocksPool.GetShardBlock(hash)
		if err != nil {
			return nil, err
		}

		log.Info("dataAggregator: get block", "hash", hash)

		notarizedShardOutportBlock := &data.NotarizedHeaderOutportData{
			NotarizedHeaderHash: notarizedHash,
			OutportBlock:        shardOutportBlock,
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

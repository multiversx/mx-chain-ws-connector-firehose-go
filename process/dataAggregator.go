package process

import (
	"encoding/hex"
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"

	data "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
)

type dataAggregator struct {
	blocksPool DataPool
	converter  OutportBlockConverter
}

// NewDataAggregator will create a new data aggregator instance
func NewDataAggregator(
	blocksPool DataPool,
	converter OutportBlockConverter,
) (*dataAggregator, error) {
	if check.IfNil(blocksPool) {
		return nil, ErrNilBlocksPool
	}
	if check.IfNil(converter) {
		return nil, ErrNilOutportBlocksConverter
	}

	return &dataAggregator{
		blocksPool: blocksPool,
		converter:  converter,
	}, nil
}

// ProcessHyperBlock will process meta outport block. It will try to fetch and aggregate
// notarized shards data
func (da *dataAggregator) ProcessHyperBlock(outportBlock *outport.OutportBlock) (*data.HyperOutportBlock, error) {
	if outportBlock.ShardID != core.MetachainShardId {
		return nil, ErrInvalidOutportBlock
	}

	hyperOutportBlock := &data.HyperOutportBlock{}
	metaOutportBlock, err := da.converter.HandleMetaOutportBlock(outportBlock)
	if err != nil {
		return nil, fmt.Errorf("couldn't process meta outportBlock: %w", err)
	}
	hyperOutportBlock.MetaOutportBlock = metaOutportBlock

	notarizedShardOutportBlocks := make([]*data.NotarizedHeaderOutportData, 0)

	log.Info("dataAggregator: notarizedHashes", "block hash", outportBlock.BlockData.HeaderHash,
		"num notarizedHashes", len(outportBlock.NotarizedHeadersHashes))

	for _, notarizedHash := range outportBlock.NotarizedHeadersHashes {
		hash, err := hex.DecodeString(notarizedHash)
		if err != nil {
			return nil, err
		}

		outportBlockShard, err := da.blocksPool.GetBlock(hash)
		if err != nil {
			return nil, err
		}

		log.Info("dataAggregator: get block", "hash", hash)

		shardBlock, err := da.converter.HandleShardOutportBlock(outportBlockShard)
		if err != nil {
			return nil, fmt.Errorf("couldn't process shard outportBlock: %w", err)
		}

		notarizedShardOutportBlock := &data.NotarizedHeaderOutportData{
			NotarizedHeaderHash: notarizedHash,
			OutportBlock:        shardBlock,
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

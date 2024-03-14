package process

import (
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
)

type dataAggregator struct {
	blocksPool BlocksPool
}

func NewDataAggregator(
	blocksPool BlocksPool,
) (*dataAggregator, error) {
	if check.IfNil(blocksPool) {
		return nil, errNilBlocksPool
	}

	return &dataAggregator{
		blocksPool: blocksPool,
	}, nil
}

func (da *dataAggregator) ProcessHyperBlock(outportBlock *outport.OutportBlock) (*data.HyperOutportBlock, error) {
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

	return hyperOutportBlock, nil
}

// IsInterfaceNil returns true if there is no value under interface
func (da *dataAggregator) IsInterfaceNil() bool {
	return da == nil
}

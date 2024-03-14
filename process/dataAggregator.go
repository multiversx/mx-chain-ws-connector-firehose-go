package process

import (
	"encoding/hex"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	coreData "github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
)

type dataAggregator struct {
	blockCreator BlockContainerHandler
	blocksPool   BlocksPool
	marshaller   marshal.Marshalizer
}

func NewDataAggregator(
	blockCreator BlockContainerHandler,
	blocksPool BlocksPool,
	marshaller marshal.Marshalizer,
) (*dataAggregator, error) {
	if check.IfNil(blockCreator) {
		return nil, errNilBlockCreator
	}
	if check.IfNil(blocksPool) {
		return nil, errNilBlocksPool
	}
	if check.IfNil(marshaller) {
		return nil, errNilMarshaller
	}

	return &dataAggregator{
		blockCreator: blockCreator,
		blocksPool:   blocksPool,
		marshaller:   marshaller,
	}, nil
}

func (da *dataAggregator) ProcessHyperBlock(outportBlock *outport.OutportBlock) (coreData.HeaderHandler, []byte, error) {
	blockCreator, err := da.blockCreator.Get(core.HeaderType(outportBlock.BlockData.HeaderType))
	if err != nil {
		return nil, nil, err
	}

	header, err := block.GetHeaderFromBytes(da.marshaller, blockCreator, outportBlock.BlockData.HeaderBytes)
	if err != nil {
		return nil, nil, err
	}

	hyperOutportBlock := &data.HyperOutportBlock{}
	hyperOutportBlock.MetaOutportBlock = outportBlock

	notarizedShardOutportBlocks := make([]*data.NotarizedHeaderOutportData, 0)
	for _, notarizedHash := range outportBlock.NotarizedHeadersHashes {
		hash, err := hex.DecodeString(notarizedHash)
		if err != nil {
			return nil, nil, err
		}

		outportBlockShard, err := da.blocksPool.GetBlock(hash)
		if err != nil {
			return nil, nil, err
		}

		notarizedShardOutportBlock := &data.NotarizedHeaderOutportData{
			NotarizedHeaderHash: notarizedHash,
			OutportBlock:        outportBlockShard,
		}

		notarizedShardOutportBlocks = append(notarizedShardOutportBlocks, notarizedShardOutportBlock)
	}

	marshalledData, err := da.marshaller.Marshal(hyperOutportBlock)
	if err != nil {
		return nil, nil, err
	}

	return header, marshalledData, nil
}

// IsInterfaceNil returns true if there is no value under interface
func (da *dataAggregator) IsInterfaceNil() bool {
	return da == nil
}

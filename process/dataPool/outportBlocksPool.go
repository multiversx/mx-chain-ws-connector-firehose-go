package dataPool

import (
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
)

type outportBlocksPool struct {
	marshaller marshal.Marshalizer
	dataPool   process.DataPool
}

// NewOutportBlocksPool will create a new outport block pool component
func NewOutportBlocksPool(
	dataPool process.DataPool,
	marshaller marshal.Marshalizer,
) (*outportBlocksPool, error) {
	if check.IfNil(dataPool) {
		return nil, ErrNilDataPool
	}
	if check.IfNil(marshaller) {
		return nil, process.ErrNilMarshaller
	}

	return &outportBlocksPool{
		dataPool:   dataPool,
		marshaller: marshaller,
	}, nil
}

func (bp *outportBlocksPool) UpdateMetaState(round uint64) {
	bp.dataPool.UpdateMetaState(round)
}

// PutBlock will put the provided outport block data to the pool
func (bp *outportBlocksPool) PutBlock(hash []byte, outportBlock *outport.OutportBlock, currentRound uint64) error {
	outportBlockBytes, err := bp.marshaller.Marshal(outportBlock)
	if err != nil {
		return err
	}

	return bp.dataPool.PutBlock(hash, outportBlockBytes, currentRound, outportBlock.ShardID)
}

// GetBlock will return outport block data from the pool
func (bp *outportBlocksPool) GetBlock(hash []byte) (*outport.OutportBlock, error) {
	marshalledData, err := bp.dataPool.GetBlock(hash)
	if err != nil {
		return nil, err
	}

	outportBlock := &outport.OutportBlock{}
	err = bp.marshaller.Unmarshal(outportBlock, marshalledData)
	if err != nil {
		return nil, err
	}

	return outportBlock, nil
}

// Close will trigger close on blocks pool component
func (bp *outportBlocksPool) Close() error {
	return nil
}

// IsInterfaceNil returns nil if there is no value under the interface
func (bp *outportBlocksPool) IsInterfaceNil() bool {
	return bp == nil
}

package dataPool

import (
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
)

type hyperOutportBlocksPool struct {
	marshaller marshal.Marshalizer
	dataPool   process.DataPool
}

// NewHyperOutportBlocksPool will create a new hyper outport blocks pool component
func NewHyperOutportBlocksPool(
	dataPool process.DataPool,
	marshaller marshal.Marshalizer,
) (*hyperOutportBlocksPool, error) {
	if check.IfNil(dataPool) {
		return nil, ErrNilDataPool
	}
	if check.IfNil(marshaller) {
		return nil, process.ErrNilMarshaller
	}

	return &hyperOutportBlocksPool{
		dataPool:   dataPool,
		marshaller: marshaller,
	}, nil
}

// UpdateMetaState will trigger data pool update state
func (bp *hyperOutportBlocksPool) UpdateMetaState(round uint64) {
	bp.dataPool.UpdateMetaState(round)
}

// PutBlock will put the provided hyper outport block data to the pool
func (bp *hyperOutportBlocksPool) PutBlock(hash []byte, outportBlock *data.HyperOutportBlock, currentRound uint64) error {
	if outportBlock.MetaOutportBlock == nil {
		return ErrNilMetaOutportBlock
	}

	shardID := outportBlock.MetaOutportBlock.ShardID

	outportBlockBytes, err := bp.marshaller.Marshal(outportBlock)
	if err != nil {
		return err
	}

	return bp.dataPool.PutBlock(hash, outportBlockBytes, currentRound, shardID)
}

// GetBlock will return the hyper outport block data from the pool
func (bp *hyperOutportBlocksPool) GetBlock(hash []byte) (*data.HyperOutportBlock, error) {
	marshalledData, err := bp.dataPool.Get(hash)
	if err != nil {
		return nil, err
	}

	outportBlock := &data.HyperOutportBlock{}
	err = bp.marshaller.Unmarshal(outportBlock, marshalledData)
	if err != nil {
		return nil, err
	}

	return outportBlock, nil
}

// Close will trigger close on data pool component
func (bp *hyperOutportBlocksPool) Close() error {
	return bp.dataPool.Close()
}

// IsInterfaceNil returns nil if there is no value under the interface
func (bp *hyperOutportBlocksPool) IsInterfaceNil() bool {
	return bp == nil
}

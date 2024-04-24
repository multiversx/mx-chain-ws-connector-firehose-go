package process

import (
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

type hyperOutportBlocksPool struct {
	marshaller marshal.Marshalizer
	dataPool   DataPool
}

// NewHyperOutportBlocksPool will create a new instance of hyper outport blocks pool
func NewHyperOutportBlocksPool(
	dataPool DataPool,
	marshaller marshal.Marshalizer,
) (*hyperOutportBlocksPool, error) {
	if check.IfNil(dataPool) {
		return nil, ErrNilDataPool
	}
	if check.IfNil(marshaller) {
		return nil, ErrNilMarshaller
	}

	return &hyperOutportBlocksPool{
		dataPool:   dataPool,
		marshaller: marshaller,
	}, nil
}

// UpdateMetaState will triiger meta state update from base data pool
func (bp *hyperOutportBlocksPool) UpdateMetaState(checkpoint *data.BlockCheckpoint) error {
	return bp.dataPool.UpdateMetaState(checkpoint)
}

// Get will trigger data pool get operation
func (bp *hyperOutportBlocksPool) Get(key []byte) ([]byte, error) {
	return bp.dataPool.Get(key)
}

// PutMetaBlock will put the provided meta outport block data to the pool
func (bp *hyperOutportBlocksPool) PutMetaBlock(hash []byte, outportBlock *hyperOutportBlocks.MetaOutportBlock) error {
	shardID := outportBlock.ShardID
	currentRound := outportBlock.BlockData.Header.GetRound()

	outportBlockBytes, err := bp.marshaller.Marshal(outportBlock)
	if err != nil {
		return err
	}

	return bp.dataPool.PutBlock(hash, outportBlockBytes, currentRound, shardID)
}

// PutShardBlock will put the provided shard outport block data to the pool
func (bp *hyperOutportBlocksPool) PutShardBlock(hash []byte, outportBlock *hyperOutportBlocks.ShardOutportBlock) error {
	shardID := outportBlock.ShardID
	currentRound := outportBlock.BlockData.Header.GetRound()

	outportBlockBytes, err := bp.marshaller.Marshal(outportBlock)
	if err != nil {
		return err
	}

	return bp.dataPool.PutBlock(hash, outportBlockBytes, currentRound, shardID)
}

// GetMetaBlock will return the meta outport block data from the pool
func (bp *hyperOutportBlocksPool) GetMetaBlock(hash []byte) (*hyperOutportBlocks.MetaOutportBlock, error) {
	marshalledData, err := bp.dataPool.Get(hash)
	if err != nil {
		return nil, err
	}

	metaOutportBlock := &hyperOutportBlocks.MetaOutportBlock{}
	err = bp.marshaller.Unmarshal(metaOutportBlock, marshalledData)
	if err != nil {
		return nil, err
	}

	return metaOutportBlock, nil
}

// GetShardBlock will return the shard outport block data from the pool
func (bp *hyperOutportBlocksPool) GetShardBlock(hash []byte) (*hyperOutportBlocks.ShardOutportBlock, error) {
	marshalledData, err := bp.dataPool.Get(hash)
	if err != nil {
		return nil, err
	}

	shardOutportBlock := &hyperOutportBlocks.ShardOutportBlock{}
	err = bp.marshaller.Unmarshal(shardOutportBlock, marshalledData)
	if err != nil {
		return nil, err
	}

	return shardOutportBlock, nil
}

// Close will trigger close on data pool component
func (bp *hyperOutportBlocksPool) Close() error {
	return bp.dataPool.Close()
}

// IsInterfaceNil returns nil if there is no value under the interface
func (bp *hyperOutportBlocksPool) IsInterfaceNil() bool {
	return bp == nil
}

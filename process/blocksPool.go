package process

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
)

type blocksPool struct {
	marshaller marshal.Marshalizer
	dataPool   DataPool
}

// NewBlocksPool will create a new instance of hyper outport blocks pool
func NewBlocksPool(
	dataPool DataPool,
	marshaller marshal.Marshalizer,
) (*blocksPool, error) {
	if check.IfNil(dataPool) {
		return nil, ErrNilDataPool
	}
	if check.IfNil(marshaller) {
		return nil, ErrNilMarshaller
	}

	return &blocksPool{
		dataPool:   dataPool,
		marshaller: marshaller,
	}, nil
}

// UpdateMetaState will triger meta state update from base data pool
func (bp *blocksPool) UpdateMetaState(checkpoint *data.BlockCheckpoint) error {
	return bp.dataPool.UpdateMetaState(checkpoint)
}

// GetLastCheckpoint returns last checkpoint data
func (bp *blocksPool) GetLastCheckpoint() (*data.BlockCheckpoint, error) {
	return bp.dataPool.GetLastCheckpoint()
}

// Get will trigger data pool get operation
func (bp *blocksPool) Get(key []byte) ([]byte, error) {
	data, err := bp.dataPool.Get(key)
	if err != nil {
		return nil, fmt.Errorf("failed to get from data pool: %w", err)
	}

	return data, nil
}

// PutBlock will put the provided outport block data to the pool
func (bp *blocksPool) PutBlock(hash []byte, outportBlock OutportBlockHandler) error {
	shardID := outportBlock.GetShardID()
	currentIndex, err := outportBlock.GetHeaderNonce()
	if err != nil {
		return err
	}

	previousIndex := outportBlock.GetHighestFinalBlockNonce()
	isHigherIndex := currentIndex >= previousIndex
	if !isHigherIndex {
		return fmt.Errorf("%w: new meta index should be higher than previous, previous index %d, new index %d",
			ErrFailedToPutBlockDataToPool, previousIndex, currentIndex)
	}

	// TODO: marshall data only when/if saving to storage
	//  we do not necessarily need to marshall outport block if it is saved only to cache
	outportBlockBytes, err := bp.marshaller.Marshal(outportBlock)
	if err != nil {
		return err
	}

	return bp.dataPool.PutBlock(hash, outportBlockBytes, currentIndex, shardID)
}

// GetMetaBlock will return the meta outport block data from the pool
func (bp *blocksPool) GetMetaBlock(hash []byte) (*hyperOutportBlocks.MetaOutportBlock, error) {
	marshalledData, err := bp.dataPool.Get(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get meta block from pool: %w", err)
	}

	metaOutportBlock := &hyperOutportBlocks.MetaOutportBlock{}
	err = bp.marshaller.Unmarshal(metaOutportBlock, marshalledData)
	if err != nil {
		return nil, err
	}

	return metaOutportBlock, nil
}

// GetShardBlock will return the shard outport block data from the pool
func (bp *blocksPool) GetShardBlock(hash []byte) (*hyperOutportBlocks.ShardOutportBlock, error) {
	marshalledData, err := bp.dataPool.Get(hash)
	if err != nil {
		return nil, fmt.Errorf("failed to get shard block from pool: %w", err)
	}

	shardOutportBlock := &hyperOutportBlocks.ShardOutportBlock{}
	err = bp.marshaller.Unmarshal(shardOutportBlock, marshalledData)
	if err != nil {
		return nil, err
	}

	return shardOutportBlock, nil
}

// Close will trigger close on data pool component
func (bp *blocksPool) Close() error {
	return bp.dataPool.Close()
}

// IsInterfaceNil returns nil if there is no value under the interface
func (bp *blocksPool) IsInterfaceNil() bool {
	return bp == nil
}

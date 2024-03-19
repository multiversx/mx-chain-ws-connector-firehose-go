package process

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-storage-go/types"
)

// TODO: refactor to use storer

type blocksPool struct {
	cacher types.Cacher
}

// NewBlocksPool will create a new blocks pool instance
func NewBlocksPool(cacher types.Cacher) (*blocksPool, error) {
	return &blocksPool{cacher: cacher}, nil
}

// PutBlock will put the provided outport block data to the pool
func (bp *blocksPool) PutBlock(hash []byte, outportBlock *outport.OutportBlock) error {
	_ = bp.cacher.Put(hash, outportBlock, 0)
	return nil
}

// GetBlock will return outport block data from the pool
func (bp *blocksPool) GetBlock(hash []byte) (*outport.OutportBlock, error) {
	data, ok := bp.cacher.Get(hash)
	if !ok {
		// TODO: handle retry/fallback mechanism
		return nil, fmt.Errorf("failed to get data from pool")
	}

	outportBlock, ok := data.(*outport.OutportBlock)
	if !ok {
		return nil, ErrWrongTypeAssertion
	}

	return outportBlock, nil
}

// IsInterfaceNil returns nil if there is no value under the interface
func (bp *blocksPool) IsInterfaceNil() bool {
	return bp == nil
}

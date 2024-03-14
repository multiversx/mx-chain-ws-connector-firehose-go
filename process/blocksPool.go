package process

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-go/storage"
)

type blocksPool struct {
	cacher storage.Cacher
}

func NewBlocksPool(cacher storage.Cacher) (*blocksPool, error) {
	return &blocksPool{cacher: cacher}, nil
}

func (bp *blocksPool) PutBlock(hash []byte, outportBlock *outport.OutportBlock) error {
	_ = bp.cacher.Put(hash, outportBlock, 0)
	return nil
}

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

func (bp *blocksPool) IsInterfaceNil() bool {
	return bp == nil
}

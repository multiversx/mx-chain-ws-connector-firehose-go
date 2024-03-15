package process

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-storage-go/storageUnit"
	"github.com/multiversx/mx-chain-storage-go/types"
)

type blocksData struct {
	blocksCacher types.Cacher
	delta        uint32
}

type blocksPool struct {
	cacher    types.Cacher
	blocksMap map[uint32]*blocksData
}

func NewBlocksPool(cacher types.Cacher) (*blocksPool, error) {
	bp := &blocksPool{
		cacher: cacher,
	}

	numOfShards := uint32(3)

	bp.createBlocksMap(numOfShards)

	return bp, nil
}

func (bp *blocksPool) createBlocksMap(numOfShards uint32) error {
	blocksMap := make(map[uint32]*blocksData)

	for i := uint32(0); i < numOfShards; i++ {
		cacher, err := bp.createCacher()
		if err != nil {
			return err
		}

		blocksMap[i] = &blocksData{
			blocksCacher: cacher,
			delta:        0,
		}
	}

	cacher, err := bp.createCacher()
	if err != nil {
		return err
	}

	blocksMap[core.MetachainShardId] = &blocksData{
		blocksCacher: cacher,
	}

	bp.blocksMap = blocksMap

	return nil
}

func (bp *blocksPool) createCacher() (types.Cacher, error) {
	cacheConfig := storageUnit.CacheConfig{
		Type:        storageUnit.SizeLRUCache,
		SizeInBytes: 209715200, // 200MB
		Capacity:    100,
	}

	cacher, err := storageUnit.NewCache(cacheConfig)
	if err != nil {
		return nil, err
	}

	return cacher, nil
}

func (bp *blocksPool) PutBlock(hash []byte, outportBlock *outport.OutportBlock) error {
	blocksData, ok := bp.blocksMap[outportBlock.ShardID]
	if !ok {
		return fmt.Errorf("did not find shard id %d in blocksMap", outportBlock.ShardID)
	}

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

package process

import (
	"fmt"
	"math"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-storage-go/types"
)

type blocksPool struct {
	storer       types.Storer
	blockCreator BlockContainerHandler
	marshaller   marshal.Marshalizer
	maxDelta     uint64

	roundsMap map[uint32]uint64
	mutMap    sync.RWMutex
}

func NewBlocksPool(
	storer types.Storer,
	blockCreator BlockContainerHandler,
	marshaller marshal.Marshalizer,
) (*blocksPool, error) {
	numberOfShards := uint32(3)

	roundsMap := make(map[uint32]uint64)
	for shardID := uint32(0); shardID < numberOfShards; shardID++ {
		roundsMap[shardID] = 0
	}
	roundsMap[core.MetachainShardId] = 0

	bp := &blocksPool{
		storer:       storer,
		blockCreator: blockCreator,
		roundsMap:    roundsMap,
		marshaller:   marshaller,
		maxDelta:     10,
	}

	return bp, nil
}

func (bp *blocksPool) UpdateMetaRound(round uint64) {
	bp.mutMap.Lock()
	defer bp.mutMap.Unlock()

	bp.roundsMap[core.MetachainShardId] = round
}

func (bp *blocksPool) PutBlock(hash []byte, outportBlock *outport.OutportBlock) error {
	bp.mutMap.Lock()
	defer bp.mutMap.Unlock()

	shardID := outportBlock.ShardID

	round, ok := bp.roundsMap[shardID]
	if !ok {
		return fmt.Errorf("did not find shard id %d in blocksMap", shardID)
	}

	if round == 0 {
		bp.putOutportBlock(hash, outportBlock)
	}

	metaRound := bp.roundsMap[core.MetachainShardId]

	if !bp.shouldPutOutportBlock(round, metaRound) {
		log.Error("failed to put outport block", "hash", hash, "round", round, "metaRound", metaRound)
		return fmt.Errorf("failed to put outport block", "hash", hash, "round", round, "metaRound", metaRound)
	}

	return bp.putOutportBlock(hash, outportBlock)
}

// should be run under mutex
func (bp *blocksPool) shouldPutOutportBlock(round, metaRound uint64) bool {
	diff := float64(int64(round) - int64(metaRound))
	delta := math.Abs(diff)

	if math.Abs(delta) > float64(bp.maxDelta) {
		return false
	}

	return true
}

// should be run under mutex
func (bp *blocksPool) putOutportBlock(hash []byte, outportBlock *outport.OutportBlock) error {
	shardID := outportBlock.ShardID

	blockCreator, err := bp.blockCreator.Get(core.HeaderType(outportBlock.BlockData.HeaderType))
	if err != nil {
		return err
	}

	header, err := block.GetHeaderFromBytes(bp.marshaller, blockCreator, outportBlock.BlockData.HeaderBytes)
	if err != nil {
		return err
	}

	outportBlockBytes, err := bp.marshaller.Marshal(outportBlock)
	if err != nil {
		return err
	}

	err = bp.storer.Put(hash, outportBlockBytes)
	if err != nil {
		return err
	}

	bp.roundsMap[shardID] = header.GetRound()

	return nil
}

func (bp *blocksPool) GetBlock(hash []byte) (*outport.OutportBlock, error) {
	data, err := bp.storer.Get(hash)
	if err != nil {
		// TODO: handle retry/fallback mechanism
		return nil, fmt.Errorf("failed to get data from pool")
	}

	outportBlock := &outport.OutportBlock{}

	err = bp.marshaller.Unmarshal(outportBlock, data)
	if err != nil {
		return nil, err
	}

	return outportBlock, nil
}

func (bp *blocksPool) IsInterfaceNil() bool {
	return bp == nil
}

package process

import (
	"fmt"
	"math"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
)

type blocksPool struct {
	storer          PruningStorer
	marshaller      marshal.Marshalizer
	maxDelta        uint64
	numOfShards     uint32
	cleanupInterval uint64

	roundsMap map[uint32]uint64
	mutMap    sync.RWMutex
}

func NewBlocksPool(
	storer PruningStorer,
	marshaller marshal.Marshalizer,
	numOfShards uint32,
	maxDelta uint64,
	cleanupInterval uint64,
) (*blocksPool, error) {
	if check.IfNil(storer) {
		return nil, ErrNilPruningStorer
	}
	if check.IfNil(marshaller) {
		return nil, ErrNilMarshaller
	}

	bp := &blocksPool{
		storer:          storer,
		marshaller:      marshaller,
		maxDelta:        maxDelta,
		numOfShards:     numOfShards,
		cleanupInterval: cleanupInterval,
	}

	bp.initRoundsMap()

	return bp, nil
}

func (bp *blocksPool) initRoundsMap() {
	roundsMap := make(map[uint32]uint64)
	for shardID := uint32(0); shardID < bp.numOfShards; shardID++ {
		roundsMap[shardID] = 0
	}
	roundsMap[core.MetachainShardId] = 0

	bp.roundsMap = roundsMap
}

func (bp *blocksPool) UpdateMetaState(round uint64) {
	bp.mutMap.Lock()
	defer bp.mutMap.Unlock()

	bp.roundsMap[core.MetachainShardId] = round

	err := bp.prunePersister(round)
	if err != nil {
		log.Warn("failed to prune persister", "error", err.Error())
	}
}

func (bp *blocksPool) prunePersister(round uint64) error {
	if round%bp.cleanupInterval != 0 {
		return nil
	}

	return bp.storer.Prune(round)
}

// PutBlock will put the provided outport block data to the pool
func (bp *blocksPool) PutBlock(hash []byte, outportBlock *outport.OutportBlock, currentRound uint64) error {
	bp.mutMap.Lock()
	defer bp.mutMap.Unlock()

	shardID := outportBlock.ShardID

	round, ok := bp.roundsMap[shardID]
	if !ok {
		return fmt.Errorf("did not find shard id %d in blocksMap", shardID)
	}

	if round == 0 {
		return bp.putOutportBlock(hash, outportBlock, currentRound)
	}

	metaRound := bp.roundsMap[core.MetachainShardId]

	if !bp.shouldPutOutportBlock(round, metaRound) {
		log.Error("failed to put outport block", "hash", hash, "round", round, "metaRound", metaRound)
		return fmt.Errorf("failed to put outport block")
	}

	return bp.putOutportBlock(hash, outportBlock, currentRound)
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
func (bp *blocksPool) putOutportBlock(
	hash []byte,
	outportBlock *outport.OutportBlock,
	currentRound uint64,
) error {
	shardID := outportBlock.ShardID

	outportBlockBytes, err := bp.marshaller.Marshal(outportBlock)
	if err != nil {
		return err
	}

	err = bp.storer.Put(hash, outportBlockBytes)
	if err != nil {
		return err
	}

	bp.roundsMap[shardID] = currentRound

	return nil
}

// GetBlock will return outport block data from the pool
func (bp *blocksPool) GetBlock(hash []byte) (*outport.OutportBlock, error) {
	data, err := bp.storer.Get(hash)
	if err != nil {
		return nil, err
	}

	outportBlock := &outport.OutportBlock{}
	err = bp.marshaller.Unmarshal(outportBlock, data)
	if err != nil {
		return nil, err
	}

	return outportBlock, nil
}

// IsInterfaceNil returns nil if there is no value under the interface
func (bp *blocksPool) IsInterfaceNil() bool {
	return bp == nil
}

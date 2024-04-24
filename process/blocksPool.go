package process

import (
	"fmt"
	"math"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
)

const initIndex = 0
const metaCheckpointKey = "lastMetaRound"

type blocksPool struct {
	storer               PruningStorer
	marshaller           marshal.Marshalizer
	maxDelta             uint64
	cleanupInterval      uint64
	firstCommitableBlock uint64

	previousIndexesMap map[uint32]uint64
	mutMap             sync.RWMutex
}

// NewBlocksPool will create a new blocks pool instance
func NewBlocksPool(
	storer PruningStorer,
	marshaller marshal.Marshalizer,
	maxDelta uint64,
	cleanupInterval uint64,
	firstCommitableBlock uint64,
) (*blocksPool, error) {
	if check.IfNil(storer) {
		return nil, ErrNilPruningStorer
	}
	if check.IfNil(marshaller) {
		return nil, ErrNilMarshaller
	}

	bp := &blocksPool{
		storer:               storer,
		marshaller:           marshaller,
		maxDelta:             maxDelta,
		cleanupInterval:      cleanupInterval,
		firstCommitableBlock: firstCommitableBlock,
	}

	bp.initIndexesMap()

	return bp, nil
}

func (bp *blocksPool) initIndexesMap() {
	lastCheckpoint, err := bp.getLastCheckpoint()
	if err != nil || lastCheckpoint == nil || lastCheckpoint.LastRounds == nil {
		indexesMap := make(map[uint32]uint64)
		bp.previousIndexesMap = indexesMap
		return
	}

	log.Info("initIndexesMap", "lastCheckpoint", lastCheckpoint)

	indexesMap := make(map[uint32]uint64)
	for shardID, round := range lastCheckpoint.LastRounds {
		indexesMap[shardID] = round
	}
	bp.previousIndexesMap = indexesMap
}

// Put will put value into storer
func (bp *blocksPool) Put(key []byte, value []byte) error {
	return bp.storer.Put(key, value)
}

// Get will get value from storer
func (bp *blocksPool) Get(key []byte) ([]byte, error) {
	return bp.storer.Get(key)
}

func (bp *blocksPool) UpdateMetaState(checkpoint *data.BlockCheckpoint) {
	index, ok := checkpoint.LastRounds[core.MetachainShardId]
	if !ok {
		index = initIndex
	}

	if index >= bp.firstCommitableBlock {
		err := bp.setCheckpoint(checkpoint)
		if err != nil {
			log.Warn("failed to set checkpoint", "error", err.Error())
		}
	}

	err := bp.pruneStorer(index)
	if err != nil {
		log.Warn("failed to prune storer", "error", err.Error())
	}
}

func (bp *blocksPool) pruneStorer(index uint64) error {
	if index%bp.cleanupInterval != 0 {
		return nil
	}

	return bp.storer.Prune(index)
}

// PutBlock will put the provided outport block data to the pool
func (bp *blocksPool) PutBlock(hash []byte, outportBlock *outport.OutportBlock, newIndex uint64) error {
	shardID := outportBlock.ShardID

	outportBlockBytes, err := bp.marshaller.Marshal(outportBlock)
	if err != nil {
		return err
	}

	bp.mutMap.Lock()
	defer bp.mutMap.Unlock()

	previousIndex, ok := bp.previousIndexesMap[shardID]
	if !ok {
		bp.previousIndexesMap[shardID] = initIndex
		previousIndex = initIndex
	}

	if previousIndex == initIndex {
		err := bp.storer.Put(hash, outportBlockBytes)
		if err != nil {
			return err
		}

		bp.previousIndexesMap[shardID] = newIndex

		return nil
	}

	isSuccesiveIndex := previousIndex+1 == newIndex
	if !isSuccesiveIndex {
		return fmt.Errorf("%w: new index should succesive, previous index %d, new index %d",
			ErrFailedToPutBlockDataToPool, previousIndex, newIndex)
	}

	if !bp.shouldPutBlockData(previousIndex) {
		return ErrFailedToPutBlockDataToPool
	}

	err = bp.storer.Put(hash, outportBlockBytes)
	if err != nil {
		return err
	}

	bp.previousIndexesMap[shardID] = newIndex

	return nil
}

func (bp *blocksPool) shouldPutBlockData(index uint64) bool {
	baseIndex := bp.previousIndexesMap[core.MetachainShardId]

	diff := float64(int64(index) - int64(baseIndex))
	delta := math.Abs(diff)

	return math.Abs(delta) < float64(bp.maxDelta)
}

// GetBlock will return outport block data from the pool
func (bp *blocksPool) GetBlock(hash []byte) (*outport.OutportBlock, error) {
	marshalledData, err := bp.storer.Get(hash)
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

func (bp *blocksPool) setCheckpoint(checkpoint *data.BlockCheckpoint) error {
	checkpointBytes, err := bp.marshaller.Marshal(checkpoint)
	if err != nil {
		return err
	}

	log.Debug("setCheckpoint", "checkpoint", checkpoint)

	return bp.storer.Put([]byte(metaCheckpointKey), checkpointBytes)
}

func (bp *blocksPool) getLastCheckpoint() (*data.BlockCheckpoint, error) {
	checkpointBytes, err := bp.storer.Get([]byte(metaCheckpointKey))
	if err != nil {
		return nil, err
	}

	checkpoint := &data.BlockCheckpoint{}
	err = bp.marshaller.Unmarshal(checkpoint, checkpointBytes)
	if err != nil {
		return nil, err
	}

	return checkpoint, nil
}

// Close will trigger close on blocks pool component
func (bp *blocksPool) Close() error {
	err := bp.storer.Close()
	if err != nil {
		return err
	}

	return nil
}

// IsInterfaceNil returns nil if there is no value under the interface
func (bp *blocksPool) IsInterfaceNil() bool {
	return bp == nil
}

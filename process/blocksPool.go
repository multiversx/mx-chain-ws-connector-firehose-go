package process

import (
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
	storer          PruningStorer
	marshaller      marshal.Marshalizer
	maxDelta        uint64
	cleanupInterval uint64

	indexesMap map[uint32]uint64
	mutMap     sync.RWMutex
}

// NewBlocksPool will create a new blocks pool instance
func NewBlocksPool(
	storer PruningStorer,
	marshaller marshal.Marshalizer,
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
		cleanupInterval: cleanupInterval,
	}

	bp.initIndexesMap()

	return bp, nil
}

func (bp *blocksPool) initIndexesMap() {
	lastCheckpoint, err := bp.getLastCheckpoint()
	if err != nil || lastCheckpoint == nil || lastCheckpoint.LastRounds == nil {
		bp.initIndexedMapFromStart()
		return
	}

	log.Info("initIndexesMap", "lastCheckpoint", lastCheckpoint)

	indexesMap := make(map[uint32]uint64)
	for shardID, round := range lastCheckpoint.LastRounds {
		indexesMap[shardID] = round
	}
	bp.indexesMap = indexesMap
}

func (bp *blocksPool) initIndexedMapFromStart() {
	log.Info("initialized blocks pool indexes map from start")

	indexesMap := make(map[uint32]uint64)
	indexesMap[core.MetachainShardId] = initIndex

	bp.indexesMap = indexesMap
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
	index := checkpoint.LastRounds[core.MetachainShardId]

	bp.mutMap.Lock()
	bp.indexesMap[core.MetachainShardId] = index
	bp.mutMap.Unlock()

	err := bp.setCheckpoint(checkpoint)
	if err != nil {
		log.Warn("failed to set checkpoint", "error", err.Error())
	}

	err = bp.pruneStorer(index)
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
func (bp *blocksPool) PutBlock(hash []byte, outportBlock *outport.OutportBlock, index uint64) error {
	shardID := outportBlock.ShardID

	outportBlockBytes, err := bp.marshaller.Marshal(outportBlock)
	if err != nil {
		return err
	}

	bp.mutMap.Lock()
	defer bp.mutMap.Unlock()

	currentIndex, ok := bp.indexesMap[shardID]
	if !ok {
		bp.indexesMap[shardID] = initIndex
		currentIndex = initIndex
	}

	if currentIndex == initIndex {
		err := bp.storer.Put(hash, outportBlockBytes)
		if err != nil {
			return err
		}

		bp.indexesMap[shardID] = index

		return nil
	}

	metaIndex := bp.indexesMap[core.MetachainShardId]

	if !bp.shouldPutBlockData(currentIndex, metaIndex) {
		return ErrFailedToPutBlockDataToPool
	}

	err = bp.storer.Put(hash, outportBlockBytes)
	if err != nil {
		return err
	}

	bp.indexesMap[shardID] = index

	return nil
}

func (bp *blocksPool) shouldPutBlockData(index, baseIndex uint64) bool {
	diff := float64(int64(index) - int64(baseIndex))
	delta := math.Abs(diff)

	if math.Abs(delta) >= float64(bp.maxDelta) {
		return false
	}

	return true
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

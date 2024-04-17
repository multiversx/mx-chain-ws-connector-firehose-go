package dataPool

import (
	"fmt"
	"math"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
)

var log = logger.GetOrCreate("dataPool")

const initIndex = 0

type blocksPool struct {
	storer          process.PruningStorer
	marshaller      marshal.Marshalizer
	maxDelta        uint64
	numOfShards     uint32
	cleanupInterval uint64

	indexesMap map[uint32]uint64
	mutMap     sync.RWMutex
}

func NewBlocksPool(
	storer process.PruningStorer,
	marshaller marshal.Marshalizer,
	numOfShards uint32,
	maxDelta uint64,
	cleanupInterval uint64,
) (*blocksPool, error) {
	if check.IfNil(storer) {
		return nil, process.ErrNilPruningStorer
	}
	if check.IfNil(marshaller) {
		return nil, process.ErrNilMarshaller
	}

	bp := &blocksPool{
		storer:          storer,
		marshaller:      marshaller,
		maxDelta:        maxDelta,
		numOfShards:     numOfShards,
		cleanupInterval: cleanupInterval,
	}

	bp.initIndexesMap()

	return bp, nil
}

func (bp *blocksPool) initIndexesMap() {
	indexesMap := make(map[uint32]uint64)
	for shardID := uint32(0); shardID < bp.numOfShards; shardID++ {
		indexesMap[shardID] = initIndex
	}
	indexesMap[core.MetachainShardId] = initIndex

	bp.indexesMap = indexesMap
}

func (bp *blocksPool) Put(key []byte, value []byte) error {
	return bp.storer.Put(key, value)
}

func (bp *blocksPool) Get(key []byte) ([]byte, error) {
	return bp.storer.Get(key)
}

func (bp *blocksPool) UpdateMetaState(index uint64) {
	bp.mutMap.Lock()
	defer bp.mutMap.Unlock()

	bp.indexesMap[core.MetachainShardId] = index

	err := bp.storer.SetCheckpoint(index)
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
func (bp *blocksPool) PutBlock(hash []byte, value []byte, index uint64, shardID uint32) error {
	bp.mutMap.Lock()
	defer bp.mutMap.Unlock()

	currentIndex, ok := bp.indexesMap[shardID]
	if !ok {
		return fmt.Errorf("did not find shard id %d in blocksMap", shardID)
	}

	if currentIndex == initIndex {
		err := bp.storer.Put(hash, value)
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

	err := bp.storer.Put(hash, value)
	if err != nil {
		return err
	}

	bp.indexesMap[shardID] = index

	return nil
}

// should be run under mutex
func (bp *blocksPool) shouldPutBlockData(index, baseIndex uint64) bool {
	diff := float64(int64(index) - int64(baseIndex))
	delta := math.Abs(diff)

	if math.Abs(delta) >= float64(bp.maxDelta) {
		return false
	}

	return true
}

// GetBlock will return outport block data from the pool
func (bp *blocksPool) GetBlock(hash []byte) ([]byte, error) {
	return bp.storer.Get(hash)
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

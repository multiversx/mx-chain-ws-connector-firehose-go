package process

import (
	"math"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
)

const initIndex = 0

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

func (bp *blocksPool) UpdateMetaState(index uint64) {
	bp.mutMap.Lock()
	bp.indexesMap[core.MetachainShardId] = index
	bp.mutMap.Unlock()

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

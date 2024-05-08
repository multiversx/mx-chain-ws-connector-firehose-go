package process

import (
	"fmt"
	"math"
	"sync"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data"
)

const (
	initIndex         = 0
	metaCheckpointKey = "lastMetaRound"
	minDelta          = 3
)

type blocksPool struct {
	storer               PruningStorer
	marshaller           marshal.Marshalizer
	maxDelta             uint64
	cleanupInterval      uint64
	firstCommitableBlock uint64

	previousIndexesMap map[uint32]uint64
	mutMap             sync.RWMutex
}

// BlocksPoolArgs defines the arguments needed to create the blocks pool component
type BlocksPoolArgs struct {
	Storer               PruningStorer
	Marshaller           marshal.Marshalizer
	MaxDelta             uint64
	CleanupInterval      uint64
	FirstCommitableBlock uint64
}

// NewBlocksPool will create a new blocks pool instance
func NewBlocksPool(args BlocksPoolArgs) (*blocksPool, error) {
	if check.IfNil(args.Storer) {
		return nil, ErrNilPruningStorer
	}
	if check.IfNil(args.Marshaller) {
		return nil, ErrNilMarshaller
	}
	if args.MaxDelta < minDelta {
		return nil, fmt.Errorf("%w for max delta, provided %d, min required %d", ErrInvalidValue, args.MaxDelta, minDelta)
	}
	if args.CleanupInterval <= args.MaxDelta {
		return nil, fmt.Errorf("%w for cleanup interval, provided %d, min required %d", ErrInvalidValue, args.MaxDelta, args.MaxDelta)
	}

	bp := &blocksPool{
		storer:               args.Storer,
		marshaller:           args.Marshaller,
		maxDelta:             args.MaxDelta,
		cleanupInterval:      args.CleanupInterval,
		firstCommitableBlock: args.FirstCommitableBlock,
	}

	bp.initIndexesMap()

	return bp, nil
}

func (bp *blocksPool) initIndexesMap() {
	lastCheckpoint, err := bp.getLastCheckpoint()
	if err != nil || lastCheckpoint == nil || lastCheckpoint.LastNonces == nil {
		indexesMap := make(map[uint32]uint64)
		bp.previousIndexesMap = indexesMap
		return
	}

	log.Info("initIndexesMap", "lastCheckpoint", lastCheckpoint)

	indexesMap := make(map[uint32]uint64, len(lastCheckpoint.LastNonces))
	for shardID, round := range lastCheckpoint.LastNonces {
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

// UpdateMetaState will update internal meta state
func (bp *blocksPool) UpdateMetaState(checkpoint *data.BlockCheckpoint) error {
	index, ok := checkpoint.LastNonces[core.MetachainShardId]
	if !ok {
		index = initIndex
	}

	if index >= bp.firstCommitableBlock {
		err := bp.setCheckpoint(checkpoint)
		if err != nil {
			return fmt.Errorf("%w, failed to set checkpoint", err)
		}
	}

	err := bp.pruneStorer(index)
	if err != nil {
		return fmt.Errorf("%w, failed to prune storer", err)
	}

	return nil
}

func (bp *blocksPool) pruneStorer(index uint64) error {
	if index%bp.cleanupInterval != 0 {
		return nil
	}

	return bp.storer.Prune(index)
}

func (bp *blocksPool) getLastIndex(shardID uint32) uint64 {
	lastCheckpoint, err := bp.getLastCheckpoint()
	if err == nil {
		baseIndex, ok := lastCheckpoint.LastNonces[shardID]
		if ok {
			return baseIndex
		}
	}

	lastIndex, ok := bp.previousIndexesMap[shardID]
	if !ok {
		return initIndex
	}

	return lastIndex
}

// PutBlock will put the provided outport block data to the pool
func (bp *blocksPool) PutBlock(hash []byte, value []byte, newIndex uint64, shardID uint32) error {
	bp.mutMap.Lock()
	defer bp.mutMap.Unlock()

	lastIndex := bp.getLastIndex(shardID)
	if lastIndex == initIndex {
		bp.previousIndexesMap[shardID] = newIndex
		return bp.storer.Put(hash, value)
	}

	if !bp.shouldPutBlockData(newIndex, lastIndex) {
		return fmt.Errorf("%w: not within required delta, last index %d, new index %d",
			ErrFailedToPutBlockDataToPool, lastIndex, newIndex)
	}

	err := bp.storer.Put(hash, value)
	if err != nil {
		return fmt.Errorf("failed to put into storer: %w", err)
	}

	bp.previousIndexesMap[shardID] = newIndex

	return nil
}

func (bp *blocksPool) shouldPutBlockData(index, baseIndex uint64) bool {
	diff := float64(int64(index) - int64(baseIndex))
	delta := math.Abs(diff)

	return math.Abs(delta) < float64(bp.maxDelta)
}

func (bp *blocksPool) setCheckpoint(checkpoint *data.BlockCheckpoint) error {
	checkpointBytes, err := bp.marshaller.Marshal(checkpoint)
	if err != nil {
		return fmt.Errorf("failed to marshall checkpoint data: %w", err)
	}

	log.Debug("setCheckpoint", "checkpoint", checkpoint)

	return bp.storer.Put([]byte(metaCheckpointKey), checkpointBytes)
}

func (bp *blocksPool) getLastCheckpoint() (*data.BlockCheckpoint, error) {
	checkpointBytes, err := bp.storer.Get([]byte(metaCheckpointKey))
	if err != nil {
		return nil, fmt.Errorf("failed to get checkpoint data from storer: %w", err)
	}

	checkpoint := &data.BlockCheckpoint{}
	err = bp.marshaller.Unmarshal(checkpoint, checkpointBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall checkpoint data: %w", err)
	}

	if checkpoint == nil || checkpoint.LastNonces == nil {
		return nil, fmt.Errorf("nil checkpoint data has been provided")
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

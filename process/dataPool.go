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
	// MetaCheckpointKey defines the meta checkpoint key
	MetaCheckpointKey = "lastMetaIndex"

	softMetaCheckpointKey = "softLastCheckpoint"

	initIndex = 0
	minDelta  = 3
)

type dataPool struct {
	storer                PruningStorer
	marshaller            marshal.Marshalizer
	maxDelta              uint64
	cleanupInterval       uint64
	firstCommitableBlocks map[uint32]uint64

	// soft checkpoint data is being used only at startup (without any previous data)
	// until there is a valid hard checkpoint (set by publisher)
	softCheckpointMap map[uint32]uint64
	mutMap            sync.RWMutex
}

// DataPoolArgs defines the arguments needed to create the blocks pool component
type DataPoolArgs struct {
	Storer                PruningStorer
	Marshaller            marshal.Marshalizer
	MaxDelta              uint64
	CleanupInterval       uint64
	FirstCommitableBlocks map[uint32]uint64
}

// NewDataPool will create a new data pool instance
func NewDataPool(args DataPoolArgs) (*dataPool, error) {
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

	bp := &dataPool{
		storer:                args.Storer,
		marshaller:            args.Marshaller,
		maxDelta:              args.MaxDelta,
		cleanupInterval:       args.CleanupInterval,
		firstCommitableBlocks: args.FirstCommitableBlocks,
	}

	bp.initIndexesMap()

	return bp, nil
}

func (bp *dataPool) initIndexesMap() {
	lastCheckpoint, err := bp.GetLastCheckpoint()
	if err == nil {
		log.Info("initIndexesMap", "lastCheckpoint", lastCheckpoint)
		bp.softCheckpointMap = lastCheckpoint.GetLastNonces()
		return
	}

	log.Warn("failed to get last checkpoint, will try to set soft last checkpoint", "error", err)

	softCheckpoint, err := bp.getLastSoftCheckpoint()
	if err == nil {
		log.Info("initIndexesMap", "softCheckpoint", softCheckpoint)
		bp.softCheckpointMap = softCheckpoint.GetLastNonces()
		return
	}

	log.Warn("failed to get last soft checkpoint, will set empty soft checkpoint", "error", err)
	bp.softCheckpointMap = make(map[uint32]uint64)
}

// Put will put value into storer
func (bp *dataPool) Put(key []byte, value []byte) error {
	return bp.storer.Put(key, value)
}

// Get will get value from storer
func (bp *dataPool) Get(key []byte) ([]byte, error) {
	return bp.storer.Get(key)
}

// UpdateMetaState will update internal meta state
func (bp *dataPool) UpdateMetaState(checkpoint *data.BlockCheckpoint) error {
	index, ok := checkpoint.LastNonces[core.MetachainShardId]
	if !ok {
		index = initIndex
	}

	if index >= bp.firstCommitableBlocks[core.MetachainShardId] {
		err := bp.setCheckpoint(checkpoint)
		if err != nil {
			return fmt.Errorf("%w, failed to set checkpoint", err)
		}
	}

	// TODO: prune storer based on epoch change
	err := bp.pruneStorer(index)
	if err != nil {
		return fmt.Errorf("%w, failed to prune storer", err)
	}

	return nil
}

func (bp *dataPool) pruneStorer(index uint64) error {
	if index%bp.cleanupInterval != 0 {
		return nil
	}

	return bp.storer.Prune(index)
}

func (bp *dataPool) getLastIndex(shardID uint32) uint64 {
	lastCheckpoint, err := bp.GetLastCheckpoint()
	if err == nil {
		baseIndex, ok := lastCheckpoint.LastNonces[shardID]
		if ok {
			return baseIndex
		}
	}

	lastIndex, ok := bp.softCheckpointMap[shardID]
	if !ok {
		return initIndex
	}

	return lastIndex
}

// PutBlock will put the provided outport block data to the pool
func (bp *dataPool) PutBlock(hash []byte, value []byte, newIndex uint64, shardID uint32) error {
	bp.mutMap.Lock()
	defer bp.mutMap.Unlock()

	firstCommitableBlock, ok := bp.firstCommitableBlocks[shardID]
	if !ok {
		return fmt.Errorf("failed to get first commitable block for shard %d", shardID)
	}

	if newIndex < firstCommitableBlock {
		log.Trace("do not commit block", "newIndex", newIndex, "firstCommitableBlock", firstCommitableBlock)
		return nil
	}

	lastIndex := bp.getLastIndex(shardID)
	if lastIndex == initIndex {
		bp.softCheckpointMap[shardID] = newIndex
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

	return nil
}

func (bp *dataPool) shouldPutBlockData(index, baseIndex uint64) bool {
	diff := float64(int64(index) - int64(baseIndex))
	delta := math.Abs(diff)

	return math.Abs(delta) < float64(bp.maxDelta)
}

func (bp *dataPool) setCheckpoint(checkpoint *data.BlockCheckpoint) error {
	checkpointBytes, err := bp.marshaller.Marshal(checkpoint)
	if err != nil {
		return fmt.Errorf("failed to marshall checkpoint data: %w", err)
	}

	log.Debug("setCheckpoint", "checkpoint", checkpoint)

	return bp.storer.Put([]byte(MetaCheckpointKey), checkpointBytes)
}

// GetLastCheckpoint returns last checkpoint data
func (bp *dataPool) GetLastCheckpoint() (*data.BlockCheckpoint, error) {
	return bp.getCheckpointData(MetaCheckpointKey)
}

func (bp *dataPool) getLastSoftCheckpoint() (*data.BlockCheckpoint, error) {
	return bp.getCheckpointData(softMetaCheckpointKey)
}

func (bp *dataPool) getCheckpointData(checkpointKey string) (*data.BlockCheckpoint, error) {
	checkpointBytes, err := bp.storer.Get([]byte(checkpointKey))
	if err != nil {
		return nil, fmt.Errorf("failed to get checkpoint data (key = %s) from storer: %w", checkpointKey, err)
	}

	checkpoint := &data.BlockCheckpoint{}
	err = bp.marshaller.Unmarshal(checkpoint, checkpointBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall checkpoint (key = %s) data: %w", checkpointKey, err)
	}

	if checkpoint == nil || checkpoint.LastNonces == nil {
		return nil, fmt.Errorf("nil checkpoint data (key = %s) has been provided", checkpointKey)
	}

	return checkpoint, nil
}

func (bp *dataPool) saveLastSoftCheckpoint() error {
	bp.mutMap.RLock()
	softCheckpointMap := bp.softCheckpointMap
	bp.mutMap.RUnlock()

	softCheckpoint := &data.BlockCheckpoint{}
	softCheckpoint.LastNonces = softCheckpointMap

	checkpointBytes, err := bp.marshaller.Marshal(softCheckpoint)
	if err != nil {
		return fmt.Errorf("failed to marshall publish soft checkpoint data: %w", err)
	}

	log.Debug("saveLastSoftCheckpoint", "previousIndexesMap", softCheckpointMap)

	return bp.storer.Put([]byte(softMetaCheckpointKey), checkpointBytes)
}

// Close will trigger close on blocks pool component
func (bp *dataPool) Close() error {
	err := bp.saveLastSoftCheckpoint()
	if err != nil {
		return err
	}

	err = bp.storer.Close()
	if err != nil {
		return err
	}

	return nil
}

// IsInterfaceNil returns nil if there is no value under the interface
func (bp *dataPool) IsInterfaceNil() bool {
	return bp == nil
}

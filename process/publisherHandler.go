package process

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
)

const (
	publishCheckpointKey = "publishCheckpoint"

	minRetryDudrationInMilliseconds = 100
)

type publisherHandler struct {
	handler               HyperBlockPublisher
	outportBlocksPool     BlocksPool
	dataAggregator        DataAggregator
	retryDuration         time.Duration
	marshaller            marshal.Marshalizer
	firstCommitableBlocks map[uint32]uint64

	blocksChan chan []byte
	cancelFunc func()
	closeChan  chan struct{}

	checkpoint    *data.PublishCheckpoint
	checkpointMut sync.RWMutex
}

// PublisherHandlerArgs defines the argsuments needed to create a new publisher handler
type PublisherHandlerArgs struct {
	Handler                     HyperBlockPublisher
	OutportBlocksPool           BlocksPool
	DataAggregator              DataAggregator
	RetryDurationInMilliseconds uint64
	Marshalizer                 marshal.Marshalizer
	FirstCommitableBlocks       map[uint32]uint64
}

// NewPublisherHandler creates a new publisher handler component
func NewPublisherHandler(args PublisherHandlerArgs) (*publisherHandler, error) {
	if check.IfNil(args.Handler) {
		return nil, ErrNilPublisher
	}
	if check.IfNil(args.OutportBlocksPool) {
		return nil, ErrNilBlocksPool
	}
	if check.IfNil(args.DataAggregator) {
		return nil, ErrNilDataAggregator
	}
	if check.IfNil(args.Marshalizer) {
		return nil, ErrNilMarshaller
	}
	if args.FirstCommitableBlocks == nil {
		return nil, ErrNilFirstCommitableBlocks
	}
	if args.RetryDurationInMilliseconds < minRetryDudrationInMilliseconds {
		return nil, fmt.Errorf("%w for retry duration: provided %d, min required %d",
			ErrInvalidValue, args.RetryDurationInMilliseconds, minRetryDudrationInMilliseconds)
	}

	ph := &publisherHandler{
		handler:               args.Handler,
		outportBlocksPool:     args.OutportBlocksPool,
		dataAggregator:        args.DataAggregator,
		marshaller:            args.Marshalizer,
		retryDuration:         time.Duration(args.RetryDurationInMilliseconds) * time.Millisecond,
		firstCommitableBlocks: args.FirstCommitableBlocks,
		checkpoint:            &data.PublishCheckpoint{},
		blocksChan:            make(chan []byte),
		closeChan:             make(chan struct{}),
	}

	var ctx context.Context
	ctx, ph.cancelFunc = context.WithCancel(context.Background())

	go ph.startListener(ctx)

	return ph, nil
}

func (ph *publisherHandler) startListener(ctx context.Context) {
	log.Debug("starting publisherHandler listerer")

	err := ph.handleLastCheckpointOnInit()
	if err != nil {
		log.Error("failed to handle last publish checkpoint on init", "error", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Debug("closing publisherHandler listener")
			return
		case headerHash := <-ph.blocksChan:
			ph.handlePublishEvent(headerHash)
		}
	}
}

// PublishBlock will push header hash to blocks channel
func (ph *publisherHandler) PublishBlock(headerHash []byte) error {
	select {
	case ph.blocksChan <- headerHash:
	case <-ph.closeChan:
	}

	ph.setPublishCheckpoint(headerHash)

	return nil
}

func (ph *publisherHandler) handlePublishEvent(headerHash []byte) {
	for {
		err := ph.handlerHyperOutportBlock(headerHash)
		if err == nil {
			log.Info("published hyper block", "headerHash", headerHash)
			return
		}

		log.Error("failed to publish hyper block event", "headerHash", headerHash, "error", err)
		time.Sleep(ph.retryDuration)
	}
}

func (ph *publisherHandler) setPublishCheckpoint(headerHash []byte) {
	publishCheckpoint := &data.PublishCheckpoint{
		HeaderHash: headerHash,
		Published:  false,
	}

	ph.checkpointMut.Lock()
	ph.checkpoint = publishCheckpoint
	ph.checkpointMut.Unlock()
}

func (ph *publisherHandler) updatePublishCheckpoint() {
	ph.checkpointMut.Lock()
	if ph.checkpoint != nil {
		ph.checkpoint.Published = true
	}
	ph.checkpointMut.Unlock()
}

func (ph *publisherHandler) handleLastCheckpointOnInit() error {
	checkpoint, err := ph.getLastPublishCheckpoint()
	if err != nil {
		return err
	}

	if checkpoint.Published {
		return nil
	}

	log.Debug("trying to publish from last checkpoint", "headerHash", checkpoint.HeaderHash)

	ph.handlePublishEvent(checkpoint.HeaderHash)

	return nil
}

func (ph *publisherHandler) getLastPublishCheckpoint() (*data.PublishCheckpoint, error) {
	checkpointBytes, err := ph.outportBlocksPool.Get([]byte(publishCheckpointKey))
	if err != nil || len(checkpointBytes) == 0 {
		// if there is no last publish checkpoint we assume that it was not set;
		// in case it was not saved properly to storage on close, we need to reprocess/import it manually

		// TODO: evaluate adding full persister mode - saving/updating publishCheckpoint state
		// to storage on each event

		log.Trace("no previous last publish checkpoint, probably starting from scratch")

		return &data.PublishCheckpoint{
			Published: true,
		}, nil
	}

	checkpoint := &data.PublishCheckpoint{}
	err = ph.marshaller.Unmarshal(checkpoint, checkpointBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall publish checkpoint data: %w", err)
	}

	return checkpoint, nil
}

func (ph *publisherHandler) savePublishCheckpoint() error {
	ph.checkpointMut.RLock()
	publishCheckpoint := ph.checkpoint
	ph.checkpointMut.RUnlock()

	if publishCheckpoint == nil {
		log.Trace("no previous last publish checkpoint, probably closing without handling any events")
		return nil
	}

	checkpointBytes, err := ph.marshaller.Marshal(publishCheckpoint)
	if err != nil {
		return fmt.Errorf("failed to marshall publish checkpoint data: %w", err)
	}

	log.Debug("savePublishCheckpoint", "headerHash", publishCheckpoint.HeaderHash)

	return ph.outportBlocksPool.Put([]byte(publishCheckpointKey), checkpointBytes)
}

func (ph *publisherHandler) handlerHyperOutportBlock(headerHash []byte) error {
	metaOutportBlock, err := ph.outportBlocksPool.GetMetaBlock(headerHash)
	if err != nil {
		return err
	}

	err = checkMetaOutportBlockHeader(metaOutportBlock)
	if err != nil {
		return err
	}

	metaNonce := metaOutportBlock.BlockData.Header.GetNonce()
	shardID := metaOutportBlock.GetShardID()

	firstCommitableBlock, ok := ph.firstCommitableBlocks[shardID]
	if !ok {
		return fmt.Errorf("failed to get first commitable block for shard %d", shardID)
	}

	if metaNonce < firstCommitableBlock {
		// do not try to aggregate or publish hyper outport block

		log.Trace("do not commit block", "currentNonce", metaNonce, "firstCommitableNonce", firstCommitableBlock)

		return nil
	}

	hyperOutportBlock, err := ph.dataAggregator.ProcessHyperBlock(metaOutportBlock)
	if err != nil {
		return err
	}

	lastBlockCheckpoint, err := ph.getLastBlockCheckpoint(hyperOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to get last round data: %w", err)
	}

	err = ph.outportBlocksPool.UpdateMetaState(lastBlockCheckpoint)
	if err != nil {
		return err
	}

	err = ph.handler.PublishHyperBlock(hyperOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to publish hyperblock: %w", err)
	}

	ph.updatePublishCheckpoint()

	return nil
}

func (ph *publisherHandler) getLastBlockCheckpoint(hyperOutportBlock *hyperOutportBlocks.HyperOutportBlock) (*data.BlockCheckpoint, error) {
	if hyperOutportBlock == nil {
		return nil, ErrNilHyperOutportBlock
	}
	err := checkMetaOutportBlockHeader(hyperOutportBlock.MetaOutportBlock)
	if err != nil {
		return nil, err
	}

	checkpoint := &data.BlockCheckpoint{
		LastNonces: make(map[uint32]uint64),
	}

	metaBlock := hyperOutportBlock.MetaOutportBlock.BlockData.Header
	checkpoint.LastNonces[core.MetachainShardId] = metaBlock.GetNonce()

	for _, outportBlockData := range hyperOutportBlock.NotarizedHeadersOutportData {
		header := outportBlockData.OutportBlock.BlockData.Header
		checkpoint.LastNonces[outportBlockData.OutportBlock.ShardID] = header.GetNonce()
	}

	return checkpoint, nil
}

func checkMetaOutportBlockHeader(metaOutportBlock *hyperOutportBlocks.MetaOutportBlock) error {
	if metaOutportBlock == nil {
		return fmt.Errorf("%w for metaOutportBlock", ErrNilOutportBlockData)
	}
	if metaOutportBlock.BlockData == nil {
		return fmt.Errorf("%w for blockData", ErrNilOutportBlockData)
	}
	if metaOutportBlock.BlockData.Header == nil {
		return fmt.Errorf("%w for blockData header", ErrNilOutportBlockData)
	}

	return nil
}

// Close will close the internal writer
func (ph *publisherHandler) Close() error {
	err := ph.savePublishCheckpoint()
	if err != nil {
		return err
	}

	err = ph.handler.Close()
	if err != nil {
		return err
	}

	ph.cancelFunc()

	close(ph.closeChan)

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (ph *publisherHandler) IsInterfaceNil() bool {
	return ph == nil
}

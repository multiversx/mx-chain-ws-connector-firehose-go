package process

import (
	"context"
	"fmt"
	"time"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
)

type publisherHandler struct {
	handler               HyperBlockPublisher
	outportBlocksPool     BlocksPool
	dataAggregator        DataAggregator
	retryDuration         time.Duration
	firstCommitableBlocks map[uint32]uint64

	blocksChan chan []byte
	cancelFunc func()
	closeChan  chan struct{}
}

// PublisherHandlerArgs defines the argsuments needed to create a new publisher handler
type PublisherHandlerArgs struct {
	Handler                     HyperBlockPublisher
	OutportBlocksPool           BlocksPool
	DataAggregator              DataAggregator
	RetryDurationInMilliseconds uint64
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
	if args.FirstCommitableBlocks == nil {
		return nil, ErrNilFirstCommitableBlocks
	}

	ph := &publisherHandler{
		handler:               args.Handler,
		outportBlocksPool:     args.OutportBlocksPool,
		dataAggregator:        args.DataAggregator,
		retryDuration:         time.Duration(args.RetryDurationInMilliseconds) * time.Millisecond,
		firstCommitableBlocks: args.FirstCommitableBlocks,
		blocksChan:            make(chan []byte),
		closeChan:             make(chan struct{}),
	}

	var ctx context.Context
	ctx, ph.cancelFunc = context.WithCancel(context.Background())

	go ph.startListener(ctx)

	return ph, nil
}

func (ph *publisherHandler) startListener(ctx context.Context) {
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

	return nil
}

func (ph *publisherHandler) handlePublishEvent(headerHash []byte) {
	for {
		err := ph.handlerHyperOutportBlock(headerHash)
		if err == nil {
			return
		}

		log.Error("failed to publish hyper block event", "error", err)
		time.Sleep(ph.retryDuration)
	}
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

	lastCheckpoint, err := ph.getLastCheckpointData(hyperOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to get last round data: %w", err)
	}

	err = ph.handler.PublishHyperBlock(hyperOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to publish hyperblock: %w", err)
	}

	return ph.outportBlocksPool.UpdateMetaState(lastCheckpoint)
}

func (ph *publisherHandler) getLastCheckpointData(hyperOutportBlock *hyperOutportBlocks.HyperOutportBlock) (*data.BlockCheckpoint, error) {
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
	err := ph.handler.Close()
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

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
	handler              HyperBlockPublisher
	outportBlocksPool    HyperBlocksPool
	dataAggregator       DataAggregator
	retryDuration        time.Duration
	firstCommitableBlock uint64

	blocksChan chan []byte
	cancelFunc func()
	closeChan  chan struct{}
}

// NewPublisherHandler creates a new publisher handler component
func NewPublisherHandler(
	handler HyperBlockPublisher,
	outportBlocksPool HyperBlocksPool,
	dataAggregator DataAggregator,
	retryDurationInMiliseconds int64,
	firstCommitableBlock uint64,
) (*publisherHandler, error) {
	if check.IfNil(handler) {
		return nil, ErrNilPublisher
	}
	if check.IfNil(outportBlocksPool) {
		return nil, ErrNilHyperBlocksPool
	}
	if check.IfNil(dataAggregator) {
		return nil, ErrNilDataAggregator
	}

	ph := &publisherHandler{
		handler:              handler,
		outportBlocksPool:    outportBlocksPool,
		dataAggregator:       dataAggregator,
		retryDuration:        time.Duration(retryDurationInMiliseconds) * time.Millisecond,
		firstCommitableBlock: firstCommitableBlock,
		blocksChan:           make(chan []byte),
		closeChan:            make(chan struct{}),
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
			log.Debug("closing commonPublisher listener")
			return
		case headerHash := <-ph.blocksChan:
			ph.handlePublishEvent(headerHash)
		}
	}
}

// PublishHyperBlock will push aggregated outport block data to the firehose writer
func (ph *publisherHandler) PublishBlock(headerHash []byte) error {
	select {
	case ph.blocksChan <- headerHash:
	case <-ph.closeChan:
	}

	return nil
}

func (ph *publisherHandler) handlePublishEvent(headerHash []byte) {
	// TODO: evaluate max retries and exit failure
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

	metaRound := metaOutportBlock.BlockData.Header.GetRound()

	if metaRound < ph.firstCommitableBlock {
		// do not try to aggregate or publish hyper outport block
		// update only blocks pool state

		log.Trace("do not commit block", "currentRound", metaRound, "firstCommitableRound", ph.firstCommitableBlock)

		lastCheckpoint := &data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				core.MetachainShardId: metaRound,
			},
		}
		err := ph.outportBlocksPool.UpdateMetaState(lastCheckpoint)
		if err != nil {
			return err
		}

		return nil
	}

	hyperOutportBlock, err := ph.dataAggregator.ProcessHyperBlock(metaOutportBlock)
	if err != nil {
		return err
	}

	lastCheckpoint, err := ph.getLastRoundsData(hyperOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to get last round data: %w", err)
	}

	err = ph.handler.PublishHyperBlock(hyperOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to publish hyperblock: %w", err)
	}

	return ph.outportBlocksPool.UpdateMetaState(lastCheckpoint)
}

func (ph *publisherHandler) getLastRoundsData(hyperOutportBlock *hyperOutportBlocks.HyperOutportBlock) (*data.BlockCheckpoint, error) {
	if hyperOutportBlock == nil {
		return nil, ErrNilHyperOutportBlock
	}
	err := checkMetaOutportBlockHeader(hyperOutportBlock.MetaOutportBlock)
	if err != nil {
		return nil, err
	}

	checkpoint := &data.BlockCheckpoint{
		LastRounds: make(map[uint32]uint64),
	}

	metaBlock := hyperOutportBlock.MetaOutportBlock.BlockData.Header
	checkpoint.LastRounds[core.MetachainShardId] = metaBlock.GetRound()

	for _, outportBlockData := range hyperOutportBlock.NotarizedHeadersOutportData {
		header := outportBlockData.OutportBlock.BlockData.Header
		checkpoint.LastRounds[outportBlockData.OutportBlock.ShardID] = header.GetNonce()
	}

	return checkpoint, nil
}

func checkMetaOutportBlockHeader(metaOutportBlock *hyperOutportBlocks.MetaOutportBlock) error {
	if metaOutportBlock == nil {
		return fmt.Errorf("%w for metaOutportBlock", ErrNilHyperOutportBlock)
	}
	if metaOutportBlock.BlockData == nil {
		return fmt.Errorf("%w for blockData", ErrNilHyperOutportBlock)
	}
	if metaOutportBlock.BlockData.Header == nil {
		return fmt.Errorf("%w for blockData header", ErrNilHyperOutportBlock)
	}

	return nil
}

// Close will close the internal writer
func (ph *publisherHandler) Close() error {
	err := ph.handler.Close()
	if err != nil {
		return err
	}

	if ph.cancelFunc != nil {
		ph.cancelFunc()
	}

	close(ph.closeChan)

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (ph *publisherHandler) IsInterfaceNil() bool {
	return ph == nil
}

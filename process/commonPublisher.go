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

type commonPublisher struct {
	handler              Publisher
	outportBlocksPool    HyperBlocksPool
	dataAggregator       DataAggregator
	retryDuration        time.Duration
	firstCommitableBlock uint64

	blocksChan chan []byte
	cancelFunc func()
	closeChan  chan struct{}
}

// NewPublisher creates a new publisher component
func NewPublisher(
	handler Publisher,
	outportBlocksPool HyperBlocksPool,
	dataAggregator DataAggregator,
	retryDurationInMiliseconds int64,
	firstCommitableBlock uint64,
) (*commonPublisher, error) {
	if check.IfNil(outportBlocksPool) {
		return nil, ErrNilHyperBlocksPool
	}
	if check.IfNil(handler) {
		return nil, ErrNilPublisher
	}
	if check.IfNil(dataAggregator) {
		return nil, ErrNilDataAggregator
	}

	cp := &commonPublisher{
		handler:              handler,
		retryDuration:        time.Duration(retryDurationInMiliseconds) * time.Millisecond,
		firstCommitableBlock: firstCommitableBlock,
	}

	var ctx context.Context
	ctx, cp.cancelFunc = context.WithCancel(context.Background())

	go cp.startListener(ctx)

	return cp, nil
}

func (cp *commonPublisher) startListener(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Debug("closing commonPublisher listener")
			return
		case headerHash := <-cp.blocksChan:
			cp.handlePublishEvent(headerHash)
		}
	}
}

func (cp *commonPublisher) handlePublishEvent(headerHash []byte) {
	for {
		err := cp.handlerHyperOutportBlock(headerHash)
		if err == nil {
			return
		}

		log.Error("failed to publish hyper block event", "error", err)
		time.Sleep(cp.retryDuration)
	}
}

func (cp *commonPublisher) handlerHyperOutportBlock(headerHash []byte) error {
	metaOutportBlock, err := cp.outportBlocksPool.GetMetaBlock(headerHash)
	if err != nil {
		return nil
	}

	metaRound := metaOutportBlock.BlockData.Header.GetRound()

	if metaRound < cp.firstCommitableBlock {
		// do not try to aggregate or publish hyper outport block
		// update only blocks pool state

		log.Trace("do not commit block", "currentRound", metaRound, "firstCommitableRound", cp.firstCommitableBlock)

		lastCheckpoint := &data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				core.MetachainShardId: metaRound,
			},
		}
		err := cp.outportBlocksPool.UpdateMetaState(lastCheckpoint)
		if err != nil {
			return err
		}

		return nil
	}

	hyperOutportBlock, err := cp.dataAggregator.ProcessHyperBlock(metaOutportBlock)
	if err != nil {
		return err
	}

	lastCheckpoint, err := cp.getLastRoundsData(hyperOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to get last round data: %w", err)
	}

	err = cp.handler.PublishHyperBlock(hyperOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to publish hyperblock: %w", err)
	}

	return cp.outportBlocksPool.UpdateMetaState(lastCheckpoint)
}

func (cp *commonPublisher) getLastRoundsData(hyperOutportBlock *hyperOutportBlocks.HyperOutportBlock) (*data.BlockCheckpoint, error) {
	if hyperOutportBlock == nil {
		return nil, ErrNilHyperOutportBlock
	}
	if hyperOutportBlock.MetaOutportBlock == nil {
		return nil, fmt.Errorf("%w for metaOutportBlock", ErrNilHyperOutportBlock)
	}
	if hyperOutportBlock.MetaOutportBlock.BlockData == nil {
		return nil, fmt.Errorf("%w for blockData", ErrNilHyperOutportBlock)
	}
	if hyperOutportBlock.MetaOutportBlock.BlockData.Header == nil {
		return nil, fmt.Errorf("%w for blockData header", ErrNilHyperOutportBlock)
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

// PublishHyperBlock will push aggregated outport block data to the firehose writer
func (cp *commonPublisher) PublishBlock(headerHash []byte) error {
	select {
	case cp.blocksChan <- headerHash:
	case <-cp.closeChan:
	}

	return nil
}

// Close will close the internal writer
func (cp *commonPublisher) Close() error {
	if cp.cancelFunc != nil {
		cp.cancelFunc()
	}

	close(cp.closeChan)

	return nil
}

// IsInterfaceNil returns true if there is no value under the interface
func (cp *commonPublisher) IsInterfaceNil() bool {
	return cp == nil
}

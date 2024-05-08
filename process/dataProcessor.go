package process

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
)

type dataProcessor struct {
	marshaller            marshal.Marshalizer
	operationHandlers     map[string]func(marshalledData []byte) error
	publisher             Publisher
	outportBlocksPool     BlocksPool
	dataAggregator        DataAggregator
	outportBlockConverter OutportBlockConverter
	firstCommitableBlock  uint64
}

// NewDataProcessor creates a data processor able to receive data from a ws outport driver and handle blocks
func NewDataProcessor(
	publisher Publisher,
	marshaller marshal.Marshalizer,
	blocksPool BlocksPool,
	dataAggregator DataAggregator,
	outportBlockConverter OutportBlockConverter,
	firstCommitableBlock uint64,
) (DataProcessor, error) {
	if check.IfNil(publisher) {
		return nil, ErrNilPublisher
	}
	if check.IfNil(blocksPool) {
		return nil, ErrNilBlocksPool
	}
	if check.IfNil(marshaller) {
		return nil, ErrNilMarshaller
	}
	if check.IfNil(dataAggregator) {
		return nil, ErrNilDataAggregator
	}
	if check.IfNil(outportBlockConverter) {
		return nil, ErrNilOutportBlocksConverter
	}

	dp := &dataProcessor{
		marshaller:            marshaller,
		publisher:             publisher,
		outportBlocksPool:     blocksPool,
		dataAggregator:        dataAggregator,
		outportBlockConverter: outportBlockConverter,
		firstCommitableBlock:  firstCommitableBlock,
	}

	dp.operationHandlers = map[string]func(marshalledData []byte) error{
		outport.TopicSaveBlock: dp.saveBlock,
	}

	return dp, nil
}

// ProcessPayload will process the received payload only for TopicSaveBlock, otherwise ignores it.
func (dp *dataProcessor) ProcessPayload(payload []byte, topic string, _ uint32) error {
	operationHandler, found := dp.operationHandlers[topic]
	if !found {
		return nil
	}

	return operationHandler(payload)
}

func (dp *dataProcessor) saveBlock(marshalledData []byte) error {
	outportBlock := &outport.OutportBlock{}
	err := dp.marshaller.Unmarshal(outportBlock, marshalledData)
	if err != nil {
		return err
	}

	if outportBlock == nil || outportBlock.BlockData == nil {
		return ErrNilOutportBlockData
	}

	if outportBlock.ShardID == core.MetachainShardId {
		return dp.handleMetaOutportBlock(outportBlock)
	}

	return dp.handleShardOutportBlock(outportBlock)
}

func (dp *dataProcessor) handleMetaOutportBlock(outportBlock *outport.OutportBlock) error {
	metaOutportBlock, err := dp.outportBlockConverter.HandleMetaOutportBlock(outportBlock)
	if err != nil {
		return fmt.Errorf("failed to handle metaOutportBlock: %w", err)
	}
	if metaOutportBlock == nil {
		return ErrInvalidOutportBlock
	}
	if metaOutportBlock.BlockData == nil {
		return fmt.Errorf("%w for blockData", ErrInvalidOutportBlock)
	}
	if metaOutportBlock.BlockData.Header == nil {
		return fmt.Errorf("%w for blockData header", ErrInvalidOutportBlock)
	}
	metaNonce := metaOutportBlock.BlockData.Header.GetNonce()

	log.Info("saving meta outport block",
		"hash", metaOutportBlock.BlockData.GetHeaderHash(),
		"round", metaNonce,
		"shardID", metaOutportBlock.ShardID)

	headerHash := outportBlock.BlockData.HeaderHash
	err = dp.outportBlocksPool.PutBlock(headerHash, metaOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to put metablock: %w", err)
	}

	if metaNonce < dp.firstCommitableBlock {
		// do not try to aggregate or publish hyper outport block
		// update only blocks pool state

		log.Trace("do not commit block", "currentRound", metaNonce, "firstCommitableRound", dp.firstCommitableBlock)

		lastCheckpoint := &data.BlockCheckpoint{
			LastNonces: map[uint32]uint64{
				core.MetachainShardId: metaNonce,
			},
		}
		err = dp.outportBlocksPool.UpdateMetaState(lastCheckpoint)
		if err != nil {
			return err
		}

		return nil
	}

	hyperOutportBlock, err := dp.dataAggregator.ProcessHyperBlock(metaOutportBlock)
	if err != nil {
		return err
	}

	lastCheckpoint, err := dp.getLastRoundsData(hyperOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to get last round data: %w", err)
	}

	err = dp.publisher.PublishHyperBlock(hyperOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to publish hyperblock: %w", err)
	}

	return dp.outportBlocksPool.UpdateMetaState(lastCheckpoint)
}

func (dp *dataProcessor) getLastRoundsData(hyperOutportBlock *hyperOutportBlocks.HyperOutportBlock) (*data.BlockCheckpoint, error) {
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

func (dp *dataProcessor) handleShardOutportBlock(outportBlock *outport.OutportBlock) error {
	shardOutportBlock, err := dp.outportBlockConverter.HandleShardOutportBlock(outportBlock)
	if err != nil {
		return fmt.Errorf("failed to handle shardOutportBlock: %w", err)
	}
	if shardOutportBlock.BlockData == nil {
		return fmt.Errorf("%w for blockData", ErrInvalidOutportBlock)
	}
	if shardOutportBlock.BlockData.Header == nil {
		return fmt.Errorf("%w for blockData header", ErrInvalidOutportBlock)
	}
	round := shardOutportBlock.BlockData.Header.GetRound()

	log.Info("saving shard outport block",
		"hash", shardOutportBlock.BlockData.GetHeaderHash(),
		"round", round,
		"shardID", shardOutportBlock.ShardID)

	headerHash := outportBlock.BlockData.HeaderHash

	return dp.outportBlocksPool.PutBlock(headerHash, shardOutportBlock)
}

// Close will close the internal writer
func (dp *dataProcessor) Close() error {
	return dp.publisher.Close()
}

// IsInterfaceNil checks if the underlying pointer is nil
func (dp *dataProcessor) IsInterfaceNil() bool {
	return dp == nil
}

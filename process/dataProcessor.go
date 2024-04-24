package process

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
)

type dataProcessor struct {
	marshaller            marshal.Marshalizer
	operationHandlers     map[string]func(marshalledData []byte) error
	publisher             Publisher
	outportBlocksPool     HyperBlocksPool
	dataAggregator        DataAggregator
	outportBlockConverter OutportBlockConverter
	firstCommitableBlock  uint64
}

// NewDataProcessor creates a data processor able to receive data from a ws outport driver and handle blocks
func NewDataProcessor(
	publisher Publisher,
	marshaller marshal.Marshalizer,
	outportBlocksPool HyperBlocksPool,
	dataAggregator DataAggregator,
	converter OutportBlockConverter,
	firstCommitableBlock uint64,
) (DataProcessor, error) {
	if check.IfNil(publisher) {
		return nil, ErrNilPublisher
	}
	if check.IfNil(outportBlocksPool) {
		return nil, ErrNilHyperBlocksPool
	}
	if check.IfNil(marshaller) {
		return nil, ErrNilMarshaller
	}
	if check.IfNil(dataAggregator) {
		return nil, ErrNilDataAggregator
	}
	if check.IfNil(converter) {
		return nil, ErrNilOutportBlocksConverter
	}

	dp := &dataProcessor{
		marshaller:            marshaller,
		publisher:             publisher,
		outportBlocksPool:     outportBlocksPool,
		dataAggregator:        dataAggregator,
		outportBlockConverter: converter,
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
		return err
	}
	if metaOutportBlock.BlockData == nil {
		return fmt.Errorf("%w for blockData", ErrInvalidOutportBlock)
	}
	if metaOutportBlock.BlockData.Header == nil {
		return fmt.Errorf("%w for blockData header", ErrInvalidOutportBlock)
	}
	metaRound := metaOutportBlock.BlockData.Header.GetRound()

	log.Info("saving meta outport block",
		"hash", metaOutportBlock.BlockData.GetHeaderHash(),
		"round", metaRound,
		"shardID", metaOutportBlock.ShardID)

	headerHash := outportBlock.BlockData.HeaderHash
	err = dp.outportBlocksPool.PutMetaBlock(headerHash, metaOutportBlock)
	if err != nil {
		return err
	}

	if metaRound < dp.firstCommitableBlock {
		// do not try to aggregate or publish hyper outport block
		// update only blocks pool state

		log.Trace("do not commit block", "currentRound", metaRound, "firstCommitableRound", dp.firstCommitableBlock)

		lastCheckpoint := &data.BlockCheckpoint{
			LastRounds: map[uint32]uint64{
				core.MetachainShardId: metaRound,
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
		return err
	}

	err = dp.publisher.PublishHyperBlock(hyperOutportBlock)
	if err != nil {
		return err
	}

	return dp.outportBlocksPool.UpdateMetaState(lastCheckpoint)
}

// TODO: update to use latest data structures
func (dp *dataProcessor) getLastRoundsData(hyperOutportBlock *hyperOutportBlocks.HyperOutportBlock) (*data.BlockCheckpoint, error) {
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

func (dp *dataProcessor) handleShardOutportBlock(outportBlock *outport.OutportBlock) error {
	shardOutportBlock, err := dp.outportBlockConverter.HandleShardOutportBlock(outportBlock)
	if err != nil {
		return err
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

	return dp.outportBlocksPool.PutShardBlock(headerHash, shardOutportBlock)
}

// Close will close the internal writer
func (dp *dataProcessor) Close() error {
	return dp.publisher.Close()
}

// IsInterfaceNil checks if the underlying pointer is nil
func (dp *dataProcessor) IsInterfaceNil() bool {
	return dp == nil
}

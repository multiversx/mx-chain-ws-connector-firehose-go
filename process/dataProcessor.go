package process

import (
	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
)

type dataProcessor struct {
	marshaller        marshal.Marshalizer
	operationHandlers map[string]func(marshalledData []byte) error
	publisher         Publisher
	blocksPool        BlocksPool
	dataAggregator    DataAggregator
	blockCreator      BlockContainerHandler
}

// NewDataProcessor creates a data processor able to receive data from a ws outport driver and handle blocks
func NewDataProcessor(
	publisher Publisher,
	marshaller marshal.Marshalizer,
	blocksPool BlocksPool,
	dataAggregator DataAggregator,
	blockCreator BlockContainerHandler,
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
	if check.IfNil(blockCreator) {
		return nil, ErrNilBlockCreator
	}

	dp := &dataProcessor{
		marshaller:     marshaller,
		publisher:      publisher,
		blocksPool:     blocksPool,
		dataAggregator: dataAggregator,
		blockCreator:   blockCreator,
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

	log.Info("saving block", "hash", outportBlock.BlockData.GetHeaderHash(), "shardID", outportBlock.ShardID)

	if outportBlock.ShardID == core.MetachainShardId {
		return dp.handleMetaOutportBlock(outportBlock)
	}

	return dp.handleShardOutportBlock(outportBlock)
}

func (dp *dataProcessor) handleMetaOutportBlock(outportBlock *outport.OutportBlock) error {
	hyperOutportBlock, err := dp.dataAggregator.ProcessHyperBlock(outportBlock)
	if err != nil {
		return err
	}

	err = dp.publisher.PublishHyperBlock(hyperOutportBlock)
	if err != nil {
		return err
	}

	round, err := dp.getHeaderRound(hyperOutportBlock.MetaOutportBlock)
	if err != nil {
		return err
	}

	dp.blocksPool.UpdateMetaRound(round)

	return nil
}

func (dp *dataProcessor) handleShardOutportBlock(outportBlock *outport.OutportBlock) error {
	blockHash := outportBlock.BlockData.HeaderHash

	round, err := dp.getHeaderRound(outportBlock)
	if err != nil {
		return err
	}

	err = dp.blocksPool.PutBlock(blockHash, outportBlock, round)
	if err != nil {
		return err
	}

	return nil
}

func (dp *dataProcessor) getHeaderRound(outportBlock *outport.OutportBlock) (uint64, error) {
	blockCreator, err := dp.blockCreator.Get(core.HeaderType(outportBlock.BlockData.HeaderType))
	if err != nil {
		return 0, err
	}

	header, err := block.GetHeaderFromBytes(dp.marshaller, blockCreator, outportBlock.BlockData.HeaderBytes)
	if err != nil {
		return 0, err
	}

	return header.GetRound(), nil
}

// Close will close the internal writer
func (dp *dataProcessor) Close() error {
	return dp.publisher.Close()
}

// IsInterfaceNil checks if the underlying pointer is nil
func (dp *dataProcessor) IsInterfaceNil() bool {
	return dp == nil
}

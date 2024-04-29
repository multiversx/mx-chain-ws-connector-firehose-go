package process

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
)

type dataProcessor struct {
	marshaller            marshal.Marshalizer
	operationHandlers     map[string]func(marshalledData []byte) error
	publisher             Publisher
	outportBlocksPool     HyperBlocksPool
	outportBlockConverter OutportBlockConverter
}

// NewDataProcessor creates a data processor able to receive data from a ws outport driver and handle blocks
func NewDataProcessor(
	publisher Publisher,
	marshaller marshal.Marshalizer,
	outportBlocksPool HyperBlocksPool,
	outportBlockConverter OutportBlockConverter,
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
	if check.IfNil(outportBlockConverter) {
		return nil, ErrNilOutportBlocksConverter
	}

	dp := &dataProcessor{
		marshaller:            marshaller,
		publisher:             publisher,
		outportBlocksPool:     outportBlocksPool,
		outportBlockConverter: outportBlockConverter,
	}

	dp.operationHandlers = map[string]func(marshalledData []byte) error{
		outport.TopicSaveBlock:          dp.saveBlock,
		outport.TopicRevertIndexedBlock: dp.revertBlock,
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
		return fmt.Errorf("failed to convert outportBlock to metaOutportBlock: %w", err)
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
	metaRound := metaOutportBlock.BlockData.Header.GetRound()

	log.Info("saving meta outport block",
		"hash", metaOutportBlock.BlockData.GetHeaderHash(),
		"round", metaRound,
		"shardID", metaOutportBlock.ShardID)

	headerHash := metaOutportBlock.BlockData.HeaderHash
	err = dp.outportBlocksPool.PutMetaBlock(headerHash, metaOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to put metablock: %w", err)
	}

	dp.publisher.PublishBlock(headerHash)

	return nil
}

func (dp *dataProcessor) handleShardOutportBlock(outportBlock *outport.OutportBlock) error {
	shardOutportBlock, err := dp.outportBlockConverter.HandleShardOutportBlock(outportBlock)
	if err != nil {
		return fmt.Errorf("failed to convert outportBlock to shardOutportBlock: %w", err)
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

func (dp *dataProcessor) revertBlock(marshalledData []byte) error {
	blockData := &outport.BlockData{}
	err := dp.marshaller.Unmarshal(blockData, marshalledData)
	if err != nil {
		return err
	}

	if blockData == nil {
		return fmt.Errorf("nil block data")
	}

	dp.publisher.PublishBlock(blockData.HeaderHash)

	return nil
}

// Close will close the internal writer
func (dp *dataProcessor) Close() error {
	return dp.publisher.Close()
}

// IsInterfaceNil checks if the underlying pointer is nil
func (dp *dataProcessor) IsInterfaceNil() bool {
	return dp == nil
}

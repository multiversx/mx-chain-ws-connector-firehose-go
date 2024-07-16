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
	outportBlocksPool     BlocksPool
	outportBlockConverter OutportBlockConverter
	firstCommitableBlocks map[uint32]uint64
}

// NewDataProcessor creates a data processor able to receive data from a ws outport driver and handle blocks
func NewDataProcessor(
	publisher Publisher,
	marshaller marshal.Marshalizer,
	blocksPool BlocksPool,
	outportBlockConverter OutportBlockConverter,
	firstCommitableBlocks map[uint32]uint64,
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
	if check.IfNil(outportBlockConverter) {
		return nil, ErrNilOutportBlocksConverter
	}
	if firstCommitableBlocks == nil {
		return nil, ErrNilFirstCommitableBlocks
	}

	dp := &dataProcessor{
		marshaller:            marshaller,
		publisher:             publisher,
		outportBlocksPool:     blocksPool,
		outportBlockConverter: outportBlockConverter,
		firstCommitableBlocks: firstCommitableBlocks,
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
	headerHash := metaOutportBlock.BlockData.GetHeaderHash()

	log.Info("saving meta outport block",
		"hash", headerHash,
		"nonce", metaNonce,
		"shardID", metaOutportBlock.ShardID)

	err = dp.outportBlocksPool.PutBlock(headerHash, metaOutportBlock)
	if err != nil {
		return fmt.Errorf("failed to put metablock: %w", err)
	}

	shardID := metaOutportBlock.GetShardID()
	firstCommitableBlock, ok := dp.firstCommitableBlocks[shardID]
	if !ok {
		return fmt.Errorf("failed to get first commitable block for shard %d", shardID)
	}

	if metaNonce < firstCommitableBlock {
		// do not try to aggregate or publish hyper outport block

		log.Trace("do not publish block", "currentNonce", metaNonce, "firstCommitableNonce", firstCommitableBlock)

		return nil
	}

	err = dp.publisher.PublishBlock(headerHash)
	if err != nil {
		return fmt.Errorf("failed to publish block: %w", err)
	}

	return nil
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
	nonce := shardOutportBlock.BlockData.Header.GetNonce()

	headerHash := shardOutportBlock.BlockData.GetHeaderHash()

	log.Info("saving shard outport block",
		"hash", headerHash,
		"nonce", nonce,
		"shardID", shardOutportBlock.ShardID)

	return dp.outportBlocksPool.PutBlock(headerHash, shardOutportBlock)
}

func (dp *dataProcessor) revertBlock(marshalledData []byte) error {
	blockData := &outport.BlockData{}
	err := dp.marshaller.Unmarshal(blockData, marshalledData)
	if err != nil {
		return err
	}

	err = dp.publisher.PublishBlock(blockData.HeaderHash)
	if err != nil {
		return fmt.Errorf("failed to publish block: %w", err)
	}

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

package process

import (
	"encoding/binary"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
)

type dataProcessor struct {
	marshaller        marshal.Marshalizer
	operationHandlers map[string]func(marshalledData []byte) error
	publisher         Publisher
	outportBlocksPool DataPool
	dataAggregator    DataAggregator
	blockCreator      BlockContainerHandler
}

// NewDataProcessor creates a data processor able to receive data from a ws outport driver and handle blocks
func NewDataProcessor(
	publisher Publisher,
	marshaller marshal.Marshalizer,
	outportBlocksPool DataPool,
	dataAggregator DataAggregator,
	blockCreator BlockContainerHandler,
) (DataProcessor, error) {
	if check.IfNil(publisher) {
		return nil, ErrNilPublisher
	}
	if check.IfNil(outportBlocksPool) {
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
		marshaller:        marshaller,
		publisher:         publisher,
		outportBlocksPool: outportBlocksPool,
		dataAggregator:    dataAggregator,
		blockCreator:      blockCreator,
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

	err = dp.pushToBlocksPool(outportBlock)
	if err != nil {
		return err
	}

	if outportBlock.ShardID == core.MetachainShardId {
		return dp.handleMetaOutportBlock(outportBlock)
	}

	return nil
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

	header, err := dp.getHeader(hyperOutportBlock.MetaOutportBlock)
	if err != nil {
		return err
	}

	err = dp.putMetaNonce(header.GetNonce(), outportBlock.BlockData.GetHeaderHash())
	if err != nil {
		return err
	}

	dp.outportBlocksPool.UpdateMetaState(header.GetRound())

	return nil
}

func (dp *dataProcessor) putMetaNonce(nonce uint64, hash []byte) error {
	nonceBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(nonceBytes, nonce)

	return dp.outportBlocksPool.Put(nonceBytes, hash)
}

func (dp *dataProcessor) pushToBlocksPool(outportBlock *outport.OutportBlock) error {
	blockHash := outportBlock.BlockData.HeaderHash

	header, err := dp.getHeader(outportBlock)
	if err != nil {
		return err
	}

	err = dp.outportBlocksPool.PutBlock(blockHash, outportBlock, header.GetRound())
	if err != nil {
		return err
	}

	return nil
}

func (dp *dataProcessor) getHeader(outportBlock *outport.OutportBlock) (data.HeaderHandler, error) {
	blockCreator, err := dp.blockCreator.Get(core.HeaderType(outportBlock.BlockData.HeaderType))
	if err != nil {
		return nil, err
	}

	header, err := block.GetHeaderFromBytes(dp.marshaller, blockCreator, outportBlock.BlockData.HeaderBytes)
	if err != nil {
		return nil, err
	}

	return header, nil
}

// Close will close the internal writer
func (dp *dataProcessor) Close() error {
	return dp.publisher.Close()
}

// IsInterfaceNil checks if the underlying pointer is nil
func (dp *dataProcessor) IsInterfaceNil() bool {
	return dp == nil
}

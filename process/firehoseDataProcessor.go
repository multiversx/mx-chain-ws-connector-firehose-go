package process

import (
	"encoding/hex"
	"fmt"
	"io"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("firehose")

const (
	firehosePrefix   = "FIRE"
	beginBlockPrefix = "BLOCK_BEGIN"
	endBlockPrefix   = "BLOCK_END"
)

type dataProcessor struct {
	marshaller        marshal.Marshalizer
	operationHandlers map[string]func(marshalledData []byte) error
	writer            io.Writer
	blockCreator      BlockContainerHandler
}

// NewFirehoseDataProcessor creates a data processor able to receive data from a ws outport driver and print saved blocks
func NewFirehoseDataProcessor(
	writer io.Writer,
	blockCreator BlockContainerHandler,
	marshaller marshal.Marshalizer,
) (DataProcessor, error) {
	if writer == nil {
		return nil, errNilWriter
	}
	if check.IfNil(blockCreator) {
		return nil, errNilBlockCreator
	}
	if check.IfNil(marshaller) {
		return nil, errNilMarshaller
	}

	dp := &dataProcessor{
		marshaller:   marshaller,
		writer:       writer,
		blockCreator: blockCreator,
	}

	dp.operationHandlers = map[string]func(marshalledData []byte) error{
		outport.TopicSaveBlock:             dp.saveBlock,
		outport.TopicRevertIndexedBlock:    noOpHandler,
		outport.TopicSaveRoundsInfo:        noOpHandler,
		outport.TopicSaveValidatorsRating:  noOpHandler,
		outport.TopicSaveValidatorsPubKeys: noOpHandler,
		outport.TopicSaveAccounts:          noOpHandler,
		outport.TopicFinalizedBlock:        noOpHandler,
	}

	return dp, nil
}

// ProcessPayload will process the received payload, if the topic is known.
func (dp *dataProcessor) ProcessPayload(payload []byte, topic string) error {
	operationHandler, found := dp.operationHandlers[topic]
	if !found {
		return fmt.Errorf("%w, operation type for topic = %s, received data = %s",
			errInvalidOperationType, topic, string(payload))
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
		return errNilOutportBlockData
	}

	blockCreator, err := dp.blockCreator.Get(core.HeaderType(outportBlock.BlockData.HeaderType))
	if err != nil {
		return err
	}

	header, err := block.GetHeaderFromBytes(dp.marshaller, blockCreator, outportBlock.BlockData.HeaderBytes)
	if err != nil {
		return err
	}

	log.Info("firehose: saving block", "nonce", header.GetNonce(), "hash", outportBlock.BlockData.HeaderHash)

	_, err = fmt.Fprintf(dp.writer, "%s %s %d\n",
		firehosePrefix,
		beginBlockPrefix,
		header.GetNonce(),
	)
	if err != nil {
		return fmt.Errorf("could not write %s prefix , err: %w", beginBlockPrefix, err)
	}

	_, err = fmt.Fprintf(dp.writer, "%s %s %d %s %d %x\n",
		firehosePrefix,
		endBlockPrefix,
		header.GetNonce(),
		hex.EncodeToString(header.GetPrevHash()),
		header.GetTimeStamp(),
		marshalledData,
	)
	if err != nil {
		return fmt.Errorf("could not write %s prefix , err: %w", endBlockPrefix, err)
	}

	return nil
}

func noOpHandler(_ []byte) error {
	return nil
}

// Close will signal via a log that the data processor is closed
func (dp *dataProcessor) Close() error {
	log.Info("data processor closed")
	return nil
}

// IsInterfaceNil checks if the underlying pointer is nil
func (dp *dataProcessor) IsInterfaceNil() bool {
	return dp == nil
}

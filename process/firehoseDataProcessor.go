package process

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("firehose")

const (
	firehosePrefix = "FIRE"
	blockPrefix    = "BLOCK"
)

type dataProcessor struct {
	marshaller        marshal.Marshalizer
	operationHandlers map[string]func(marshalledData []byte) error
	writer            Writer
	blockCreator      BlockContainerHandler
}

// NewFirehoseDataProcessor creates a data processor able to receive data from a ws outport driver and print saved blocks
func NewFirehoseDataProcessor(
	writer Writer,
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

	log.Info("saving block", "nonce", header.GetNonce(), "hash", outportBlock.BlockData.HeaderHash)

	blockNum := header.GetNonce()
	parentNum := blockNum - 1
	if blockNum == 0 {
		parentNum = 0
	}
	encodedMarshalledData := base64.StdEncoding.EncodeToString(marshalledData)

	_, err = fmt.Fprintf(dp.writer, "%s %s %d %s %d %s %d %d %s\n",
		firehosePrefix,
		blockPrefix,
		blockNum,
		hex.EncodeToString(outportBlock.BlockData.HeaderHash),
		parentNum,
		hex.EncodeToString(header.GetPrevHash()),
		outportBlock.HighestFinalBlockNonce,
		header.GetTimeStamp(),
		encodedMarshalledData,
	)

	return err
}

// Close will close the internal writer
func (dp *dataProcessor) Close() error {
	return dp.writer.Close()
}

// IsInterfaceNil checks if the underlying pointer is nil
func (dp *dataProcessor) IsInterfaceNil() bool {
	return dp == nil
}

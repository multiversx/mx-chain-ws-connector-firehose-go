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
	"github.com/multiversx/mx-chain-ws-connector-template-go/pbmultiversx"
)

var log = logger.GetOrCreate("firehose")

const (
	firehosePrefix = "FIRE"
	blockPrefix    = "BLOCK"
	initPrefix     = "INIT"

	protocolReaderVersion = "1.0"
	protoMessageType      = "type.googleapis.com/sf.multiversx.type.v1.Block"
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

	_, _ = fmt.Fprintf(dp.writer, "%s %s %s %s\n", firehosePrefix, initPrefix, protocolReaderVersion, protoMessageType)

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

	// TODO: update to use aggregated data structure
	tmpBlock := &pbmultiversx.Block{}

	blockHeader := &pbmultiversx.BlockHeader{
		Height:       blockNum,
		Hash:         hex.EncodeToString(outportBlock.BlockData.HeaderHash),
		PreviousNum:  parentNum,
		PreviousHash: hex.EncodeToString(header.GetPrevHash()),
		FinalNum:     outportBlock.HighestFinalBlockNonce,
		FinalHash:    hex.EncodeToString(outportBlock.HighestFinalBlockHash),
		Timestamp:    header.GetTimeStamp(),
	}

	tmpBlock.MultiversxBlock = outportBlock
	tmpBlock.Header = blockHeader

	tmpBlockMarshalled, err := dp.marshaller.Marshal(tmpBlock)
	if err != nil {
		return err
	}

	encodedMarshalledData := base64.StdEncoding.EncodeToString(tmpBlockMarshalled)

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

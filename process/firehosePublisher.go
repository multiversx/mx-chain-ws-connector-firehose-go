package process

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
)

var log = logger.GetOrCreate("process")

const (
	firehosePrefix = "FIRE"
	blockPrefix    = "BLOCK"
	initPrefix     = "INIT"

	protocolReaderVersion = "1.0"
	protoMessageType      = "type.googleapis.com/proto.OutportBlock"
)

type firehosePublisher struct {
	marshaller   marshal.Marshalizer
	writer       Writer
	blockCreator BlockContainerHandler
}

// NewFirehosePublisher creates a data processor able to receive data from a ws outport driver and print saved blocks
func NewFirehosePublisher(
	writer Writer,
	blockCreator BlockContainerHandler,
	marshaller marshal.Marshalizer,
) (*firehosePublisher, error) {
	if writer == nil {
		return nil, ErrNilWriter
	}
	if check.IfNil(blockCreator) {
		return nil, ErrNilBlockCreator
	}
	if check.IfNil(marshaller) {
		return nil, ErrNilMarshaller
	}

	fp := &firehosePublisher{
		marshaller:   marshaller,
		writer:       writer,
		blockCreator: blockCreator,
	}

	_, err := fmt.Fprintf(fp.writer, "%s %s %s %s\n", firehosePrefix, initPrefix, protocolReaderVersion, protoMessageType)
	if err != nil {
		return nil, err
	}

	return fp, nil
}

// PublishHyperBlock will push aggregated outport block data to the firehose writer
func (fp *firehosePublisher) PublishHyperBlock(hyperOutportBlock *data.HyperOutportBlock) error {
	outportBlock := hyperOutportBlock.MetaOutportBlock

	blockCreator, err := fp.blockCreator.Get(core.HeaderType(outportBlock.BlockData.HeaderType))
	if err != nil {
		return err
	}

	header, err := block.GetHeaderFromBytes(fp.marshaller, blockCreator, outportBlock.BlockData.HeaderBytes)
	if err != nil {
		return err
	}

	blockNum := header.GetNonce()
	parentNum := blockNum - 1
	if blockNum == 0 {
		parentNum = 0
	}

	marshalledData, err := fp.marshaller.Marshal(hyperOutportBlock)
	if err != nil {
		return err
	}

	encodedMarshalledData := base64.StdEncoding.EncodeToString(marshalledData)

	_, err = fmt.Fprintf(fp.writer, "%s %s %d %s %d %s %d %d %s\n",
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
func (fp *firehosePublisher) Close() error {
	return fp.writer.Close()
}

// IsInterfaceNil checks if the underlying pointer is nil
func (fp *firehosePublisher) IsInterfaceNil() bool {
	return fp == nil
}

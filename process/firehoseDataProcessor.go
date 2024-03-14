package process

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/multiversx/mx-chain-core-go/data"
	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("firehose")

const (
	firehosePrefix = "FIRE"
	blockPrefix    = "BLOCK"
	initPrefix     = "INIT"

	protocolReaderVersion = "1.0"
	protoMessageType      = "type.googleapis.com/bstream.pb.sf.bstream.v1.OutportBlock"
)

type firehosePublisher struct {
	writer Writer
}

// NewFirehosePublisher creates a data processor able to receive data from a ws outport driver and print saved blocks
func NewFirehosePublisher(
	writer Writer,
) (*firehosePublisher, error) {
	if writer == nil {
		return nil, errNilWriter
	}

	fp := &firehosePublisher{
		writer: writer,
	}

	_, err := fmt.Fprintf(fp.writer, "%s %s %s %s\n", firehosePrefix, initPrefix, protocolReaderVersion, protoMessageType)
	if err != nil {
		return nil, err
	}

	return fp, nil
}

func (fp *firehosePublisher) PublishHyperBlock(header data.HeaderHandler, headerHash []byte, marshalledData []byte) error {
	log.Info("saving block", "nonce", header.GetNonce(), "hash", headerHash)

	blockNum := header.GetNonce()
	parentNum := blockNum - 1
	if blockNum == 0 {
		parentNum = 0
	}

	encodedMarshalledData := base64.StdEncoding.EncodeToString(marshalledData)

	_, err := fmt.Fprintf(fp.writer, "%s %s %d %s %d %s %d %d %s\n",
		firehosePrefix,
		blockPrefix,
		blockNum,
		hex.EncodeToString(headerHash),
		parentNum,
		hex.EncodeToString(header.GetPrevHash()),
		// outportBlock.HighestFinalBlockNonce,
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

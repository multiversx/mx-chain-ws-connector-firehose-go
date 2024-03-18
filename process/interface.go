package process

import (
	"io"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-ws-connector-template-go/data"
)

// WSConnector defines a ws connector that receives incoming data and can be closed
type WSConnector interface {
	Close() error
}

// DataProcessor defines a payload processor for incoming ws data
type DataProcessor interface {
	ProcessPayload(payload []byte, topic string, version uint32) error
	Close() error
	IsInterfaceNil() bool
}

// Logger defines the behavior of a data logger component
type Logger interface {
	Info(message string, args ...interface{})
	IsInterfaceNil() bool
}

// BlockContainerHandler defines a block creator container
type BlockContainerHandler interface {
	Add(headerType core.HeaderType, creator block.EmptyBlockCreator) error
	Get(headerType core.HeaderType) (block.EmptyBlockCreator, error)
	IsInterfaceNil() bool
}

// Writer defines a handler for the Write method
type Writer interface {
	io.Writer
	Close() error
}

type Publisher interface {
	PublishHyperBlock(hyperOutportBlock *data.HyperOutportBlock) error
	Close() error
}

type DataAggregator interface {
	ProcessHyperBlock(outportBlock *outport.OutportBlock) (*data.HyperOutportBlock, error)
	IsInterfaceNil() bool
}

type BlocksPool interface {
	PutBlock(hash []byte, outportBlock *outport.OutportBlock) error
	GetBlock(hash []byte) (*outport.OutportBlock, error)
	UpdateMetaRound(round uint64)
	IsInterfaceNil() bool
}

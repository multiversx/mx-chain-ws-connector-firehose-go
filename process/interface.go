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

// Publisher defines the behaviour of an aggregated outport block publisher component
type Publisher interface {
	PublishHyperBlock(hyperOutportBlock *data.HyperOutportBlock) error
	Close() error
	IsInterfaceNil() bool
}

// DataAggregator defines the behaviour of a component that is able to aggregate outport
// block data for shards
type DataAggregator interface {
	ProcessHyperBlock(outportBlock *outport.OutportBlock) (*data.HyperOutportBlock, error)
	IsInterfaceNil() bool
}

// DataPool defines the behaviour of a data pool handler component
type DataPool interface {
	PutBlock(hash []byte, data []byte, index uint64, shardID uint32) error
	GetBlock(hash []byte) ([]byte, error)
	UpdateMetaState(index uint64)
	IsInterfaceNil() bool
}

// OutportBlocksPool defines the behaviour of an outport blocks pool handler component
type OutportBlocksPool interface {
	PutBlock(hash []byte, outportBlock *outport.OutportBlock, round uint64) error
	GetBlock(hash []byte) (*outport.OutportBlock, error)
	UpdateMetaState(round uint64)
	IsInterfaceNil() bool
}

// HyperOutportBlocksPool defines the behaviour of a hyper blocks pool handler component
type HyperOutportBlocksPool interface {
	PutBlock(hash []byte, outportBlock *data.HyperOutportBlock, round uint64) error
	GetBlock(hash []byte) (*data.HyperOutportBlock, error)
	UpdateMetaState(round uint64)
	IsInterfaceNil() bool
}

// PruningStorer defines the behaviour of a pruning storer component
type PruningStorer interface {
	Get(key []byte) ([]byte, error)
	Put(key, data []byte) error
	Prune(index uint64) error
	Dump() error
	SetCheckpoint(round uint64) error
	Close() error
	IsInterfaceNil() bool
}

type Server interface {
	Start() error
}

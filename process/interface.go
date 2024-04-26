package process

import (
	"io"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/data/block"
	"github.com/multiversx/mx-chain-core-go/data/outport"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
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
	PublishHyperBlock(hyperOutportBlock *hyperOutportBlocks.HyperOutportBlock) error
	Close() error
	IsInterfaceNil() bool
}

// DataAggregator defines the behaviour of a component that is able to aggregate outport
// block data for shards
type DataAggregator interface {
	ProcessHyperBlock(outportBlock *hyperOutportBlocks.MetaOutportBlock) (*hyperOutportBlocks.HyperOutportBlock, error)
	IsInterfaceNil() bool
}

// DataPool defines the behaviour of a data pool handler component
type DataPool interface {
	Put(key []byte, value []byte) error
	Get(key []byte) ([]byte, error)
	PutBlock(hash []byte, value []byte, round uint64, shardID uint32) error
	UpdateMetaState(checkpoint *data.BlockCheckpoint) error
	Close() error
	IsInterfaceNil() bool
}

// HyperBlocksPool defines the behaviour of a blocks pool handler component that is
// able to handle meta and shard outport blocks data
type HyperBlocksPool interface {
	Get(key []byte) ([]byte, error)
	PutMetaBlock(hash []byte, outportBlock *hyperOutportBlocks.MetaOutportBlock) error
	PutShardBlock(hash []byte, outportBlock *hyperOutportBlocks.ShardOutportBlock) error
	GetMetaBlock(hash []byte) (*hyperOutportBlocks.MetaOutportBlock, error)
	GetShardBlock(hash []byte) (*hyperOutportBlocks.ShardOutportBlock, error)
	UpdateMetaState(checkpoint *data.BlockCheckpoint) error
	Close() error
	IsInterfaceNil() bool
}

// PruningStorer defines the behaviour of a pruning storer component
type PruningStorer interface {
	Get(key []byte) ([]byte, error)
	Put(key, data []byte) error
	Prune(index uint64) error
	Dump() error
	Close() error
	Destroy() error
	IsInterfaceNil() bool
}

// OutportBlockConverter handles the conversion between gogo and google proto buffer definitions.
type OutportBlockConverter interface {
	HandleShardOutportBlock(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.ShardOutportBlock, error)
	HandleMetaOutportBlock(outportBlock *outport.OutportBlock) (*hyperOutportBlocks.MetaOutportBlock, error)
	IsInterfaceNil() bool
}

// GRPCBlocksHandler defines the behaviour of handling block via gRPC
type GRPCBlocksHandler interface {
	FetchHyperBlockByHash(hash []byte) (*hyperOutportBlocks.HyperOutportBlock, error)
	FetchHyperBlockByNonce(nonce uint64) (*hyperOutportBlocks.HyperOutportBlock, error)
	IsInterfaceNil() bool
}

// GRPCServer is the server that will serve the stored hyperOutportBlocks.
type GRPCServer interface {
	Start()
	Close()
}

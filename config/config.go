package config

// Config holds general configuration
type Config struct {
	WebSocket            WebSocketConfig
	DataPool             DataPoolConfig
	OutportBlocksStorage StorageConfig
	GRPC                 GRPCConfig
}

// WebSocketConfig holds web sockets config
type WebSocketConfig struct {
	URL                        string
	MarshallerType             string
	Mode                       string
	RetryDurationInSec         uint32
	WithAcknowledge            bool
	AcknowledgeTimeoutInSec    int
	BlockingAckOnError         bool
	DropMessagesIfNoConnection bool
	Version                    uint32
}

// DataPoolConfig will map data poil configuration
type DataPoolConfig struct {
	MaxDelta             uint64
	PruningWindow        uint64
	NumPersistersToKeep  int
	FirstCommitableBlock uint64
}

// StorageConfig will map the storage unit configuration
type StorageConfig struct {
	Cache CacheConfig
	DB    DBConfig
}

// CacheConfig will map the cache configuration
type CacheConfig struct {
	Name        string
	Type        string
	Capacity    uint32
	SizeInBytes uint64
}

// DBConfig will map the database configuration
type DBConfig struct {
	FilePath          string
	Type              string
	BatchDelaySeconds int
	MaxBatchSize      int
	MaxOpenFiles      int
}

// GRPCConfig will map the gRPC server configuration
type GRPCConfig struct {
	URL string
}

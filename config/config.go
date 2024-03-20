package config

// Config holds general configuration
type Config struct {
	WebSocketConfig      WebSocketConfig `toml:"web_socket"`
	DataPoolConfig       DataPoolConfig
	OutportBlocksStorage StorageConfig
}

// WebSocketConfig holds web sockets config
type WebSocketConfig struct {
	Url                        string `toml:"url"`
	MarshallerType             string `toml:"marshaller_type"`
	Mode                       string `toml:"mode"`
	RetryDuration              uint32 `toml:"retry_duration"`
	WithAcknowledge            bool   `toml:"with_acknowledge"`
	AcknowledgeTimeoutInSec    int    `toml:"acknowledge_timeout_in_sec"`
	BlockingAckOnError         bool   `toml:"blocking_ack_on_error"`
	DropMessagesIfNoConnection bool   `toml:"drop_messages_if_no_connection"` // Set to `true` to drop messages if there is no active WebSocket connection to send to.
	Version                    uint32 `toml:"version"`                        // Defines the payload version.
}

// DataPoolConfig will map data poil configuration
type DataPoolConfig struct {
	NumberOfShards uint32
	MaxDelta       uint64
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

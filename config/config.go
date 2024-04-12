package config

// Config holds general configuration
type Config struct {
	WebSocketConfig WebSocketConfig `toml:"web_socket"`
	GRPCConfig      GRPCConfig      `toml:"grpc_server"`
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

// GRPCConfig holds gRPC server config
type GRPCConfig struct {
	URL string `toml:"url"`
}

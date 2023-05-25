package config

// Config holds general configuration
type Config struct {
	WebSocketConfig WebSocketConfig `toml:"web_socket"`
}

// WebSocketConfig holds web sockets config
type WebSocketConfig struct {
	Url                string `toml:"url"`
	MarshallerType     string `toml:"marshaller_type"`
	Mode               string `toml:"mode"`
	RetryDuration      uint32 `toml:"retry_duration"`
	WithAcknowledge    bool   `toml:"with_acknowledge"`
	BlockingAckOnError bool   `toml:"blocking_ack_on_error"`
}

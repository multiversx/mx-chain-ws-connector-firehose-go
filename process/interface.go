package process

// WSConnector defines a ws connector that receives incoming data and can be closed
type WSConnector interface {
	Close() error
}

// DataProcessor defines a payload processor for incoming ws data
type DataProcessor interface {
	ProcessPayload(payload []byte, topic string) error
	Close() error
	IsInterfaceNil() bool
}

// Logger defines the behavior of a data logger component
type Logger interface {
	Info(message string, args ...interface{})
	IsInterfaceNil() bool
}

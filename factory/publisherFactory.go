package factory

import (
	"fmt"
	"os"

	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
)

// CreatePublisher will return the required Publisher implementation based on whether the hyperOutportBlock are
// served via gRPC or stdout.
func CreatePublisher(enableGRPCServer bool, blockContainer process.BlockContainerHandler) (process.HyperBlockPublisher, error) {
	if enableGRPCServer {
		return process.NewGRPCBlockPublisher(), nil
	}

	publisher, err := process.NewFirehosePublisher(
		os.Stdout,
		blockContainer,
		&process.ProtoMarshaller{},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create firehose publisher: %w", err)
	}

	return publisher, nil
}

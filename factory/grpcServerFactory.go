package factory

import (
	"fmt"
	"os"
	"time"

	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/server"
)

// CreateGRPCServer will create the gRPC server along with the required handlers.
func CreateGRPCServer(
	enableGrpcServer bool,
	cfg config.GRPCConfig,
	outportBlocksPool process.DataPool,
	dataAggregator process.DataAggregator) (process.GRPCServer, process.Writer, error) {
	if !enableGrpcServer {
		return nil, os.Stdout, nil
	}

	handler, err := process.NewGRPCBlocksHandler(outportBlocksPool, dataAggregator)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create grpc blocks handler: %w", err)
	}
	s := server.New(cfg, handler)
	err = s.Start()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start grpc server: %w", err)
	}

	return s, &fakeWriter{}, nil

}

type fakeWriter struct {
	err      error
	duration time.Duration
}

// Write is a mock writer.
func (f *fakeWriter) Write(p []byte) (int, error) {
	time.Sleep(f.duration)
	if f.err != nil {
		return 0, f.err
	}

	return len(p), nil
}

// Close is mock closer.
func (f *fakeWriter) Close() error {
	return nil
}

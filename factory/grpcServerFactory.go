package factory

import (
	"fmt"

	"github.com/multiversx/mx-chain-ws-connector-firehose-go/config"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/server"
)

// CreateGRPCServer will create the gRPC server along with the required handlers.
func CreateGRPCServer(
	enableGrpcServer bool,
	cfg config.GRPCConfig,
	outportBlocksPool process.HyperBlocksPool,
	dataAggregator process.DataAggregator) (process.GRPCServer, error) {
	if !enableGrpcServer {
		return nil, nil
	}

	handler, err := process.NewGRPCBlocksHandler(outportBlocksPool, dataAggregator)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc blocks handler: %w", err)
	}
	s, err := server.New(cfg, handler)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc server: %w", err)
	}
	s.Start()

	return s, nil

}

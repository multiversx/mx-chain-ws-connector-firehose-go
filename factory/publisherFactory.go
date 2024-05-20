package factory

import (
	"fmt"
	"os"

	"github.com/multiversx/mx-chain-ws-connector-firehose-go/config"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/server"
	"google.golang.org/grpc"
)

// CreateHyperBlockPublisher will create a new hyper block publisher component
func CreateHyperBlockPublisher(
	cfg *config.Config,
	enableGRPCServer bool,
	blockContainer process.BlockContainerHandler,
	outportBlocksPool process.BlocksPool,
	dataAggregator process.DataAggregator,
) (process.HyperBlockPublisher, error) {
	if enableGRPCServer {
		handler, err := process.NewGRPCBlocksHandler(outportBlocksPool, dataAggregator)
		if err != nil {
			return nil, fmt.Errorf("failed to create grpc blocks handler: %w", err)
		}

		grpcServer := grpc.NewServer()

		s, err := server.NewGRPCServerWrapper(grpcServer, cfg.GRPC, handler)
		if err != nil {
			return nil, fmt.Errorf("failed to create grpc server: %w", err)
		}

		return process.NewGRPCBlockPublisher(s)
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

package server

import (
	"context"
	"fmt"
	"net"

	logger "github.com/multiversx/mx-chain-logger-go"
	"google.golang.org/grpc/reflection"

	"github.com/multiversx/mx-chain-ws-connector-firehose-go/config"
	data "github.com/multiversx/mx-chain-ws-connector-firehose-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/process"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/service/hyperOutportBlock"
)

var (
	log = logger.GetOrCreate("server")
)

type grpcServer interface {
	reflection.GRPCServer
	Serve(lis net.Listener) error
	GracefulStop()
}

type grpcServerWrapper struct {
	server grpcServer
	config config.GRPCConfig

	cancelFunc context.CancelFunc
}

// NewGRPCServerWrapper instantiates the underlying grpc server handling rpc requests.
func NewGRPCServerWrapper(
	grpcServer grpcServer,
	config config.GRPCConfig,
	blocksHandler process.GRPCBlocksHandler,
) (*grpcServerWrapper, error) {
	ctx, cancelFunc := context.WithCancel(context.Background())
	service, err := hyperOutportBlock.NewService(ctx, blocksHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}
	data.RegisterBlockStreamServer(grpcServer, service)
	reflection.Register(grpcServer)

	return &grpcServerWrapper{
		server:     grpcServer,
		config:     config,
		cancelFunc: cancelFunc,
	}, nil
}

// Start will start the grpc server on the configured URL.
func (s *grpcServerWrapper) Start() {
	go func() {
		err := s.run()
		if err != nil {
			log.Error("failed to start grpc server", "error", err)
		}
	}()
}

func (s *grpcServerWrapper) run() error {
	lis, err := net.Listen("tcp", s.config.URL)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	if err = s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// Close will gracefully stop the grpc server.
func (s *grpcServerWrapper) Close() {
	if s.cancelFunc != nil {
		s.cancelFunc()
	}

	s.server.GracefulStop()
}

// IsInterfaceNil checks if the underlying server is nil.
func (s *grpcServerWrapper) IsInterfaceNil() bool {
	return s == nil
}

package factory

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/service/hyperOutportBlock"
)

type grpcServer struct {
	server *grpc.Server
	config config.GRPCConfig
}

// NewServer instantiates the underlying grpc server handling rpc requests.
func NewServer(config config.GRPCConfig, blocksHandler process.GrpcBlocksHandler) *grpcServer {
	s := grpc.NewServer()

	service := hyperOutportBlock.NewService(blocksHandler)
	data.RegisterHyperOutportBlockServiceServer(s, service)
	reflection.Register(s)

	return &grpcServer{s, config}
}

// Start will start the grpc server on the configured URL.
func (s *grpcServer) Start() error {
	lis, err := net.Listen("tcp", s.config.URL)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	if err = s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to serve: %v", err)
	}

	return nil
}

// Stop will gracefully stop the grpc server.
func (s *grpcServer) Stop() {
	s.server.GracefulStop()
}

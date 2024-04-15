package factory

import (
	"fmt"
	"net"

	logger "github.com/multiversx/mx-chain-logger-go"
	"google.golang.org/grpc"

	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/service/hyperOutportBlock"
)

var serverLog = logger.GetOrCreate("service")

type grpcServer struct {
	server *grpc.Server
	config config.GRPCConfig
}

func NewServer(config config.GRPCConfig, pool process.HyperOutportBlocksPool) *grpcServer {
	s := grpc.NewServer()

	hyperService := hyperOutportBlock.Service{BlocksPool: pool}
	hyperService.Register(s)

	return &grpcServer{s, config}
}

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

func (s *grpcServer) Stop() {
	s.server.GracefulStop()
}

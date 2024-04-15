package factory

import (
	"fmt"
	"net"

	logger "github.com/multiversx/mx-chain-logger-go"
	"google.golang.org/grpc"

	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	"github.com/multiversx/mx-chain-ws-connector-template-go/service/hyperOutportBlock"
)

var serverLog = logger.GetOrCreate("service")

type grpcServer struct {
	server *grpc.Server
	config config.GRPCConfig
}

func NewServer(config config.GRPCConfig) *grpcServer {
	s := grpc.NewServer()

	hyperService := hyperOutportBlock.Service{}
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

// CreateGRPCServer will create a gRPC service able to process incoming request.
func CreateGRPCServer(cfg config.GRPCConfig) (*grpcServer, error) {
	server := NewServer(cfg)

	go func() {
		if err := server.Start(); err != nil {
			serverLog.Error("failed to start service", "error", err)
			return
		}
	}()

	return server, nil
}

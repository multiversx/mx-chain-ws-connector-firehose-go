package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	logger "github.com/multiversx/mx-chain-logger-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/multiversx/mx-chain-ws-connector-template-go/config"
	data "github.com/multiversx/mx-chain-ws-connector-template-go/data/hyperOutportBlocks"
	"github.com/multiversx/mx-chain-ws-connector-template-go/process"
	"github.com/multiversx/mx-chain-ws-connector-template-go/service/hyperOutportBlock"
)

var (
	log              = logger.GetOrCreate("server")
	errServerStarted = errors.New("server is already started")
)

type grpcServer struct {
	server *grpc.Server
	config config.GRPCConfig

	cancelFunc func()
	closeChan  chan struct{}
	mutState   sync.RWMutex
}

// New instantiates the underlying grpc server handling rpc requests.
func New(config config.GRPCConfig, blocksHandler process.GRPCBlocksHandler) *grpcServer {
	s := grpc.NewServer()

	service := hyperOutportBlock.NewService(blocksHandler)
	data.RegisterHyperOutportBlockServiceServer(s, service)
	reflection.Register(s)

	return &grpcServer{
		server:    s,
		config:    config,
		closeChan: make(chan struct{}),
	}
}

// Start will start the grpc server on the configured URL.
func (s *grpcServer) Start() error {
	s.mutState.Lock()
	defer s.mutState.Unlock()

	if s.cancelFunc != nil {
		return errServerStarted
	}

	var (
		ctx context.Context
		err error
	)
	ctx, s.cancelFunc = context.WithCancel(context.Background())

	go func() {
		err = s.run(ctx)
		if err != nil {
			log.Error("failed to serve server", "err", err)
			return
		}
	}()

	return err
}

func (s *grpcServer) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			s.server.GracefulStop()

		default:
			lis, err := net.Listen("tcp", s.config.URL)
			if err != nil {
				return fmt.Errorf("failed to listen: %v", err)
			}

			if err = s.server.Serve(lis); err != nil {
				return fmt.Errorf("failed to serve: %v", err)
			}
		}
	}
}

// Close will gracefully stop the grpc server.
func (s *grpcServer) Close() {
	s.mutState.RLock()
	defer s.mutState.RUnlock()

	if s.cancelFunc != nil {
		s.cancelFunc()
	}

	close(s.closeChan)
}

package server

import (
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/multiversx/mx-chain-ws-connector-template-go/server/dummy"
)

type server struct {
	grpcServer *grpc.Server
}

// New returns the server instance with the services registered.
func New() *server {
	s := grpc.NewServer()
	service := dummy.BlockService{}
	service.Register(s)

	return &server{s}
}

// Start serving the server on the listener.
func (s *server) Start(lis net.Listener) {
	if err := s.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

package server_test

import (
	"net"
	"testing"

	"github.com/multiversx/mx-chain-ws-connector-firehose-go/config"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/server"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/testscommon"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

type GRPCServerStub struct{}

func (g *GRPCServerStub) RegisterService(desc *grpc.ServiceDesc, impl any) {
}

func (g *GRPCServerStub) GetServiceInfo() map[string]grpc.ServiceInfo {
	return make(map[string]grpc.ServiceInfo)
}

func (g *GRPCServerStub) Serve(lis net.Listener) error {
	return nil
}

func (g *GRPCServerStub) GracefulStop() {
}

func TestNewGRPCServerWrapper(t *testing.T) {
	t.Parallel()

	gsv, err := server.NewGRPCServerWrapper(
		&GRPCServerStub{},
		config.GRPCConfig{
			URL: "localhost:8081",
		},
		&testscommon.GRPCBlocksHandlerStub{},
	)
	require.Nil(t, err)
	require.False(t, gsv.IsInterfaceNil())

	gsv.Start()
	gsv.Close()
}

package server_test

import (
	"net"
	"sync/atomic"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-ws-connector-firehose-go/config"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/server"
	"github.com/multiversx/mx-chain-ws-connector-firehose-go/testscommon"
	"github.com/stretchr/testify/require"
)

func TestNewGRPCServerWrapper(t *testing.T) {
	t.Parallel()

	serveCalls := uint32(0)
	stopCalled := false
	gsv, err := server.NewGRPCServerWrapper(
		&testscommon.GRPCServerMock{
			ServeCalled: func(lis net.Listener) error {
				atomic.AddUint32(&serveCalls, 1)
				return nil
			},
			GracefulStopCalled: func() {
				stopCalled = true
			},
		},
		config.GRPCConfig{
			URL: ":8081",
		},
		&testscommon.GRPCBlocksHandlerMock{},
	)
	require.Nil(t, err)
	require.False(t, gsv.IsInterfaceNil())

	gsv.Start()

	time.Sleep(1000 * time.Millisecond)

	require.Equal(t, uint32(1), atomic.LoadUint32(&serveCalls))

	gsv.Close()

	require.True(t, stopCalled)
}

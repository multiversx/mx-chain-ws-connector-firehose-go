package testscommon

import (
	"net"

	"google.golang.org/grpc"
)

// GRPCServerMock -
type GRPCServerMock struct {
	ServeCalled        func(lis net.Listener) error
	GracefulStopCalled func()
}

// RegisterService -
func (g *GRPCServerMock) RegisterService(desc *grpc.ServiceDesc, impl any) {
}

// GetServiceInfo -
func (g *GRPCServerMock) GetServiceInfo() map[string]grpc.ServiceInfo {
	return make(map[string]grpc.ServiceInfo)
}

// Serve -
func (g *GRPCServerMock) Serve(lis net.Listener) error {
	if g.ServeCalled != nil {
		return g.ServeCalled(lis)
	}

	return nil
}

// GracefulStop -
func (g *GRPCServerMock) GracefulStop() {
	if g.GracefulStopCalled != nil {
		g.GracefulStopCalled()
	}
}

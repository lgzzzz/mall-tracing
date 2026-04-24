package grpcutil

import (
	"context"

	"github.com/go-kratos/kratos/v2/registry"
	kgrpc "github.com/go-kratos/kratos/v2/transport/grpc"
	"google.golang.org/grpc"
)

// NewInsecureClient creates a gRPC client connection with service discovery.
func NewInsecureClient(
	ctx context.Context,
	discovery registry.Discovery,
	endpoint string,
	opts ...kgrpc.ClientOption,
) (*grpc.ClientConn, error) {
	dialOpts := []kgrpc.ClientOption{
		kgrpc.WithEndpoint(endpoint),
		kgrpc.WithDiscovery(discovery),
	}
	dialOpts = append(dialOpts, opts...)

	return kgrpc.DialInsecure(ctx, dialOpts...)
}

// NewDirectClient creates a gRPC client connection to a direct endpoint (no discovery).
func NewDirectClient(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
	return kgrpc.DialInsecure(ctx, kgrpc.WithEndpoint(endpoint))
}

package grpcutil

import (
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

// ServerBuilder builds a gRPC server with fluent API.
type ServerBuilder struct {
	address    string
	timeout    time.Duration
	middleware []middleware.Middleware
	registers  []func(*grpc.Server)
}

// NewServerBuilder creates a new ServerBuilder.
func NewServerBuilder() *ServerBuilder {
	return &ServerBuilder{}
}

// WithAddress sets the server bind address.
func (b *ServerBuilder) WithAddress(addr string) *ServerBuilder {
	b.address = addr
	return b
}

// WithTimeout sets the request timeout.
func (b *ServerBuilder) WithTimeout(timeout time.Duration) *ServerBuilder {
	b.timeout = timeout
	return b
}

// WithMiddleware adds middleware to the server.
func (b *ServerBuilder) WithMiddleware(mw ...middleware.Middleware) *ServerBuilder {
	b.middleware = append(b.middleware, mw...)
	return b
}

// RegisterService registers a gRPC service.
func (b *ServerBuilder) RegisterService(register func(*grpc.Server)) *ServerBuilder {
	b.registers = append(b.registers, register)
	return b
}

// Build creates the gRPC server.
func (b *ServerBuilder) Build() *grpc.Server {
	var opts []grpc.ServerOption

	if b.address != "" {
		opts = append(opts, grpc.Address(b.address))
	}
	if b.timeout > 0 {
		opts = append(opts, grpc.Timeout(b.timeout))
	}
	if len(b.middleware) > 0 {
		opts = append(opts, grpc.Middleware(b.middleware...))
	}

	srv := grpc.NewServer(opts...)
	for _, register := range b.registers {
		register(srv)
	}
	return srv
}

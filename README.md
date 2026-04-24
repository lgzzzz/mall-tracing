# mall-tracing

Shared tracing and utility library for mall-kratos microservices.

## Features

- **OpenTelemetry Tracing**: Initialize TracerProvider with OTLP exporter to Jaeger
- **Shared Middleware**: JWT auth, response error handling, tracing middleware
- **gRPC Utilities**: Client creation with service discovery, server builder
- **Data Layer**: GORM setup with tracing plugin, etcd discovery
- **Kafka Tracing**: Producer/consumer wrappers with context propagation

## Usage

### Tracing Initialization

```go
tp, err := tracing.Init(tracing.Config{
    ServiceName:  "order-service",
    Version:      "v1.0.0",
    OTLPEndpoint: "localhost:4317",
    SampleRatio:  0.1,
    Insecure:     true,
})
if err != nil {
    log.Fatal(err)
}
defer tracing.Shutdown(context.Background(), tp)

tracer := tracing.NewTracer("order-service")
```

### Server Middleware

```go
grpcutil.NewServerBuilder().
    WithMiddleware(
        recovery.Recovery(),
        mallmiddleware.ServerMiddleware(tracer),
        mallmiddleware.ResponseError(newStatus),
        mallmiddleware.ServerAuth(secret),
    ).
    Build()
```

### Client Middleware

```go
conn, err := grpcutil.NewInsecureClient(ctx, discovery, endpoint,
    grpc.WithMiddleware(mallmiddleware.ClientMiddleware(tracer)),
)
```

## License

MIT

package tracing

import (
	_ "github.com/go-kratos/kratos/v2"
	_ "github.com/golang-jwt/jwt/v5"
	_ "github.com/segmentio/kafka-go"
	_ "go.etcd.io/etcd/client/v3"
	_ "go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	_ "go.opentelemetry.io/otel/sdk"
	_ "go.opentelemetry.io/otel/trace"
	_ "google.golang.org/grpc"
	_ "gorm.io/gorm"
)

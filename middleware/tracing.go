package middleware

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// ServerMiddleware returns a server-side tracing middleware.
func ServerMiddleware(tracer trace.Tracer) middleware.Middleware {
	if tracer == nil {
		tracer = otel.Tracer("github.com/lgzzz/mall-tracing")
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var spanName string
			if tr, ok := transport.FromServerContext(ctx); ok {
				spanName = tr.Operation()
			} else {
				spanName = "grpc.request"
			}

			ctx, span := tracer.Start(ctx, spanName,
				trace.WithSpanKind(trace.SpanKindServer),
			)
			defer span.End()

			reply, err := handler(ctx, req)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			return reply, err
		}
	}
}

// ClientMiddleware returns a client-side tracing middleware.
func ClientMiddleware(tracer trace.Tracer) middleware.Middleware {
	if tracer == nil {
		tracer = otel.Tracer("github.com/lgzzz/mall-tracing")
	}
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			var spanName string
			if tr, ok := transport.FromClientContext(ctx); ok {
				spanName = tr.Operation()
			} else {
				spanName = "grpc.call"
			}

			ctx, span := tracer.Start(ctx, spanName,
				trace.WithSpanKind(trace.SpanKindClient),
			)
			defer span.End()

			reply, err := handler(ctx, req)
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			return reply, err
		}
	}
}

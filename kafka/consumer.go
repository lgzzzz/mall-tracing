package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// TracedReader wraps a kafka.Reader with OpenTelemetry tracing.
type TracedReader struct {
	inner      *kafka.Reader
	tracer     trace.Tracer
	propagator propagation.TextMapPropagator
}

// NewTracedConsumer creates a new TracedReader.
func NewTracedConsumer(inner *kafka.Reader, tracer trace.Tracer) *TracedReader {
	if tracer == nil {
		tracer = otel.Tracer("github.com/lgzzz/mall-tracing/kafka")
	}
	return &TracedReader{
		inner:      inner,
		tracer:     tracer,
		propagator: otel.GetTextMapPropagator(),
	}
}

// FetchMessage reads a message with tracing.
func (r *TracedReader) FetchMessage(ctx context.Context) (kafka.Message, error) {
	msg, err := r.inner.ReadMessage(ctx)
	if err != nil {
		return msg, err
	}

	// Extract trace context from message headers
	carrier := &headerCarrier{headers: &msg.Headers}
	childCtx := r.propagator.Extract(ctx, carrier)

	_, span := r.tracer.Start(childCtx, "kafka.consume",
		trace.WithSpanKind(trace.SpanKindConsumer),
	)
	span.End()

	return msg, nil
}

// Close closes the underlying reader.
func (r *TracedReader) Close() error {
	return r.inner.Close()
}

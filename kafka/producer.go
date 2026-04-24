package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// TracedWriter wraps a kafka.Writer with OpenTelemetry tracing.
type TracedWriter struct {
	inner      *kafka.Writer
	tracer     trace.Tracer
	propagator propagation.TextMapPropagator
}

// NewTracedProducer creates a new TracedWriter.
func NewTracedProducer(inner *kafka.Writer, tracer trace.Tracer) *TracedWriter {
	if tracer == nil {
		tracer = otel.Tracer("github.com/lgzzz/mall-tracing/kafka")
	}
	return &TracedWriter{
		inner:      inner,
		tracer:     tracer,
		propagator: otel.GetTextMapPropagator(),
	}
}

// WriteMessages writes messages to Kafka with tracing.
func (w *TracedWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	ctx, span := w.tracer.Start(ctx, "kafka.publish",
		trace.WithSpanKind(trace.SpanKindProducer),
	)
	defer span.End()

	for i := range msgs {
		// Inject trace context into message headers
		carrier := &headerCarrier{headers: &msgs[i].Headers}
		w.propagator.Inject(ctx, carrier)
	}

	if err := w.inner.WriteMessages(ctx, msgs...); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return err
	}
	return nil
}

// Close closes the underlying writer.
func (w *TracedWriter) Close() error {
	return w.inner.Close()
}

type headerCarrier struct {
	headers *[]kafka.Header
}

func (c *headerCarrier) Get(key string) string {
	for _, h := range *c.headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func (c *headerCarrier) Set(key string, value string) {
	*c.headers = append(*c.headers, kafka.Header{
		Key:   key,
		Value: []byte(value),
	})
}

func (c *headerCarrier) Keys() []string {
	keys := make([]string, 0, len(*c.headers))
	for _, h := range *c.headers {
		keys = append(keys, h.Key)
	}
	return keys
}

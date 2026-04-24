package data

import (
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"
)

// NewGORMTracingPlugin creates a GORM plugin for OpenTelemetry tracing.
func NewGORMTracingPlugin(tracer trace.TracerProvider) gorm.Plugin {
	return tracing.NewPlugin(tracing.WithTracerProvider(tracer))
}

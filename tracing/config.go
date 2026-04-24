package tracing

// Config holds OpenTelemetry tracer provider configuration.
type Config struct {
	ServiceName  string
	Version      string
	OTLPEndpoint string
	SampleRatio  float64
	Insecure     bool
}

package config

// Tracing holds OpenTelemetry tracing configuration.
type Tracing struct {
	Enabled     bool    `json:"enabled" yaml:"enabled"`
	Endpoint    string  `json:"endpoint" yaml:"endpoint"`
	SampleRatio float64 `json:"sample_ratio" yaml:"sample_ratio"`
}

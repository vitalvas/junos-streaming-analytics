package output

type Output interface {
	// Add a metric to the output; timestamp is in milliseconds
	AddMetric(name string, labels map[string]string, value float64, timestamp int64) error

	Send() error
}

type Config struct {
	Type   string            `yaml:"type" json:"type"`
	Config map[string]string `yaml:"config" json:"config"`
}

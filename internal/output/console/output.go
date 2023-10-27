package console

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/vitalvas/junos-streaming-analytics/internal/output"
)

type Output struct {
	data []Metric

	skipZero bool
}

type Metric struct {
	Name      string            `json:"name"`
	Labels    map[string]string `json:"labels"`
	Value     float64           `json:"value"`
	Timestamp int64             `json:"timestamp"`
}

func NewOutput(config output.Config) (*Output, error) {
	output := &Output{}

	if val, ok := config.Config["skip_zero"]; ok {
		skipZero, err := strconv.ParseBool(val)
		if err != nil {
			return nil, err
		}

		output.skipZero = skipZero
	}

	return output, nil
}

func (o *Output) AddMetric(name string, labels map[string]string, value float64, timestamp int64) error {
	if o.skipZero && value == 0 {
		return nil
	}

	o.data = append(o.data, Metric{
		Name:      name,
		Labels:    labels,
		Value:     value,
		Timestamp: timestamp,
	})

	return nil
}

func (o *Output) Send() error {
	if o.data == nil {
		return nil
	}

	data, err := json.Marshal(o.data)
	if err != nil {
		return err
	}

	fmt.Println("metrics:", len(o.data), "value:", string(data))

	o.data = nil

	return nil
}

package prometheus

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/prometheus/prometheus/prompb"
)

const (
	prefix = "juniper_telemetry"
)

type Output struct {
	URL string

	metrics []prompb.TimeSeries
}

func NewOutput() *Output {
	return &Output{}
}

func (o *Output) AddMetric(name string, labels map[string]string, value float64, timestamp int64) error {
	metric := prompb.TimeSeries{
		Labels: []prompb.Label{
			{
				Name:  "__name__",
				Value: fmt.Sprintf("%s_%s", prefix, strings.ToLower(name)),
			},
		},
		Samples: []prompb.Sample{
			{
				Value:     value,
				Timestamp: timestamp,
			},
		},
	}

	for k, v := range labels {
		metric.Labels = append(metric.Labels, prompb.Label{
			Name:  k,
			Value: v,
		})
	}

	o.metrics = append(o.metrics, metric)

	return nil
}

func (o *Output) Send() error {
	request := &prompb.WriteRequest{
		Timeseries: o.metrics,
	}

	data, err := proto.Marshal(request)
	if err != nil {
		return err
	}

	resp, err := http.Post(o.URL, "application/x-protobuf", bytes.NewReader(data))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	o.metrics = nil

	return nil
}

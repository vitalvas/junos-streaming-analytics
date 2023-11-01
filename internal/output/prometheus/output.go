package prometheus

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/prometheus/prometheus/prompb"

	"github.com/vitalvas/junos-streaming-analytics/internal/output"
)

const (
	prefix = "juniper_telemetry"
)

type Output struct {
	conf   output.Config
	url    string
	client *http.Client

	metrics []prompb.TimeSeries
	lock    sync.RWMutex
}

func NewOutput(config output.Config) (*Output, error) {
	output := &Output{
		conf: config,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	if val, ok := config.Config["url"]; ok {
		link, err := url.Parse(val)
		if err != nil {
			return nil, err
		}

		output.url = link.String()

	} else {
		return nil, fmt.Errorf("missing url")
	}

	return output, nil
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

	o.lock.Lock()
	o.metrics = append(o.metrics, metric)
	o.lock.Unlock()

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

	compressedData := snappy.Encode(nil, data)

	req, err := http.NewRequest(http.MethodPost, o.url, bytes.NewReader(compressedData))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Encoding", "snappy")
	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("X-Prometheus-Remote-Write-Version", "0.1.0")
	req.Header.Set("User-Agent", "junos-streaming-analytics/0.0.0")

	resp, err := o.client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	o.metrics = nil

	return nil
}

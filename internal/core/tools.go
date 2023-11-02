package core

import (
	"strings"

	"github.com/vitalvas/junos-streaming-analytics/internal/protos/jti"
)

func getJTIHostname(ts *jti.TelemetryStream) string {
	resp := ""
	if ts.SystemId != nil {
		systemID := ts.GetSystemId()

		// format: <router name>:<ip address>
		names := strings.Split(systemID, ":")
		if len(names) > 0 {
			resp = names[0]
		}
	}

	// TODO: dns resolve if empty
	return resp
}

func addToMetrics(metrics map[string]float64, key string, ptr interface{}) {
	if ptr == nil {
		return
	}

	switch v := ptr.(type) {
	case *float32:
		if v != nil {
			metrics[key] = float64(*v)
		}

	case *float64:
		if v != nil {
			metrics[key] = *v
		}

	case *int32:
		if v != nil {
			metrics[key] = float64(*v)
		}

	case *int64:
		if v != nil {
			metrics[key] = float64(*v)
		}

	case *uint32:
		if v != nil {
			metrics[key] = float64(*v)
		}

	case *uint64:
		if v != nil {
			metrics[key] = float64(*v)
		}

	case *bool:
		if v != nil {
			if *v {
				metrics[key] = 1
			} else {
				metrics[key] = 0
			}
		}
	}
}

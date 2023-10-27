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

package core

import (
	"fmt"
	"strings"

	"github.com/vitalvas/junos-streaming-analytics/internal/protos/jti"
	"google.golang.org/protobuf/proto"
)

func (app *App) RunJTIProcess() error {
	for {
		select {
		case <-app.shutdownCh:
			return nil

		default:
			msg := <-app.jtiMessages

			if msg.Data == nil {
				continue
			}

			if err := app.processJTIMessage(msg); err != nil {
				fmt.Printf("error processing jti message: %s", err.Error())
			}
		}

	}
}

func (app *App) processJTIMessage(msg jtiMessage) error {
	stream := &jti.TelemetryStream{}

	if err := proto.Unmarshal(msg.Data, stream); err != nil {
		return err
	}

	// skip unknown messages
	if !proto.HasExtension(stream.Enterprise, jti.E_JuniperNetworks) {
		return nil
	}

	jnsIface := proto.GetExtension(stream.Enterprise, jti.E_JuniperNetworks)
	if jnsIface == nil {
		return fmt.Errorf("error getting jti extension")
	}

	baseLabels := map[string]string{
		"hostname": getJTIHostname(stream),
		"instance": msg.Instance,
	}

	timestamp := int64(stream.GetTimestamp())
	if timestamp <= 0 {
		timestamp = msg.Timestamp.UnixMilli()
	}

	switch jns := jnsIface.(type) {
	case *jti.JuniperNetworksSensors:
		if proto.HasExtension(jns, jti.E_JnprInterfaceExt) {
			extension := proto.GetExtension(jns, jti.E_JnprInterfaceExt)
			if extension == nil {
				return fmt.Errorf("error getting jti E_JnprInterfaceExt extension")
			}

			switch p := extension.(type) {
			case *jti.Port:
				if err := app.jtiParsePort(msg.Instance, p, baseLabels, timestamp); err != nil {
					return err
				}

			default:
				return fmt.Errorf("unknown jti E_JnprInterfaceExt protobuf extension type: %T", p)
			}

			// } else if proto.HasExtension(jns, jti.E_JnprFirewallExt) {
			// 	extension := proto.GetExtension(jns, jti.E_JnprFirewallExt)
			// 	if extension == nil {
			// 		return fmt.Errorf("error getting jti E_JnprFirewallExt extension")
			// 	}

		}

	default:
		return fmt.Errorf("unknown jti protobuf extension type: %T", jns)
	}

	return nil
}

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

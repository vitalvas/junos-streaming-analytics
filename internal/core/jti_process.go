package core

import (
	"fmt"

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
		if err := app.processJuniperNetworksSensors(msg, jns, baseLabels, timestamp); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown jti protobuf extension type: %T", jns)
	}

	return nil
}

func (app *App) processJuniperNetworksSensors(msg jtiMessage, jns *jti.JuniperNetworksSensors, baseLabels map[string]string, timestamp int64) error {
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
	}

	if proto.HasExtension(jns, jti.E_JnprLogicalInterfaceExt) {
		extension := proto.GetExtension(jns, jti.E_JnprLogicalInterfaceExt)
		if extension == nil {
			return fmt.Errorf("error getting jti E_JnprLogicalInterfaceExt extension")
		}

		switch p := extension.(type) {
		case *jti.LogicalPort:
			if err := app.jtiParseLogicalPort(msg.Instance, p, baseLabels, timestamp); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown jti E_JnprLogicalInterfaceExt protobuf extension type: %T", p)
		}
	}

	if proto.HasExtension(jns, jti.E_JnprFirewallExt) {
		extension := proto.GetExtension(jns, jti.E_JnprFirewallExt)
		if extension == nil {
			return fmt.Errorf("error getting jti E_JnprFirewallExt extension")
		}

		switch p := extension.(type) {
		case *jti.Firewall:
			if err := app.jtiParseFirewall(msg.Instance, p, baseLabels, timestamp); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown jti E_JnprFirewallExt protobuf extension type: %T", p)
		}
	}

	if proto.HasExtension(jns, jti.E_JnprOpticsExt) {
		extension := proto.GetExtension(jns, jti.E_JnprOpticsExt)
		if extension == nil {
			return fmt.Errorf("error getting jti E_JnprOpticsExt extension")
		}

		switch p := extension.(type) {
		case *jti.Optics:
			if err := app.jtiParseOptics(msg.Instance, p, baseLabels, timestamp); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown jti E_JnprOpticsExt protobuf extension type: %T", p)
		}
	}

	if proto.HasExtension(jns, jti.E_NpuMemoryExt) {
		extension := proto.GetExtension(jns, jti.E_NpuMemoryExt)
		if extension == nil {
			return fmt.Errorf("error getting jti E_NpuMemoryExt extension")
		}

		switch p := extension.(type) {
		case *jti.NetworkProcessorMemoryUtilization:
			if err := app.jtiParseNpuMemoryUtilization(msg.Instance, p, baseLabels, timestamp); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown jti E_NpuMemoryExt protobuf extension type: %T", p)
		}
	}

	if proto.HasExtension(jns, jti.E_JnprNpuUtilizationExt) {
		extension := proto.GetExtension(jns, jti.E_JnprNpuUtilizationExt)
		if extension == nil {
			return fmt.Errorf("error getting jti E_JnprNpuUtilizationExt extension")
		}

		switch p := extension.(type) {
		case *jti.NetworkProcessorUtilization:
			if err := app.jtiParseNpuUtilization(msg.Instance, p, baseLabels, timestamp); err != nil {
				return err
			}

		default:
			return fmt.Errorf("unknown jti E_JnprNpuUtilizationExt protobuf extension type: %T", p)
		}
	}

	return nil
}

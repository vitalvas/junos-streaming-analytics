package app

import (
	"encoding/json"
	"log"
	"os"

	"github.com/junos-streaming-analytics/internal/decoder"
	"github.com/junos-streaming-analytics/internal/jti"
	"google.golang.org/protobuf/proto"
)

func (app *App) process() {
	for {
		bytes := <-app.message

		if bytes == nil {
			continue
		}

		ts := &jti.TelemetryStream{}

		if err := proto.Unmarshal(bytes, ts); err != nil {
			log.Println("error unmarshal message", err)
			continue
		}

		if proto.HasExtension(ts.Enterprise, jti.E_JuniperNetworks) {
			jnsIface := proto.GetExtension(ts.Enterprise, jti.E_JuniperNetworks)
			if jnsIface == nil {
				log.Println("failed to get extension")
				continue
			}

			switch jns := jnsIface.(type) {
			case *jti.JuniperNetworksSensors:
				if proto.HasExtension(jns, jti.E_JnprFirewallExt) {
					firewallIface := proto.GetExtension(jns, jti.E_JnprFirewallExt)
					if firewallIface == nil {
						log.Println("failed to get extension")
						continue
					}

					switch p := firewallIface.(type) {
					case *jti.Firewall:
						data := decoder.DecodeFirewall(ts, p)

						json.NewEncoder(os.Stdout).Encode(data)

					default:
						log.Printf("found no matching firewall: %s", p)
					}

				} else {
					log.Printf("received message with extension not currently supported: %s", jns)
				}

			default:
				log.Printf("unsupported JTI protobuf extension: %s", jns)
			}

		}
	}
}

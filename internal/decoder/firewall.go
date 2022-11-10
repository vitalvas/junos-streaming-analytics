package decoder

import "github.com/junos-streaming-analytics/internal/jti"

type Firewall struct {
	SensorName     string            `json:"sensor_name"`
	SystemID       string            `json:"system_id"`
	ComponentID    uint32            `json:"component_id"`
	SubComponentID uint32            `json:"sub_component_id"`
	Timestamp      uint64            `json:"timestamp,omitempty"`
	Counters       []FirewallCounter `json:"counters,omitempty"`
	Policers       []FirewallPolicer `json:"policers,omitempty"`
}

type FirewallCounter struct {
	Timestamp   uint64 `json:"timestamp,omitempty"`
	FilterName  string `json:"filter_name"`
	CounterName string `json:"counter_name"`
	Bytes       uint64 `json:"bytes"`
	Packets     uint64 `json:"packets"`
}

type FirewallPolicer struct {
	Timestamp        uint64 `json:"timestamp,omitempty"`
	FilterName       string `json:"filter_name"`
	PolicerName      string `json:"policer_name"`
	OutOfSpecPackets uint64 `json:"out_of_spec_packets"`
	OutOfSpecBytes   uint64 `json:"out_of_spec_bytes"`

	OfferedPackets     uint64 `json:"offered_packets,omitempty"`
	OfferedBytes       uint64 `json:"offered_bytes,omitempty"`
	TransmittedPackets uint64 `json:"transmitted_packets,omitempty"`
	TransmittedBytes   uint64 `json:"transmitted_bytes,omitempty"`
}

func DecodeFirewall(ts *jti.TelemetryStream, fw *jti.Firewall) *Firewall {
	data := &Firewall{
		SensorName:     ts.GetSensorName(),
		SystemID:       ts.GetSystemId(),
		ComponentID:    ts.GetComponentId(),
		SubComponentID: ts.GetSubComponentId(),
		Timestamp:      ts.GetTimestamp(),
	}

	for _, stats := range fw.GetFirewallStats() {
		for _, info := range stats.GetCounterStats() {
			data.Counters = append(data.Counters, FirewallCounter{
				Timestamp:   stats.GetTimestamp(),
				FilterName:  stats.GetFilterName(),
				CounterName: info.GetName(),
				Bytes:       info.GetBytes(),
				Packets:     info.GetPackets(),
			})
		}

		for _, info := range stats.GetPolicerStats() {
			fp := FirewallPolicer{
				Timestamp:        stats.GetTimestamp(),
				FilterName:       stats.GetFilterName(),
				PolicerName:      info.GetName(),
				OutOfSpecPackets: info.GetOutOfSpecPackets(),
				OutOfSpecBytes:   info.GetOutOfSpecBytes(),
			}

			extStats := info.GetExtendedPolicerStats()
			if extStats != nil {
				fp.OfferedPackets = extStats.GetOfferedPackets()
				fp.OfferedBytes = extStats.GetOfferedBytes()
				fp.TransmittedPackets = extStats.GetTransmittedPackets()
				fp.TransmittedBytes = extStats.GetTransmittedBytes()
			}

			data.Policers = append(data.Policers, fp)
		}
	}

	return data
}

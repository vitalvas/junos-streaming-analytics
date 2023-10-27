package core

import (
	"github.com/vitalvas/junos-streaming-analytics/internal/output"
	"github.com/vitalvas/junos-streaming-analytics/internal/protos/jti"
	"github.com/vitalvas/junos-streaming-analytics/internal/tools"
)

const firewallSeriesName = "firewall"

func (app *App) jtiParseFirewall(instance string, data *jti.Firewall, baseLabels map[string]string, timestamp int64) error {
	for _, firewallStats := range data.GetFirewallStats() {
		labels := tools.MergeMap(baseLabels, map[string]string{
			"filter_name": firewallStats.GetFilterName(),
		})

		metrics := map[string]float64{
			"changed": float64(firewallStats.GetTimestamp()),
		}

		for name, value := range metrics {
			if err := app.addMetricToOutput(instance, output.JoinMetricName(firewallSeriesName, name), labels, value, timestamp); err != nil {
				return err
			}
		}

		for _, memoryUsage := range firewallStats.GetMemoryUsage() {
			memoryLabels := tools.MergeMap(labels, map[string]string{
				"name": memoryUsage.GetName(),
			})

			if err := app.addMetricToOutput(instance, output.JoinMetricName(firewallSeriesName, "memory_usage"), memoryLabels, float64(memoryUsage.GetAllocated()), timestamp); err != nil {
				return err
			}
		}

		for _, counterStats := range firewallStats.GetCounterStats() {
			counterLabels := tools.MergeMap(labels, map[string]string{
				"name": counterStats.GetName(),
			})

			counterMetrics := map[string]float64{
				"packets": float64(counterStats.GetPackets()),
				"bytes":   float64(counterStats.GetBytes()),
			}

			for name, value := range counterMetrics {
				if err := app.addMetricToOutput(instance, output.JoinMetricName(firewallSeriesName, "counter_stats", name), counterLabels, value, timestamp); err != nil {
					return err
				}
			}
		}

		for _, policerStats := range firewallStats.GetPolicerStats() {
			policerLabels := tools.MergeMap(labels, map[string]string{
				"name": policerStats.GetName(),
			})

			policerMetrics := map[string]float64{
				"out_of_spec_packets": float64(policerStats.GetOutOfSpecPackets()),
				"out_of_spec_bytes":   float64(policerStats.GetOutOfSpecBytes()),
			}

			if stat := policerStats.GetExtendedPolicerStats(); stat != nil {
				policerMetrics[output.JoinMetricName("extended_policer_stats", "offered_packets")] = float64(stat.GetOfferedPackets())
				policerMetrics[output.JoinMetricName("extended_policer_stats", "offered_bytes")] = float64(stat.GetOfferedBytes())
				policerMetrics[output.JoinMetricName("extended_policer_stats", "transmitted_packets")] = float64(stat.GetTransmittedPackets())
				policerMetrics[output.JoinMetricName("extended_policer_stats", "transmitted_bytes")] = float64(stat.GetTransmittedBytes())
			}

			for name, value := range policerMetrics {
				if err := app.addMetricToOutput(instance, output.JoinMetricName(firewallSeriesName, "policer_stats", name), policerLabels, value, timestamp); err != nil {
					return err
				}
			}
		}

		for _, hierarchicalPolicerStats := range firewallStats.GetHierarchicalPolicerStats() {
			hierarchicalPolicerStatsLabels := tools.MergeMap(labels, map[string]string{
				"name": hierarchicalPolicerStats.GetName(),
			})

			hierarchicalPolicerStatsMetrics := map[string]float64{
				"premium_packets":   float64(hierarchicalPolicerStats.GetPremiumPackets()),
				"premium_bytes":     float64(hierarchicalPolicerStats.GetPremiumBytes()),
				"aggregate_packets": float64(hierarchicalPolicerStats.GetAggregatePackets()),
				"aggregate_bytes":   float64(hierarchicalPolicerStats.GetAggregateBytes()),
			}

			for name, value := range hierarchicalPolicerStatsMetrics {
				if err := app.addMetricToOutput(instance, output.JoinMetricName(firewallSeriesName, "hierarchical_policer_stats", name), hierarchicalPolicerStatsLabels, value, timestamp); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

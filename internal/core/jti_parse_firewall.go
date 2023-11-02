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

		metrics := make(map[string]float64)
		addToMetrics(metrics, "changed", firewallStats.GetTimestamp)

		for name, value := range metrics {
			if err := app.addMetricToOutput(instance, output.JoinMetricName(firewallSeriesName, name), labels, value, timestamp); err != nil {
				return err
			}
		}

		for _, memoryUsage := range firewallStats.GetMemoryUsage() {
			memoryLabels := tools.MergeMap(labels, map[string]string{
				"name": memoryUsage.GetName(),
			})

			memoryMetrics := make(map[string]float64)
			addToMetrics(memoryMetrics, "memory_usage", memoryUsage.GetAllocated)

			for name, value := range memoryMetrics {
				if err := app.addMetricToOutput(instance, output.JoinMetricName(firewallSeriesName, name), memoryLabels, value, timestamp); err != nil {
					return err
				}
			}
		}

		for _, counterStats := range firewallStats.GetCounterStats() {
			counterLabels := tools.MergeMap(labels, map[string]string{
				"name": counterStats.GetName(),
			})

			counterMetrics := make(map[string]float64)
			addToMetrics(counterMetrics, "packets", counterStats.Packets)
			addToMetrics(counterMetrics, "bytes", counterStats.Bytes)

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

			policerMetrics := make(map[string]float64)
			addToMetrics(policerMetrics, "out_of_spec_packets", policerStats.OutOfSpecPackets)
			addToMetrics(policerMetrics, "out_of_spec_bytes", policerStats.OutOfSpecBytes)

			if stat := policerStats.GetExtendedPolicerStats(); stat != nil {
				addToMetrics(policerMetrics, output.JoinMetricName("extended_policer_stats", "offered_packets"), stat.OfferedPackets)
				addToMetrics(policerMetrics, output.JoinMetricName("extended_policer_stats", "offered_bytes"), stat.OfferedBytes)
				addToMetrics(policerMetrics, output.JoinMetricName("extended_policer_stats", "transmitted_packets"), stat.TransmittedPackets)
				addToMetrics(policerMetrics, output.JoinMetricName("extended_policer_stats", "transmitted_bytes"), stat.TransmittedBytes)
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

			hierarchicalPolicerStatsMetrics := make(map[string]float64)
			addToMetrics(hierarchicalPolicerStatsMetrics, "premium_packets", hierarchicalPolicerStats.PremiumPackets)
			addToMetrics(hierarchicalPolicerStatsMetrics, "premium_bytes", hierarchicalPolicerStats.PremiumBytes)
			addToMetrics(hierarchicalPolicerStatsMetrics, "aggregate_packets", hierarchicalPolicerStats.AggregatePackets)
			addToMetrics(hierarchicalPolicerStatsMetrics, "aggregate_bytes", hierarchicalPolicerStats.AggregateBytes)

			for name, value := range hierarchicalPolicerStatsMetrics {
				if err := app.addMetricToOutput(instance, output.JoinMetricName(firewallSeriesName, "hierarchical_policer_stats", name), hierarchicalPolicerStatsLabels, value, timestamp); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

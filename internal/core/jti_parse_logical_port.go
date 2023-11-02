package core

import (
	"fmt"

	"github.com/vitalvas/junos-streaming-analytics/internal/output"
	"github.com/vitalvas/junos-streaming-analytics/internal/protos/jti"
	"github.com/vitalvas/junos-streaming-analytics/internal/tools"
)

const logicPortSeriesName = "network_logic"

func (app *App) jtiParseLogicalPort(instance string, data *jti.LogicalPort, baseLabels map[string]string, timestamp int64) error {
	for _, interfaceInfo := range data.GetInterfaceInfo() {
		labels := tools.MergeMap(baseLabels, map[string]string{
			"interface": interfaceInfo.GetIfName(),
		})

		if name := interfaceInfo.GetDescription(); name != "" {
			labels["description"] = name
		}

		if name := interfaceInfo.GetParentAeName(); name != "" {
			labels["parent_ae_name"] = name
		}

		metrics := make(map[string]float64)
		addToMetrics(metrics, "init_time", interfaceInfo.InitTime)
		addToMetrics(metrics, "snmp_if_index", interfaceInfo.SnmpIfIndex)
		addToMetrics(metrics, "last_change", interfaceInfo.LastChange)
		addToMetrics(metrics, "high_speed", interfaceInfo.HighSpeed)

		if stats := interfaceInfo.GetIngressStats(); stats != nil {
			addToMetrics(metrics, output.JoinMetricName("ingress_stats", "if_packets"), stats.IfPackets)
			addToMetrics(metrics, output.JoinMetricName("ingress_stats", "if_octets"), stats.IfOctets)
			addToMetrics(metrics, output.JoinMetricName("ingress_stats", "if_ucast_packets"), stats.IfUcastPackets)
			addToMetrics(metrics, output.JoinMetricName("ingress_stats", "if_mcast_packets"), stats.IfMcastPackets)

			for _, ifFcStats := range stats.GetIfFcStats() {
				if ifFcStats == nil {
					continue
				}

				ifFcStatsLabels := tools.MergeMap(labels, map[string]string{
					"if_family": ifFcStats.GetIfFamily(),
					"fc_number": fmt.Sprintf("%d", ifFcStats.GetFcNumber()),
				})

				ifFcStatsMetrics := make(map[string]float64)
				addToMetrics(ifFcStatsMetrics, "if_packets", ifFcStats.IfPackets)
				addToMetrics(ifFcStatsMetrics, "if_octets", ifFcStats.IfOctets)

				for key, value := range ifFcStatsMetrics {
					if err := app.addMetricToOutput(instance, output.JoinMetricName(logicPortSeriesName, "ingress_stats", "if_fc_stats", key), ifFcStatsLabels, value, timestamp); err != nil {
						return err
					}
				}
			}

			for _, ifFaStats := range stats.GetIfFaStats() {
				if ifFaStats == nil {
					continue
				}

				ifFaStatsLabels := tools.MergeMap(labels, map[string]string{
					"if_family": ifFaStats.GetIfFamily(),
				})

				ifFaStatsMetrics := make(map[string]float64)
				addToMetrics(ifFaStatsMetrics, "if_packets", ifFaStats.IfPackets)
				addToMetrics(ifFaStatsMetrics, "if_octets", ifFaStats.IfOctets)
				addToMetrics(ifFaStatsMetrics, "if_v6_packets", ifFaStats.IfV6Packets)
				addToMetrics(ifFaStatsMetrics, "if_v6_octets", ifFaStats.IfV6Octets)
				addToMetrics(ifFaStatsMetrics, "if_mcast_packets", ifFaStats.IfMcastPackets)
				addToMetrics(ifFaStatsMetrics, "if_mcast_octets", ifFaStats.IfMcastOctets)

				for key, value := range ifFaStatsMetrics {
					if err := app.addMetricToOutput(instance, output.JoinMetricName(logicPortSeriesName, "ingress_stats", "if_fa_stats", key), ifFaStatsLabels, value, timestamp); err != nil {
						return err
					}
				}
			}
		}

		if stats := interfaceInfo.GetEgressStats(); stats != nil {
			addToMetrics(metrics, output.JoinMetricName("egress_stats", "if_packets"), stats.IfPackets)
			addToMetrics(metrics, output.JoinMetricName("egress_stats", "if_octets"), stats.IfOctets)
			addToMetrics(metrics, output.JoinMetricName("egress_stats", "if_ucast_packets"), stats.IfUcastPackets)
			addToMetrics(metrics, output.JoinMetricName("egress_stats", "if_mcast_packets"), stats.IfMcastPackets)

			for _, ifFaStats := range stats.GetIfFaStats() {
				if ifFaStats == nil {
					continue
				}

				ifFaStatsLabels := tools.MergeMap(labels, map[string]string{
					"if_family": ifFaStats.GetIfFamily(),
				})

				ifFaStatsMetrics := make(map[string]float64)
				addToMetrics(ifFaStatsMetrics, "if_packets", ifFaStats.IfPackets)
				addToMetrics(ifFaStatsMetrics, "if_octets", ifFaStats.IfOctets)
				addToMetrics(ifFaStatsMetrics, "if_v6_packets", ifFaStats.IfV6Packets)
				addToMetrics(ifFaStatsMetrics, "if_v6_octets", ifFaStats.IfV6Octets)
				addToMetrics(ifFaStatsMetrics, "if_mcast_packets", ifFaStats.IfMcastPackets)
				addToMetrics(ifFaStatsMetrics, "if_mcast_octets", ifFaStats.IfMcastOctets)

				for key, value := range ifFaStatsMetrics {
					if err := app.addMetricToOutput(instance, output.JoinMetricName(logicPortSeriesName, "ingress_stats", "if_fa_stats", key), ifFaStatsLabels, value, timestamp); err != nil {
						return err
					}
				}
			}
		}

		for dir, stats := range map[string][]*jti.LogicalInterfaceQueueStats{
			"ingress_queue_info": interfaceInfo.GetIngressQueueInfo(),
			"egress_queue_info":  interfaceInfo.GetEgressQueueInfo(),
		} {
			if stats == nil {
				continue
			}

			for _, queueStat := range stats {
				if queueStat == nil {
					continue
				}

				queueLabels := tools.MergeMap(labels, map[string]string{
					"queue_number": fmt.Sprintf("%d", queueStat.GetQueueNumber()),
				})

				queueMetrics := make(map[string]float64)
				addToMetrics(queueMetrics, "packets", queueStat.Packets)
				addToMetrics(queueMetrics, "bytes", queueStat.Bytes)
				addToMetrics(queueMetrics, "tail_drop_packets", queueStat.TailDropPackets)
				addToMetrics(queueMetrics, "rate_limit_drop_packets", queueStat.RateLimitDropPackets)
				addToMetrics(queueMetrics, "rate_limit_drop_bytes", queueStat.RateLimitDropBytes)
				addToMetrics(queueMetrics, "red_drop_packets", queueStat.RedDropPackets)
				addToMetrics(queueMetrics, "red_drop_bytes", queueStat.RedDropBytes)
				addToMetrics(queueMetrics, "average_buffer_occupancy", queueStat.AverageBufferOccupancy)
				addToMetrics(queueMetrics, "current_buffer_occupancy", queueStat.CurrentBufferOccupancy)
				addToMetrics(queueMetrics, "peak_buffer_occupancy", queueStat.PeakBufferOccupancy)
				addToMetrics(queueMetrics, "allocated_buffer_size", queueStat.AllocatedBufferSize)

				for queueMetricsName, queueMetricsValue := range queueMetrics {
					if err := app.addMetricToOutput(instance, output.JoinMetricName(logicPortSeriesName, dir, queueMetricsName), queueLabels, queueMetricsValue, timestamp); err != nil {
						return err
					}
				}
			}

		}

		for name, value := range metrics {
			if err := app.addMetricToOutput(instance, output.JoinMetricName(logicPortSeriesName, name), labels, value, timestamp); err != nil {
				return err
			}
		}
	}

	return nil
}

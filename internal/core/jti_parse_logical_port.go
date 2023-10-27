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

		metrics := map[string]float64{
			"init_time":     float64(interfaceInfo.GetInitTime()),
			"snmp_if_index": float64(interfaceInfo.GetSnmpIfIndex()),
			"last_change":   float64(interfaceInfo.GetLastChange()),
			"high_speed":    float64(interfaceInfo.GetHighSpeed()),
		}

		if stats := interfaceInfo.GetIngressStats(); stats != nil {
			metrics[output.JoinMetricName("ingress_stats", "if_packets")] = float64(stats.GetIfPackets())
			metrics[output.JoinMetricName("ingress_stats", "if_octets")] = float64(stats.GetIfOctets())
			metrics[output.JoinMetricName("ingress_stats", "if_ucast_packets")] = float64(stats.GetIfUcastPackets())
			metrics[output.JoinMetricName("ingress_stats", "if_mcast_packets")] = float64(stats.GetIfMcastPackets())

			for _, ifFcStats := range stats.GetIfFcStats() {
				if ifFcStats == nil {
					continue
				}

				ifFcStatsLabels := tools.MergeMap(labels, map[string]string{
					"if_family": ifFcStats.GetIfFamily(),
					"fc_number": fmt.Sprintf("%d", ifFcStats.GetFcNumber()),
				})

				ifFcStatsMetrics := map[string]float64{
					"if_packets": float64(ifFcStats.GetIfPackets()),
					"if_octets":  float64(ifFcStats.GetIfOctets()),
				}

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

				ifFaStatsMetrics := map[string]float64{
					"if_packets":       float64(ifFaStats.GetIfPackets()),
					"if_octets":        float64(ifFaStats.GetIfOctets()),
					"if_v6_packets":    float64(ifFaStats.GetIfV6Packets()),
					"if_v6_octets":     float64(ifFaStats.GetIfV6Octets()),
					"if_mcast_packets": float64(ifFaStats.GetIfMcastPackets()),
					"if_mcast_octets":  float64(ifFaStats.GetIfMcastOctets()),
				}

				for key, value := range ifFaStatsMetrics {
					if err := app.addMetricToOutput(instance, output.JoinMetricName(logicPortSeriesName, "ingress_stats", "if_fa_stats", key), ifFaStatsLabels, value, timestamp); err != nil {
						return err
					}
				}
			}
		}

		if stats := interfaceInfo.GetEgressStats(); stats != nil {
			metrics[output.JoinMetricName("egress_stats", "if_packets")] = float64(stats.GetIfPackets())
			metrics[output.JoinMetricName("egress_stats", "if_octets")] = float64(stats.GetIfOctets())
			metrics[output.JoinMetricName("egress_stats", "if_ucast_packets")] = float64(stats.GetIfUcastPackets())
			metrics[output.JoinMetricName("egress_stats", "if_mcast_packets")] = float64(stats.GetIfMcastPackets())

			for _, ifFaStats := range stats.GetIfFaStats() {
				if ifFaStats == nil {
					continue
				}

				ifFaStatsLabels := tools.MergeMap(labels, map[string]string{
					"if_family": ifFaStats.GetIfFamily(),
				})

				ifFaStatsMetrics := map[string]float64{
					"if_packets":       float64(ifFaStats.GetIfPackets()),
					"if_octets":        float64(ifFaStats.GetIfOctets()),
					"if_v6_packets":    float64(ifFaStats.GetIfV6Packets()),
					"if_v6_octets":     float64(ifFaStats.GetIfV6Octets()),
					"if_mcast_packets": float64(ifFaStats.GetIfMcastPackets()),
					"if_mcast_octets":  float64(ifFaStats.GetIfMcastOctets()),
				}

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

				queueMetrics := map[string]float64{
					"packets":                  float64(queueStat.GetPackets()),
					"bytes":                    float64(queueStat.GetBytes()),
					"tail_drop_packets":        float64(queueStat.GetTailDropPackets()),
					"rate_limit_drop_packets":  float64(queueStat.GetRateLimitDropPackets()),
					"rate_limit_drop_bytes":    float64(queueStat.GetRateLimitDropBytes()),
					"red_drop_packets":         float64(queueStat.GetRedDropPackets()),
					"red_drop_bytes":           float64(queueStat.GetRedDropBytes()),
					"average_buffer_occupancy": float64(queueStat.GetAverageBufferOccupancy()),
					"current_buffer_occupancy": float64(queueStat.GetCurrentBufferOccupancy()),
					"peak_buffer_occupancy":    float64(queueStat.GetPeakBufferOccupancy()),
					"allocated_buffer_size":    float64(queueStat.GetAllocatedBufferSize()),
				}

				for queueMetricsName, queueMetricsValue := range queueMetrics {
					if err := app.addMetricToOutput(instance, output.JoinMetricName(logicPortSeriesName, dir, queueMetricsName), queueLabels, queueMetricsValue, timestamp); err != nil {
						return err
					}
				}
			}

		}

	}

	return nil
}

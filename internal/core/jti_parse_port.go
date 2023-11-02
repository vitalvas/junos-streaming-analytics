package core

import (
	"fmt"

	"github.com/vitalvas/junos-streaming-analytics/internal/output"
	"github.com/vitalvas/junos-streaming-analytics/internal/protos/jti"
	"github.com/vitalvas/junos-streaming-analytics/internal/tools"
)

const portSeriesName = "network"

func (app *App) jtiParsePort(instance string, data *jti.Port, baseLabels map[string]string, timestamp int64) error {
	for _, interfaceStat := range data.GetInterfaceStats() {
		labels := tools.MergeMap(baseLabels, map[string]string{
			"interface": interfaceStat.GetIfName(),
		})

		if name := interfaceStat.GetIfDescription(); name != "" {
			labels["description"] = name
		}

		if name := interfaceStat.GetParentAeName(); name != "" {
			labels["parent_ae_name"] = name
		}

		metrics := make(map[string]float64)

		addToMetrics(metrics, "init_time", interfaceStat.InitTime)
		addToMetrics(metrics, "snmp_if_index", interfaceStat.SnmpIfIndex)
		addToMetrics(metrics, "if_transitions", interfaceStat.IfTransitions)
		addToMetrics(metrics, "if_last_change", interfaceStat.IfLastChange)
		addToMetrics(metrics, "if_high_speed", interfaceStat.IfHighSpeed)

		for dir, stats := range map[string]*jti.InterfaceStats{
			"ingress_stats": interfaceStat.GetIngressStats(),
			"egress_stats":  interfaceStat.GetEgressStats(),
		} {
			if stats == nil {
				continue
			}

			addToMetrics(metrics, output.JoinMetricName(dir, "if_pkts"), stats.IfPkts)
			addToMetrics(metrics, output.JoinMetricName(dir, "if_octets"), stats.IfOctets)
			addToMetrics(metrics, output.JoinMetricName(dir, "if_1sec_pkts"), stats.If_1SecPkts)
			addToMetrics(metrics, output.JoinMetricName(dir, "if_1sec_octets"), stats.If_1SecOctets)
			addToMetrics(metrics, output.JoinMetricName(dir, "if_uc_pkts"), stats.IfUcPkts)
			addToMetrics(metrics, output.JoinMetricName(dir, "if_mc_pkts"), stats.IfMcPkts)
			addToMetrics(metrics, output.JoinMetricName(dir, "if_bc_pkts"), stats.IfBcPkts)
			addToMetrics(metrics, output.JoinMetricName(dir, "if_error"), stats.IfError)
			addToMetrics(metrics, output.JoinMetricName(dir, "if_pause_pkts"), stats.IfPausePkts)
			addToMetrics(metrics, output.JoinMetricName(dir, "if_unknown_proto_pkts"), stats.IfUnknownProtoPkts)
		}

		if stats := interfaceStat.GetIngressErrors(); stats != nil {
			addToMetrics(metrics, output.JoinMetricName("ingress_errors", "if_errors"), stats.IfErrors)
			addToMetrics(metrics, output.JoinMetricName("ingress_errors", "if_in_qdrops"), stats.IfInQdrops)
			addToMetrics(metrics, output.JoinMetricName("ingress_errors", "if_in_frame_errors"), stats.IfInFrameErrors)
			addToMetrics(metrics, output.JoinMetricName("ingress_errors", "if_discards"), stats.IfDiscards)
			addToMetrics(metrics, output.JoinMetricName("ingress_errors", "if_in_runts"), stats.IfInRunts)
			addToMetrics(metrics, output.JoinMetricName("ingress_errors", "if_in_l3_incompletes"), stats.IfInL3Incompletes)
			addToMetrics(metrics, output.JoinMetricName("ingress_errors", "if_in_l2chan_errors"), stats.IfInL2ChanErrors)
			addToMetrics(metrics, output.JoinMetricName("ingress_errors", "if_in_l2_mismatch_timeouts"), stats.IfInL2MismatchTimeouts)
			addToMetrics(metrics, output.JoinMetricName("ingress_errors", "if_in_fifo_errors"), stats.IfInFifoErrors)
			addToMetrics(metrics, output.JoinMetricName("ingress_errors", "if_in_resource_errors"), stats.IfInResourceErrors)
		}

		if stats := interfaceStat.GetEgressErrors(); stats != nil {
			addToMetrics(metrics, output.JoinMetricName("egress_errors", "if_errors"), stats.IfErrors)
			addToMetrics(metrics, output.JoinMetricName("egress_errors", "if_discards"), stats.IfDiscards)
		}

		for name, value := range metrics {
			if err := app.addMetricToOutput(instance, output.JoinMetricName(portSeriesName, name), labels, value, timestamp); err != nil {
				return err
			}
		}

		for name, stats := range map[string][]*jti.QueueStats{
			"egress_queue_info":  interfaceStat.GetEgressQueueInfo(),
			"ingress_queue_info": interfaceStat.GetIngressQueueInfo(),
		} {
			if stats == nil {
				continue
			}

			for _, queueStat := range stats {
				queueLabels := tools.MergeMap(labels, map[string]string{
					"queue_number": fmt.Sprintf("%d", queueStat.GetQueueNumber()),
				})

				queueMetrics := make(map[string]float64)

				addToMetrics(queueMetrics, "packets", queueStat.Packets)
				addToMetrics(queueMetrics, "bytes", queueStat.Bytes)
				addToMetrics(queueMetrics, "tail_drop_packets", queueStat.TailDropPackets)
				addToMetrics(queueMetrics, "rl_drop_packets", queueStat.RlDropPackets)
				addToMetrics(queueMetrics, "rl_drop_bytes", queueStat.RlDropBytes)
				addToMetrics(queueMetrics, "red_drop_packets", queueStat.RedDropPackets)
				addToMetrics(queueMetrics, "red_drop_bytes", queueStat.RedDropBytes)
				addToMetrics(queueMetrics, "avg_buffer_occupancy", queueStat.AvgBufferOccupancy)
				addToMetrics(queueMetrics, "cur_buffer_occupancy", queueStat.CurBufferOccupancy)
				addToMetrics(queueMetrics, "peak_buffer_occupancy", queueStat.PeakBufferOccupancy)
				addToMetrics(queueMetrics, "allocated_buffer_size", queueStat.AllocatedBufferSize)

				for queueMetricsName, queueMetricsValue := range queueMetrics {
					if err := app.addMetricToOutput(instance, output.JoinMetricName(portSeriesName, name, queueMetricsName), queueLabels, queueMetricsValue, timestamp); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

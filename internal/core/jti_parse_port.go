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

		metrics := map[string]float64{
			"init_time":      float64(interfaceStat.GetInitTime()),
			"snmp_if_index":  float64(interfaceStat.GetSnmpIfIndex()),
			"if_transitions": float64(interfaceStat.GetIfTransitions()),
			"if_last_change": float64(interfaceStat.GetIfLastChange()),
			"if_high_speed":  float64(interfaceStat.GetIfHighSpeed()),
		}

		for dir, stats := range map[string]*jti.InterfaceStats{
			"ingress_stats": interfaceStat.GetIngressStats(),
			"egress_stats":  interfaceStat.GetEgressStats(),
		} {
			if stats == nil {
				continue
			}

			metrics[output.JoinMetricName(dir, "if_pkts")] = float64(stats.GetIfPkts())
			metrics[output.JoinMetricName(dir, "if_octets")] = float64(stats.GetIfOctets())
			metrics[output.JoinMetricName(dir, "if_1sec_pkts")] = float64(stats.GetIf_1SecPkts())
			metrics[output.JoinMetricName(dir, "if_1sec_octets")] = float64(stats.GetIf_1SecOctets())
			metrics[output.JoinMetricName(dir, "if_uc_pkts")] = float64(stats.GetIfUcPkts())
			metrics[output.JoinMetricName(dir, "if_mc_pkts")] = float64(stats.GetIfMcPkts())
			metrics[output.JoinMetricName(dir, "if_bc_pkts")] = float64(stats.GetIfBcPkts())
			metrics[output.JoinMetricName(dir, "if_error")] = float64(stats.GetIfError())
			metrics[output.JoinMetricName(dir, "if_pause_pkts")] = float64(stats.GetIfPausePkts())
			metrics[output.JoinMetricName(dir, "if_unknown_proto_pkts")] = float64(stats.GetIfUnknownProtoPkts())
		}

		if stats := interfaceStat.GetIngressErrors(); data != nil {
			metrics[output.JoinMetricName("ingress_errors", "if_errors")] = float64(stats.GetIfErrors())
			metrics[output.JoinMetricName("ingress_errors", "if_in_qdrops")] = float64(stats.GetIfInQdrops())
			metrics[output.JoinMetricName("ingress_errors", "if_in_frame_errors")] = float64(stats.GetIfInFrameErrors())
			metrics[output.JoinMetricName("ingress_errors", "if_discards")] = float64(stats.GetIfDiscards())
			metrics[output.JoinMetricName("ingress_errors", "if_in_runts")] = float64(stats.GetIfInRunts())
			metrics[output.JoinMetricName("ingress_errors", "if_in_l3_incompletes")] = float64(stats.GetIfInL3Incompletes())
			metrics[output.JoinMetricName("ingress_errors", "if_in_l2chan_errors")] = float64(stats.GetIfInL2ChanErrors())
			metrics[output.JoinMetricName("ingress_errors", "if_in_l2_mismatch_timeouts")] = float64(stats.GetIfInL2MismatchTimeouts())
			metrics[output.JoinMetricName("ingress_errors", "if_in_fifo_errors")] = float64(stats.GetIfInFifoErrors())
			metrics[output.JoinMetricName("ingress_errors", "if_in_resource_errors")] = float64(stats.GetIfInResourceErrors())
		}

		if stats := interfaceStat.GetEgressErrors(); data != nil {
			metrics[output.JoinMetricName("egress_errors", "if_errors")] = float64(stats.GetIfErrors())
			metrics[output.JoinMetricName("egress_errors", "if_discards")] = float64(stats.GetIfDiscards())
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

				queueMetrics := map[string]float64{
					"packets":               float64(queueStat.GetPackets()),
					"bytes":                 float64(queueStat.GetBytes()),
					"tail_drop_packets":     float64(queueStat.GetTailDropPackets()),
					"rl_drop_packets":       float64(queueStat.GetRlDropPackets()),
					"rl_drop_bytes":         float64(queueStat.GetRlDropBytes()),
					"red_drop_packets":      float64(queueStat.GetRedDropPackets()),
					"red_drop_bytes":        float64(queueStat.GetRedDropBytes()),
					"avg_buffer_occupancy":  float64(queueStat.GetAvgBufferOccupancy()),
					"cur_buffer_occupancy":  float64(queueStat.GetCurBufferOccupancy()),
					"peak_buffer_occupancy": float64(queueStat.GetPeakBufferOccupancy()),
					"allocated_buffer_size": float64(queueStat.GetAllocatedBufferSize()),
				}

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

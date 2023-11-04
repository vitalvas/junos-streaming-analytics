package core

import (
	"github.com/vitalvas/junos-streaming-analytics/internal/output"
	"github.com/vitalvas/junos-streaming-analytics/internal/protos/jti"
	"github.com/vitalvas/junos-streaming-analytics/internal/tools"
)

const npuUtilizationSeriesName = "npu_utilization"

func (app *App) jtiParseNpuUtilization(instance string, data *jti.NetworkProcessorUtilization, baseLabels map[string]string, timestamp int64) error {
	for _, utilStats := range data.GetNpuUtilStats() {
		if utilStats == nil {
			continue
		}

		labels := tools.MergeMap(baseLabels, map[string]string{
			"identifier": utilStats.GetIdentifier(),
		})

		metrics := make(map[string]float64)
		addToMetrics(metrics, "utilization", utilStats.Utilization)

		for key, value := range metrics {
			if err := app.addMetricToOutput(instance, output.JoinMetricName(npuUtilizationSeriesName, key), labels, value, timestamp); err != nil {
				return err
			}
		}

		for _, packets := range utilStats.GetPackets() {
			if packets == nil {
				continue
			}

			packetLabels := tools.MergeMap(labels, map[string]string{
				"internal_identifier": packets.GetIdentifier(),
			})

			packetMetrics := make(map[string]float64)
			addToMetrics(packetMetrics, "rate", packets.Rate)
			addToMetrics(packetMetrics, "average_instructions_per_packet", packets.AverageInstructionsPerPacket)
			addToMetrics(packetMetrics, "average_wait_cycles_per_packet", packets.AverageWaitCyclesPerPacket)
			addToMetrics(packetMetrics, "average_cycles_per_packet", packets.AverageCyclesPerPacket)

			for key, value := range packetMetrics {
				if err := app.addMetricToOutput(instance, output.JoinMetricName(npuUtilizationSeriesName, "packets", key), packetLabels, value, timestamp); err != nil {
					return err
				}
			}
		}

		for _, memory := range utilStats.GetMemory() {
			if memory == nil {
				continue
			}

			memoryLabels := tools.MergeMap(labels, map[string]string{
				"internal_name": memory.GetName(),
			})

			memoryMetrics := make(map[string]float64)
			addToMetrics(memoryMetrics, "average_util", memory.AverageUtil)
			addToMetrics(memoryMetrics, "highest_util", memory.HighestUtil)
			addToMetrics(memoryMetrics, "lowest_util", memory.LowestUtil)
			addToMetrics(memoryMetrics, "average_cache_hit_rate", memory.AverageCacheHitRate)
			addToMetrics(memoryMetrics, "highest_cache_hit_rate", memory.HighestCacheHitRate)
			addToMetrics(memoryMetrics, "lowest_cache_hit_rate", memory.LowestCacheHitRate)

			for key, value := range memoryMetrics {
				if err := app.addMetricToOutput(instance, output.JoinMetricName(npuUtilizationSeriesName, "memory", key), memoryLabels, value, timestamp); err != nil {
					return err
				}
			}
		}

	}
	return nil
}

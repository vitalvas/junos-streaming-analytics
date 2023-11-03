package core

import (
	"github.com/vitalvas/junos-streaming-analytics/internal/output"
	"github.com/vitalvas/junos-streaming-analytics/internal/protos/jti"
	"github.com/vitalvas/junos-streaming-analytics/internal/tools"
)

const npuMemorySeriesName = "npu_memory"

func (app *App) jtiParseNpuMemoryUtilization(instance string, data *jti.NetworkProcessorMemoryUtilization, baseLabels map[string]string, timestamp int64) error {
	for _, memoryStats := range data.GetMemoryStats() {
		if memoryStats == nil {
			continue
		}

		labels := tools.MergeMap(baseLabels, map[string]string{
			"identifier": memoryStats.GetIdentifier(),
		})

		for _, summary := range memoryStats.GetSummary() {
			if summary == nil {
				continue
			}

			summaryLabels := tools.MergeMap(labels, map[string]string{
				"resource_name": summary.GetResourceName(),
			})

			summaryMetrics := make(map[string]float64)
			addToMetrics(summaryMetrics, "size", summary.Size)
			addToMetrics(summaryMetrics, "allocated", summary.Allocated)
			addToMetrics(summaryMetrics, "utilization", summary.Utilization)

			for key, value := range summaryMetrics {
				if err := app.addMetricToOutput(instance, output.JoinMetricName(npuMemorySeriesName, "summary", key), summaryLabels, value, timestamp); err != nil {
					return err
				}
			}
		}

		for _, partition := range memoryStats.GetPartition() {
			if partition == nil {
				continue
			}

			partitionLabels := tools.MergeMap(labels, map[string]string{
				"partition_name":   partition.GetName(),
				"application_name": partition.GetApplicationName(),
			})

			partitionMetrics := make(map[string]float64)
			addToMetrics(partitionMetrics, "bytes_allocated", partition.BytesAllocated)
			addToMetrics(partitionMetrics, "allocation_count", partition.AllocationCount)
			addToMetrics(partitionMetrics, "free_count", partition.FreeCount)

			for key, value := range partitionMetrics {
				if err := app.addMetricToOutput(instance, output.JoinMetricName(npuMemorySeriesName, "partition", key), partitionLabels, value, timestamp); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

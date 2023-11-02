package core

import (
	"fmt"

	"github.com/vitalvas/junos-streaming-analytics/internal/output"
	"github.com/vitalvas/junos-streaming-analytics/internal/protos/jti"
	"github.com/vitalvas/junos-streaming-analytics/internal/tools"
)

const opticsSeriesName = "optics"

func (app *App) jtiParseOptics(instance string, data *jti.Optics, baseLabels map[string]string, timestamp int64) error {
	for _, opticsInfos := range data.GetOpticsDiag() {
		labels := tools.MergeMap(baseLabels, map[string]string{
			"interface": opticsInfos.GetIfName(),
		})

		metrics := map[string]float64{
			"snmp_if_index": float64(opticsInfos.GetSnmpIfIndex()),
		}

		if stats := opticsInfos.GetOpticsDiagStats(); stats != nil {
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "optics_type"), stats.OpticsType)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "module_temp"), stats.ModuleTemp)

			addToMetrics(metrics, output.JoinMetricName("diag_stats", "module_temp_high_alarm_threshold"), stats.ModuleTempHighAlarmThreshold)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "module_temp_low_alarm_threshold"), stats.ModuleTempLowAlarmThreshold)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "module_temp_high_warning_threshold"), stats.ModuleTempHighWarningThreshold)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "module_temp_low_warning_threshold"), stats.ModuleTempLowWarningThreshold)

			addToMetrics(metrics, output.JoinMetricName("diag_stats", "laser_output_power_high_alarm_threshold_dbm"), stats.LaserOutputPowerHighAlarmThresholdDbm)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "laser_output_power_low_alarm_threshold_dbm"), stats.LaserOutputPowerLowAlarmThresholdDbm)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "laser_output_power_high_warning_threshold_dbm"), stats.LaserOutputPowerHighWarningThresholdDbm)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "laser_output_power_low_warning_threshold_dbm"), stats.LaserOutputPowerLowWarningThresholdDbm)

			addToMetrics(metrics, output.JoinMetricName("diag_stats", "laser_rx_power_high_alarm_threshold_dbm"), stats.LaserRxPowerHighAlarmThresholdDbm)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "laser_rx_power_low_alarm_threshold_dbm"), stats.LaserRxPowerLowAlarmThresholdDbm)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "laser_rx_power_high_warning_threshold_dbm"), stats.LaserRxPowerHighWarningThresholdDbm)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "laser_rx_power_low_warning_threshold_dbm"), stats.LaserRxPowerLowWarningThresholdDbm)

			addToMetrics(metrics, output.JoinMetricName("diag_stats", "laser_bias_current_high_alarm_threshold"), stats.LaserBiasCurrentHighAlarmThreshold)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "laser_bias_current_low_alarm_threshold"), stats.LaserBiasCurrentLowAlarmThreshold)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "laser_bias_current_high_warning_threshold"), stats.LaserBiasCurrentHighWarningThreshold)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "laser_bias_current_low_warning_threshold"), stats.LaserBiasCurrentLowWarningThreshold)

			addToMetrics(metrics, output.JoinMetricName("diag_stats", "module_temp_high_alarm"), stats.ModuleTempHighAlarm)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "module_temp_low_alarm"), stats.ModuleTempLowAlarm)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "module_temp_high_warning"), stats.ModuleTempHighWarning)
			addToMetrics(metrics, output.JoinMetricName("diag_stats", "module_temp_low_warning"), stats.ModuleTempLowWarning)

			for name, value := range metrics {
				if err := app.addMetricToOutput(instance, output.JoinMetricName(opticsSeriesName, name), labels, value, timestamp); err != nil {
					return err
				}
			}

			for _, laneStats := range stats.GetOpticsLaneDiagStats() {
				if laneStats == nil {
					continue
				}

				laneLabels := tools.MergeMap(labels, map[string]string{
					"lane_number": fmt.Sprintf("%d", laneStats.GetLaneNumber()),
				})

				laneMetrics := make(map[string]float64)

				addToMetrics(laneMetrics, "lane_laser_temperature", laneStats.LaneLaserTemperature)
				addToMetrics(laneMetrics, "lane_laser_output_power_dbm", laneStats.LaneLaserOutputPowerDbm)
				addToMetrics(laneMetrics, "lane_laser_receiver_power_dbm", laneStats.LaneLaserReceiverPowerDbm)
				addToMetrics(laneMetrics, "lane_laser_bias_current", laneStats.LaneLaserBiasCurrent)

				addToMetrics(laneMetrics, "lane_laser_output_power_high_alarm", laneStats.LaneLaserOutputPowerHighAlarm)
				addToMetrics(laneMetrics, "lane_laser_output_power_low_alarm", laneStats.LaneLaserOutputPowerLowAlarm)
				addToMetrics(laneMetrics, "lane_laser_output_power_high_warning", laneStats.LaneLaserOutputPowerHighWarning)
				addToMetrics(laneMetrics, "lane_laser_output_power_low_warning", laneStats.LaneLaserOutputPowerLowWarning)

				addToMetrics(laneMetrics, "lane_laser_receiver_power_high_alarm", laneStats.LaneLaserReceiverPowerHighAlarm)
				addToMetrics(laneMetrics, "lane_laser_receiver_power_low_alarm", laneStats.LaneLaserReceiverPowerLowAlarm)
				addToMetrics(laneMetrics, "lane_laser_receiver_power_high_warning", laneStats.LaneLaserReceiverPowerHighWarning)
				addToMetrics(laneMetrics, "lane_laser_receiver_power_low_warning", laneStats.LaneLaserReceiverPowerLowWarning)

				addToMetrics(laneMetrics, "lane_laser_bias_current_high_alarm", laneStats.LaneLaserBiasCurrentHighAlarm)
				addToMetrics(laneMetrics, "lane_laser_bias_current_low_alarm", laneStats.LaneLaserBiasCurrentLowAlarm)
				addToMetrics(laneMetrics, "lane_laser_bias_current_high_warning", laneStats.LaneLaserBiasCurrentHighWarning)
				addToMetrics(laneMetrics, "lane_laser_bias_current_low_warning", laneStats.LaneLaserBiasCurrentLowWarning)
				addToMetrics(laneMetrics, "lane_tx_loss_of_signal_alarm", laneStats.LaneTxLossOfSignalAlarm)
				addToMetrics(laneMetrics, "lane_rx_loss_of_signal_alarm", laneStats.LaneRxLossOfSignalAlarm)
				addToMetrics(laneMetrics, "lane_tx_laser_disabled_alarm", laneStats.LaneTxLaserDisabledAlarm)

				for name, value := range laneMetrics {
					if err := app.addMetricToOutput(instance, output.JoinMetricName(opticsSeriesName, "lane_diag_stats", name), laneLabels, value, timestamp); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

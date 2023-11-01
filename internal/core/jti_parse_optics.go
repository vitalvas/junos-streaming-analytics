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
			metrics[output.JoinMetricName("diag_stats", "optics_type")] = float64(stats.GetOpticsType())
			metrics[output.JoinMetricName("diag_stats", "module_temp")] = stats.GetModuleTemp()

			metrics[output.JoinMetricName("diag_stats", "module_temp_high_alarm_threshold")] = stats.GetModuleTempHighAlarmThreshold()
			metrics[output.JoinMetricName("diag_stats", "module_temp_low_alarm_threshold")] = stats.GetModuleTempLowAlarmThreshold()
			metrics[output.JoinMetricName("diag_stats", "module_temp_high_warning_threshold")] = stats.GetModuleTempHighWarningThreshold()
			metrics[output.JoinMetricName("diag_stats", "module_temp_low_warning_threshold")] = stats.GetModuleTempLowWarningThreshold()

			metrics[output.JoinMetricName("diag_stats", "laser_output_power_high_alarm_threshold_dbm")] = stats.GetLaserOutputPowerHighAlarmThresholdDbm()
			metrics[output.JoinMetricName("diag_stats", "laser_output_power_low_alarm_threshold_dbm")] = stats.GetLaserOutputPowerLowAlarmThresholdDbm()
			metrics[output.JoinMetricName("diag_stats", "laser_output_power_high_warning_threshold_dbm")] = stats.GetLaserOutputPowerHighWarningThresholdDbm()
			metrics[output.JoinMetricName("diag_stats", "laser_output_power_low_warning_threshold_dbm")] = stats.GetLaserOutputPowerLowWarningThresholdDbm()

			metrics[output.JoinMetricName("diag_stats", "laser_rx_power_high_alarm_threshold_dbm")] = stats.GetLaserRxPowerHighAlarmThresholdDbm()
			metrics[output.JoinMetricName("diag_stats", "laser_rx_power_low_alarm_threshold_dbm")] = stats.GetLaserRxPowerLowAlarmThresholdDbm()
			metrics[output.JoinMetricName("diag_stats", "laser_rx_power_high_warning_threshold_dbm")] = stats.GetLaserRxPowerHighWarningThresholdDbm()
			metrics[output.JoinMetricName("diag_stats", "laser_rx_power_low_warning_threshold_dbm")] = stats.GetLaserRxPowerLowWarningThresholdDbm()

			metrics[output.JoinMetricName("diag_stats", "laser_bias_current_high_alarm_threshold")] = stats.GetLaserBiasCurrentHighAlarmThreshold()
			metrics[output.JoinMetricName("diag_stats", "laser_bias_current_low_alarm_threshold")] = stats.GetLaserBiasCurrentLowAlarmThreshold()
			metrics[output.JoinMetricName("diag_stats", "laser_bias_current_high_warning_threshold")] = stats.GetLaserBiasCurrentHighWarningThreshold()
			metrics[output.JoinMetricName("diag_stats", "laser_bias_current_low_warning_threshold")] = stats.GetLaserBiasCurrentLowWarningThreshold()

			metrics[output.JoinMetricName("diag_stats", "module_temp_high_alarm")] = boolToFloat64(stats.GetModuleTempHighAlarm())
			metrics[output.JoinMetricName("diag_stats", "module_temp_low_alarm")] = boolToFloat64(stats.GetModuleTempLowAlarm())
			metrics[output.JoinMetricName("diag_stats", "module_temp_high_warning")] = boolToFloat64(stats.GetModuleTempHighWarning())
			metrics[output.JoinMetricName("diag_stats", "module_temp_low_warning")] = boolToFloat64(stats.GetModuleTempLowWarning())

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

				laneMetrics := map[string]float64{
					"lane_laser_temperature":        laneStats.GetLaneLaserTemperature(),
					"lane_laser_output_power_dbm":   float64(laneStats.GetLaneLaserOutputPowerDbm()),
					"lane_laser_receiver_power_dbm": float64(laneStats.GetLaneLaserReceiverPowerDbm()),
					"lane_laser_bias_current":       laneStats.GetLaneLaserBiasCurrent(),

					"lane_laser_output_power_high_alarm":   boolToFloat64(laneStats.GetLaneLaserOutputPowerHighAlarm()),
					"lane_laser_output_power_low_alarm":    boolToFloat64(laneStats.GetLaneLaserOutputPowerLowAlarm()),
					"lane_laser_output_power_high_warning": boolToFloat64(laneStats.GetLaneLaserOutputPowerHighWarning()),
					"lane_laser_output_power_low_warning":  boolToFloat64(laneStats.GetLaneLaserOutputPowerLowWarning()),

					"lane_laser_receiver_power_high_alarm":   boolToFloat64(laneStats.GetLaneLaserReceiverPowerHighAlarm()),
					"lane_laser_receiver_power_low_alarm":    boolToFloat64(laneStats.GetLaneLaserReceiverPowerLowAlarm()),
					"lane_laser_receiver_power_high_warning": boolToFloat64(laneStats.GetLaneLaserReceiverPowerHighWarning()),
					"lane_laser_receiver_power_low_warning":  boolToFloat64(laneStats.GetLaneLaserReceiverPowerLowWarning()),

					"lane_laser_bias_current_high_alarm":   boolToFloat64(laneStats.GetLaneLaserBiasCurrentHighAlarm()),
					"lane_laser_bias_current_low_alarm":    boolToFloat64(laneStats.GetLaneLaserBiasCurrentLowAlarm()),
					"lane_laser_bias_current_high_warning": boolToFloat64(laneStats.GetLaneLaserBiasCurrentHighWarning()),
					"lane_laser_bias_current_low_warning":  boolToFloat64(laneStats.GetLaneLaserBiasCurrentLowWarning()),
					"lane_tx_loss_of_signal_alarm":         boolToFloat64(laneStats.GetLaneTxLossOfSignalAlarm()),
					"lane_rx_loss_of_signal_alarm":         boolToFloat64(laneStats.GetLaneRxLossOfSignalAlarm()),
					"lane_tx_laser_disabled_alarm":         boolToFloat64(laneStats.GetLaneTxLaserDisabledAlarm()),
				}

				for name, value := range laneMetrics {
					if err := app.addMetricToOutput(instance, output.JoinMetricName(opticsSeriesName, "optics_lane_diag_stats", name), laneLabels, value, timestamp); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

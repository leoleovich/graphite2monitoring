package main
import (
	"fmt"
)


type CompareMetrics struct {
	graphiteClient GraphiteClient
	metric string
	mode string
	range1From string
	range1Until string
	range2From string
	range2Until string
	thresholdWarningI float64
	thresholdCriticalI float64
	thresholdWarningD float64
	thresholdCriticalD float64
}

func (cm *CompareMetrics) percentageDiff(debug bool) (string, int) {
	var difference float64
	var message string

	firstMetricAVG := cm.graphiteClient.getValueOfMetric(cm.metric, cm.range1From, cm.range1Until, debug)
	secondMetricAVG := cm.graphiteClient.getValueOfMetric(cm.metric, cm.range2From, cm.range2Until, debug)

	// Percentage of the second metric from the first one
	percentage := secondMetricAVG*100.0/firstMetricAVG

	if debug {
		fmt.Printf("firstMetricAVG: %.2f\n", firstMetricAVG)
		fmt.Printf("secondMetricAVG: %.2f\n", secondMetricAVG)
		fmt.Printf("Percentage: Second range is %.2f% from the first one\n", percentage)
	}

	if percentage > 100 {
		// Growing of metric
		difference = percentage - 100
		message = fmt.Sprintf("Increasing of metric is: %.2f%", difference)
		if difference > cm.thresholdCriticalI && cm.thresholdCriticalI != MAGIC_DO_NOT_CARE_VALUE {
			return message, EXIT_CODE_CRITICAL
		} else if difference > cm.thresholdWarningI && cm.thresholdWarningI != MAGIC_DO_NOT_CARE_VALUE {
			return message, EXIT_CODE_WARNING
		}
	} else if percentage < 100 {
		// Decreasing of metric
		difference = 100 - percentage
		message = fmt.Sprintf("Decreasing of metric is: %.2f%", difference)
		if difference > cm.thresholdCriticalD && cm.thresholdCriticalD != MAGIC_DO_NOT_CARE_VALUE {
			return message, EXIT_CODE_CRITICAL
		} else if difference > cm.thresholdWarningD && cm.thresholdWarningD != MAGIC_DO_NOT_CARE_VALUE {
			return message, EXIT_CODE_WARNING
		}
	} else {
		message = "There is no difference between the time ranges"
	}
	return message, EXIT_CODE_OK
}

func (cm *CompareMetrics) absoluteDiff(debug bool) (string, int) {
	var message string
	firstMetricAVG := cm.graphiteClient.getValueOfMetric(cm.metric, cm.range1From, cm.range1Until, debug)
	secondMetricAVG := cm.graphiteClient.getValueOfMetric(cm.metric, cm.range2From, cm.range2Until, debug)

	difference := secondMetricAVG - firstMetricAVG

	if debug {
		fmt.Printf("firstMetricAVG: %.2f\n", firstMetricAVG)
		fmt.Printf("secondMetricAVG: %.2f\n", secondMetricAVG)
	}

	if difference > 0 {
		message = fmt.Sprintf("Increasing of metric is: %.2f", difference)
		if difference > cm.thresholdCriticalI && cm.thresholdCriticalI != MAGIC_DO_NOT_CARE_VALUE {
			return message, EXIT_CODE_CRITICAL
		} else if difference >= cm.thresholdWarningI && cm.thresholdWarningI != MAGIC_DO_NOT_CARE_VALUE {
			return message, EXIT_CODE_WARNING
		}
	} else if difference < 0 {
		message = fmt.Sprintf("Decreasing of metric is: %.2f", difference)
		if difference > cm.thresholdCriticalD && cm.thresholdCriticalD != MAGIC_DO_NOT_CARE_VALUE {
			return message, EXIT_CODE_CRITICAL
		} else if difference >= cm.thresholdWarningD && cm.thresholdWarningD != MAGIC_DO_NOT_CARE_VALUE {
			return message, EXIT_CODE_WARNING
		}
	} else {
		message = "There is no difference between the time ranges"
	}
	return message, EXIT_CODE_OK
}

func (cm *CompareMetrics) absoluteCmp(debug bool) (string, int) {
	value := cm.graphiteClient.getValueOfMetric(cm.metric, cm.range1From, cm.range1Until, debug)

	if value > cm.thresholdCriticalI && cm.thresholdCriticalI != MAGIC_DO_NOT_CARE_VALUE {
		return fmt.Sprintf("Metric is above critical threshold (%.2f > %.2f)", value, cm.thresholdCriticalI), EXIT_CODE_CRITICAL
	} else if value >= cm.thresholdWarningI && cm.thresholdWarningI != MAGIC_DO_NOT_CARE_VALUE {
		return fmt.Sprintf("Metric is above warning threshold (%.2f => %.2f)", value, cm.thresholdWarningI), EXIT_CODE_WARNING
	} else if value < cm.thresholdCriticalD && cm.thresholdCriticalD != MAGIC_DO_NOT_CARE_VALUE {
		return fmt.Sprintf("Metric is below critical threshold (%.2f < %.2f)", value, cm.thresholdCriticalD), EXIT_CODE_CRITICAL
	} else if value <= cm.thresholdWarningD && cm.thresholdWarningD != MAGIC_DO_NOT_CARE_VALUE {
		return fmt.Sprintf("Metric is below warning threshold (%.2f <= %.2f)", value, cm.thresholdWarningD), EXIT_CODE_WARNING
	} else {
		return fmt.Sprintf("Metric is ok (low limits (%.2f %.2f) < %.2f < high limits (%.2f %.2f))",
			cm.thresholdCriticalD, cm.thresholdWarningD, value, cm.thresholdWarningI, cm.thresholdCriticalI ), EXIT_CODE_OK
	}
}

func (cm *CompareMetrics) analysisOfMetrics(debug bool) (string, int)  {
	switch cm.mode {
	case "percentageDiff":
		return cm.percentageDiff(debug)
	case "absoluteDiff":
		return cm.absoluteDiff(debug)
	case "absoluteCmp":
		return cm.absoluteCmp(debug)
	default:
		return "Unsupported mode", EXIT_CODE_UNKNOWN
	}
}
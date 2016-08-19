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

func (cm *CompareMetrics) percentage(debug bool) (string, int) {
	var difference float64
	var message string

	firstMetricAVG := cm.graphiteClient.getValueOfMetric(cm.metric, cm.range1From, cm.range1Until, debug)
	secondMetricAVG := cm.graphiteClient.getValueOfMetric(cm.metric, cm.range2From, cm.range2Until, debug)

	// Percentage of the second metric from the first one
	percentage := secondMetricAVG*100.0/firstMetricAVG

	if debug {
		fmt.Printf("firstMetricAVG: %f\n", firstMetricAVG)
		fmt.Printf("secondMetricAVG: %f\n", secondMetricAVG)
		fmt.Printf("Percentage: Second range is %f% from the first one\n", percentage)
	}

	if percentage > 100 {
		// Growing of metric
		difference = percentage - 100
		message = fmt.Sprintf("Increasing of metric is: %f%", difference)
		if difference > float64(cm.thresholdCriticalI) && cm.thresholdCriticalI != 0 {
			return message, EXIT_CODE_CRITICAL
		} else if difference > float64(cm.thresholdWarningI) && cm.thresholdWarningI != 0 {
			return message, EXIT_CODE_WARNING
		}
	} else if percentage < 100 {
		// Decreasing of metric
		difference = 100 - percentage
		message = fmt.Sprintf("Decreasing of metric is: %f%", difference)
		if difference > float64(cm.thresholdCriticalD) && cm.thresholdCriticalD != 0 {
			return message, EXIT_CODE_CRITICAL
		} else if difference > float64(cm.thresholdWarningD) && cm.thresholdWarningD != 0 {
			return message, EXIT_CODE_WARNING
		}
	} else {
		message = "There is no difference between the time ranges"
	}
	return message, EXIT_CODE_OK
}

func (cm *CompareMetrics) absolute(debug bool) (string, int) {
	var message string
	firstMetricAVG := cm.graphiteClient.getValueOfMetric(cm.metric, cm.range1From, cm.range1Until, debug)
	secondMetricAVG := cm.graphiteClient.getValueOfMetric(cm.metric, cm.range2From, cm.range2Until, debug)

	difference := secondMetricAVG - firstMetricAVG

	if debug {
		fmt.Printf("firstMetricAVG: %f\n", firstMetricAVG)
		fmt.Printf("secondMetricAVG: %f\n", secondMetricAVG)
	}

	if difference > 0 {
		message = fmt.Sprintf("Increasing of metric is: %f", difference)
		if difference > float64(cm.thresholdCriticalI) && cm.thresholdCriticalI != 0 {
			return message, EXIT_CODE_CRITICAL
		} else if difference > float64(cm.thresholdWarningI) && cm.thresholdWarningI != 0 {
			return message, EXIT_CODE_WARNING
		}
	} else if difference < 0 {
		message = fmt.Sprintf("Decreasing of metric is: %f", difference)
		if difference > float64(cm.thresholdCriticalD) && cm.thresholdCriticalD != 0 {
			return message, EXIT_CODE_CRITICAL
		} else if difference > float64(cm.thresholdWarningD) && cm.thresholdWarningD != 0 {
			return message, EXIT_CODE_WARNING
		}
	} else {
		message = "There is no difference between the time ranges"
	}
	return message, EXIT_CODE_OK
}

func (cm *CompareMetrics) absoluteSingle(debug bool) (string, int) {
	value := cm.graphiteClient.getValueOfMetric(cm.metric, cm.range1From, cm.range1Until, debug)

	if value > cm.thresholdCriticalI && cm.thresholdCriticalI != 0 {
		return fmt.Sprintf("Metric is above critical threshold (%f > %f)", value, cm.thresholdCriticalI), EXIT_CODE_CRITICAL
	} else if value > cm.thresholdWarningI && cm.thresholdWarningI != 0 {
		return fmt.Sprintf("Metric is above warning threshold (%f > %f)", value, cm.thresholdWarningI), EXIT_CODE_WARNING
	} else if value < cm.thresholdCriticalD && cm.thresholdCriticalD != 0 {
		return fmt.Sprintf("Metric is below critical threshold (%f < %f)", value, cm.thresholdCriticalD), EXIT_CODE_CRITICAL
	} else if value > cm.thresholdWarningD && cm.thresholdWarningD != 0 {
		return fmt.Sprintf("Metric is below warning threshold (%f < %f)", value, cm.thresholdWarningD), EXIT_CODE_WARNING
	} else {
		return fmt.Sprintf("Metric is ok (low limits (%f %f) < %f < high limits (%f %f))",
			cm.thresholdCriticalD, cm.thresholdWarningD, value, cm.thresholdWarningI, cm.thresholdCriticalI ), EXIT_CODE_OK
	}
}

func (cm *CompareMetrics) analysisOfMetrics(debug bool) (string, int)  {
	switch cm.mode {
	case "percent":
		return cm.percentage(debug)
	case "absolute":
		return cm.absolute(debug)
	case "singleAbsolute":
		return cm.absoluteSingle(debug)
	default:
		return "Unsupported mode", EXIT_CODE_UNKNOWN
	}
}
package main
import (
	"strconv"
//	"fmt"
)


type CompareMetrics struct {
	graphiteClient GraphiteClient
	metric string
	range1From string
	range1Until string
	range2From string
	range2Until string
	thresholdWarning int
	thresholdCritical int
}
const EXIT_CODE_OK = 0;
const EXIT_CODE_WARNING = 1;
const EXIT_CODE_CRITICAL = 2;

func (cm CompareMetrics) compare() float64 {
	firstMetricAVG := cm.graphiteClient.getValueOfMetric(cm.metric, cm.range1From, cm.range1Until)
	secondMetricAVG := cm.graphiteClient.getValueOfMetric(cm.metric, cm.range2From, cm.range2Until)

	// Percentage of the second value from the first one
	return secondMetricAVG*100.0/firstMetricAVG
}

func (cm CompareMetrics) analysisOfMetrics() (string, int)  {
	percentage := cm.compare()
	var difference float64
	var message string
	if percentage > 100 {
		// Growing of metric
		difference = percentage - 100
		message = "Increasing of metrics is: " + strconv.FormatFloat(difference,'f', 1, 64) + "%"
	} else if percentage < 100 {
		// Decreasing of metric
		difference = 100 - percentage
		message = "Decreasing of metrics is: " + strconv.FormatFloat(difference,'f', 1, 64) + "%"
	} else {
		message = "There is no difference between metrics"
	}
	//fmt.Println(difference)

	if difference > float64(cm.thresholdCritical) {
		return message, EXIT_CODE_CRITICAL
	} else if difference > float64(cm.thresholdWarning) {
		return message, EXIT_CODE_WARNING
	} else {
		return message, EXIT_CODE_OK
	}
}
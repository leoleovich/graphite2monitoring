package main
import (
	"strconv"
//	"fmt"
	"fmt"
)


type CompareMetrics struct {
	graphiteClient GraphiteClient
	metric string
	range1From string
	range1Until string
	range2From string
	range2Until string
	thresholdWarningI int
	thresholdCriticalI int
	thresholdWarningD int
	thresholdCriticalD int
}
const EXIT_CODE_OK = 0;
const EXIT_CODE_WARNING = 1;
const EXIT_CODE_CRITICAL = 2;

func (cm CompareMetrics) compare(debug bool) float64 {
	firstMetricAVG := cm.graphiteClient.getValueOfMetric(cm.metric, cm.range1From, cm.range1Until, debug)
	secondMetricAVG := cm.graphiteClient.getValueOfMetric(cm.metric, cm.range2From, cm.range2Until, debug)

	if debug {
		fmt.Println("firstMetricAVG:" + strconv.FormatFloat((firstMetricAVG), 'f', 1, 64))
		fmt.Println("secondMetricAVG:" + strconv.FormatFloat((secondMetricAVG), 'f', 1, 64))
	}
	// Percentage of the second value from the first one
	return secondMetricAVG*100.0/firstMetricAVG
}

func (cm CompareMetrics) analysisOfMetrics(debug bool) (string, int)  {
	percentage := cm.compare(debug)
	var difference float64
	var message string

	if debug {
		fmt.Println("Percentage: Second range is " + strconv.FormatFloat(percentage,'f', 1, 64) + "% from the first one")
	}

	if percentage > 100 {
		// Growing of metric
		difference = percentage - 100
		message = "Increasing of metric is: " + strconv.FormatFloat(difference,'f', 1, 64) + "%"
		if difference > float64(cm.thresholdCriticalI) {
			return message, EXIT_CODE_CRITICAL
		} else if difference > float64(cm.thresholdWarningI) {
			return message, EXIT_CODE_WARNING
		}
	} else if percentage < 100 {
		// Decreasing of metric
		difference = 100 - percentage
		message = "Decreasing of metric is: " + strconv.FormatFloat(difference,'f', 1, 64) + "%"
		if difference > float64(cm.thresholdCriticalD) {
			return message, EXIT_CODE_CRITICAL
		} else if difference > float64(cm.thresholdWarningD) {
			return message, EXIT_CODE_WARNING
		}
	} else {
		message = "There is no difference between thetime ranges"
	}
	return message, EXIT_CODE_OK
}
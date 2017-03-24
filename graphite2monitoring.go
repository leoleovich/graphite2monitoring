package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

const EXIT_CODE_OK = 0
const EXIT_CODE_WARNING = 1
const EXIT_CODE_CRITICAL = 2
const EXIT_CODE_UNKNOWN = 3
const MAGIC_DO_NOT_CARE_VALUE = -0.10101

func converTime2GraphiteFormat(time2convert int) string {
	return "-" + strconv.Itoa(time2convert) + "s" // -3600s
}

func main() {
	var username, password, metric, url, mode string
	var thresholdWarningI, thresholdCriticalI, thresholdWarningD, thresholdCriticalD float64
	var range1FromAgo, range1UntilAgo, range2FromAgo, range2UntilAgo int
	var debug bool

	message := "This service will analyze metrics from graphite in one of next modes:\n" +
		"\t - percentageDiff (default):\n" +
		"\t\t takes mertic within range1 (from-until) and range2 (from-until)\n" +
		"\t\t and count as percent 2 metric from the 1st one (2nd*100/1st)\n" +
		"\t - absoluteDiff:\n" +
		"\t\t simular to percent, but with absolute values (2nd-1st)\n" +
		"\t - absoluteCmp:\n" +
		"\t\t tales metric only within range1 (from-until) and compares with absolute numbers (not diff)\n" +
		"All of these methods print the result and give you exit code according nagios standards"

	flag.StringVar(&username, "u", "graphite", "User, which has rights to access Graphite")
	flag.StringVar(&password, "p", "", "Password to access the graphite-API. For example 'qqq'")
	flag.StringVar(&metric, "m", "", "Name of metric or metric filter e.g. qqqq.test.leoleovich.currentProblems")
	flag.StringVar(&url, "U", "", "Base address of your graphite server e.g. https://graphite.protury.info/")

	flag.StringVar(&mode, "mode", "percentageDiff", "Mode of analysis of metrics. E.G. percentageDiff, absoluteDiff, absoluteCmp")

	flag.IntVar(&range1FromAgo, "range1From", 90000, "Amount of seconds ago for the 1st range (from)")
	flag.IntVar(&range1UntilAgo, "range1Until", 86400, "Amount of seconds ago for the 1st range (until)")
	flag.IntVar(&range2FromAgo, "range2From", 3600, "Amount of seconds ago for the 2st range (from)")
	flag.IntVar(&range2UntilAgo, "range2Until", 0, "Amount of seconds ago for the 2st range (until)")

	flag.Float64Var(&thresholdWarningI, "wi", MAGIC_DO_NOT_CARE_VALUE, "Increasing. Metrics above this threshold will be marked as warning")
	flag.Float64Var(&thresholdCriticalI, "ci", MAGIC_DO_NOT_CARE_VALUE, "Increasing. Metrics above this threshold will be marked as critical")
	flag.Float64Var(&thresholdWarningD, "wd", MAGIC_DO_NOT_CARE_VALUE, "Decreasing. Metrics below this threshold will be marked as warning")
	flag.Float64Var(&thresholdCriticalD, "cd", MAGIC_DO_NOT_CARE_VALUE, "Decreasing. Metrics below this threshold will be marked as critical")

	flag.BoolVar(&debug, "d", false, "Debug mode will print a lot of additinal info")
	flag.Parse()

	if metric == "" || url == "" {
		fmt.Println("URL (-U) and metric (-m) attributes are required")
		flag.Usage()
		os.Exit(EXIT_CODE_UNKNOWN)
	}

	switch mode {
	case "percentageDiff", "absoluteDiff", "absoluteCmp":
	default:
		fmt.Println("Unsupported mode " + mode)
		flag.Usage = func() {
			fmt.Println(message)
		}
		os.Exit(EXIT_CODE_UNKNOWN)
	}

	range1FromS := converTime2GraphiteFormat(range1FromAgo)
	range1UntilS := converTime2GraphiteFormat(range1UntilAgo)
	range2FromS := converTime2GraphiteFormat(range2FromAgo)
	range2UntilS := converTime2GraphiteFormat(range2UntilAgo)

	if debug {
		fmt.Println("Range1: " + range1FromS + " - " + range1UntilS)
		fmt.Println("Range2: " + range2FromS + " - " + range2UntilS)
	}

	cm := CompareMetrics{
		GraphiteClient{username, password, url + "/render?"},
		metric,
		mode,
		range1FromS,
		range1UntilS,
		range2FromS,
		range2UntilS,
		thresholdWarningI,
		thresholdCriticalI,
		thresholdWarningD,
		thresholdCriticalD}

	// Compare metrics and return result
	result, returnCode := cm.analysisOfMetrics(debug)

	// Print and exit
	fmt.Println(result)
	os.Exit(returnCode)
}

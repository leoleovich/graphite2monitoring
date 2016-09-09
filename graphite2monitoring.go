package main

import (
	"fmt"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"os"
	"strconv"
)

const EXIT_CODE_OK = 0;
const EXIT_CODE_WARNING = 1;
const EXIT_CODE_CRITICAL = 2;
const EXIT_CODE_UNKNOWN = 3;
const MAGIC_DO_NOT_CARE_VALUE = -0.10101;

func new_token(username string, timestamp int, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	unauthed_token := fmt.Sprint(timestamp,":",username)
	h.Write([]byte(unauthed_token))
	return fmt.Sprint(hex.EncodeToString(h.Sum(nil)), ":", timestamp, ":", username)
}
func converTime2GraphiteFormat(time2convert int) string {
	return "-" + strconv.Itoa(time2convert) + "s" // -3600s
}

func main() {
	var username, authToken, metric, url, mode string
	var thresholdWarningI, thresholdCriticalI, thresholdWarningD, thresholdCriticalD float64
	var range1FromAgo, range1UntilAgo, range2FromAgo, range2UntilAgo int
	var debug bool


	message := "This service will analyze metrics from graphite in one of next modes:\n" +
	"\t - percent (default):\n" +
	"\t\t takes mertic within range1 (from-until) and range2 (from-until)\n" +
	"\t\t and count as percent 2 metric from the 1st one (2nd*100/1st)\n" +
	"\t - absolute:\n" +
	"\t\t simular to percent, but with absolute values (2nd-1st)\n" +
	"\t - absoluteSingle:\n" +
	"\t\t tales metric only within range1 (from-until)\n" +
	"All of these methods print the result and give you exit code according nagios standards"


	flag.StringVar(&username, "u", "graphite", "User, which has rights to access Graphite")
	flag.StringVar(&authToken, "a", "", "AuthToken to access the graphite-API. For example 'qqq'")
	flag.StringVar(&metric, "m", "", "Name of metric or metric filter e.g. qqqq.test.leoleovich.currentProblems")
	flag.StringVar(&url, "U", "", "Base address of your graphite server e.g. https://graphite.protury.info/")

	flag.StringVar(&mode, "mode", "percent", "Mode of analysis of metrics. E.G. percent, absolute, singleAbsolute")

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

	if authToken == "" ||  metric == "" || url == "" {
		fmt.Println("URL (-U), authToken (-a) and metric (-m) attributes are required")
		flag.Usage()
		os.Exit(EXIT_CODE_UNKNOWN)
	}

	switch mode {
	case "percent":
	case "absolute":
	case "singleAbsolute":
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
		fmt.Println("Token: " + authToken)
		fmt.Println("Range1: " + range1FromS + " - " + range1UntilS)
		fmt.Println("Range2: " + range2FromS + " - " + range2UntilS)
	}

	cm := CompareMetrics{
		GraphiteClient{new_token(username, int(time.Now().Unix()), authToken), url},
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

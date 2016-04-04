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
	var username, authToken, metric, url string
	var thresholdWarningI, thresholdCriticalI, thresholdWarningD, thresholdCriticalD int
	var range1FromAgo, range1UntilAgo, range2FromAgo, range2UntilAgo int
	var debug bool
	flag.StringVar(&username, "u", "graphite", "User, which has rights to access Graphite")
	flag.StringVar(&authToken, "a", "", "AuthToken to access the graphite-API. For example 'qqq'")
	flag.StringVar(&metric, "m", "", "Name of metric or metric filter e.g. qqqq.test.leoleovich.currentProblems")
	flag.StringVar(&url, "U", "", "Base address of your graphite server e.g. https://graphite.protury.info/")

	flag.IntVar(&range1FromAgo, "range1From", 90000, "Amount of seconds ago for the 1st range (from)")
	flag.IntVar(&range1UntilAgo, "range1Until", 86400, "Amount of seconds ago for the 1st range (until)")
	flag.IntVar(&range2FromAgo, "range2From", 3600, "Amount of seconds ago for the 2st range (from)")
	flag.IntVar(&range2UntilAgo, "range2Until", 0, "Amount of seconds ago for the 2st range (until)")

	flag.IntVar(&thresholdWarningI, "wi", 20, "Metrics above this threshold will be marked as warning")
	flag.IntVar(&thresholdCriticalI, "ci", 40, "Metrics above this threshold will be marked as critical")
	flag.IntVar(&thresholdWarningD, "wd", 20, "Metrics below this threshold will be marked as warning")
	flag.IntVar(&thresholdCriticalD, "cd", 40, "Metrics below this threshold will be marked as critical")

	flag.BoolVar(&debug, "d", false, "Debug mode will print a lot of additinal info")
	flag.Parse()

	if authToken == "" ||  metric == "" || url == "" {
		fmt.Println("URL (-U), authToken (-a) and metric (-m) attributes are required")
		os.Exit(5)
	}
	if thresholdCriticalD < thresholdWarningD || thresholdCriticalI < thresholdWarningI {
		fmt.Println("Critical threshold can not be less, than warning")
		os.Exit(5)
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

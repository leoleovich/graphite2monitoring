package main

import (
	"fmt"
	"time"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"os"
	"strings"
)



func new_token(username string, timestamp int, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	unauthed_token := fmt.Sprint(timestamp,":",username)
	h.Write([]byte(unauthed_token))
	return fmt.Sprint(hex.EncodeToString(h.Sum(nil)), ":", timestamp, ":", username)
}
func converTime2GraphiteFormat(time2convert int64) string {
	timeRange := time.Unix(time.Now().Unix() - time2convert, 0)
	return strings.Replace(timeRange.Format("15:04 20060102")," ","_", -1) //10:0020150923
}

func main() {
	var username, authToken, metric string
	var thresholdWarningI, thresholdCriticalI, thresholdWarningD, thresholdCriticalD int
	var range1FromAgo, range1UntilAgo, range2FromAgo, range2UntilAgo int64
	var debug bool
	flag.StringVar(&username, "u", "graphite", "User, which has rights to access Graphite")
	flag.StringVar(&authToken, "a", "", "AuthToken to access the graphite-API. For example 'qqq'")
	flag.StringVar(&metric, "m", "", "Name of metric or metric filter e.g. qqqq.test.leoleovich.currentProblems")

	flag.Int64Var(&range1FromAgo, "range1From", 90000, "Amount of seconds ago for the 1st range (from)")
	flag.Int64Var(&range1UntilAgo, "range1Until", 86400, "Amount of seconds ago for the 1st range (until)")
	flag.Int64Var(&range2FromAgo, "range2From", 3600, "Amount of seconds ago for the 2st range (from)")
	flag.Int64Var(&range2UntilAgo, "range2Until", 0, "Amount of seconds ago for the 2st range (until)")

	flag.IntVar(&thresholdWarningI, "wi", 20, "Metrics above this threshold will be marked as warning")
	flag.IntVar(&thresholdCriticalI, "ci", 40, "Metrics above this threshold will be marked as critical")
	flag.IntVar(&thresholdWarningD, "wd", 20, "Metrics below this threshold will be marked as warning")
	flag.IntVar(&thresholdCriticalD, "cd", 40, "Metrics below this threshold will be marked as critical")

	flag.BoolVar(&debug, "d", false, "Debug mode will print a lot of additinal info")
	flag.Parse()

	if authToken == "" ||  metric == "" {
		fmt.Println("authToken (-a) and metric (-m) attributes are required")
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
		GraphiteClient{new_token(username, int(time.Now().Unix()), authToken)},
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

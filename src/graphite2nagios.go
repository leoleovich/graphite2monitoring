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
func converTime2GraphiteFormat(time2convert string) string {
	timeRange,_ := time.Parse("2006-01-02 15:04", time2convert)
	return strings.Replace(timeRange.Format("15:04 20060102")," ","_", -1) //10:0020150923
}


func main() {
	var username, authToken, metric, range1From, range1Until, range2From, range2Until string
	var thresholdWarning, thresholdCritical int
	var debug bool
        flag.StringVar(&username, "u", "graphite", "User, which has rights to access Graphite")
	flag.StringVar(&authToken, "a", "", "AuthToken to access the graphite-API. For example 'qqq'")
	flag.StringVar(&metric, "metric", "qqqq.test.leoleovich.currentProblems", "Name of metric or metric filter e.g. Character.*")
	flag.StringVar(&range1From, "range1From", time.Unix((time.Now().Unix() - 90000), 0).Format("2006-01-02 15:04"), "e.g. 2014-09-01 10:00")
	flag.StringVar(&range1Until, "range1Until", time.Unix((time.Now().Unix() - 86400), 0).Format("2006-01-02 15:04"), "e.g. 2014-09-01 11:00")
	flag.StringVar(&range2From, "range2From", time.Unix((time.Now().Unix() - 3600), 0).Format("2006-01-02 15:04"), "e.g. 2014-09-01 10:00")
	flag.StringVar(&range2Until, "range2Until", time.Now().Format("2006-01-02 15:04"), "e.g. 2014-09-01 11:00")
	flag.IntVar(&thresholdWarning, "w", 20, "Metrics above the threshold will be marked as warning")
	flag.IntVar(&thresholdCritical, "c", 40, "Metrics above the threshold will be marked as critical")
	flag.BoolVar(&debug, "d", false, "Debug mode will print a lot of additinal info")
	flag.Parse()

	if authToken == "" {
		println("authToken (-a) attribute is required")
		os.Exit(2)
	}
	range1From = converTime2GraphiteFormat(range1From)
	range1Until = converTime2GraphiteFormat(range1Until)
	range2From = converTime2GraphiteFormat(range2From)
	range2Until = converTime2GraphiteFormat(range2Until)

	if debug {
		fmt.Println(authToken)
		fmt.Println(range1From+" - "+range1Until+"\n")
		fmt.Println(range2From+" - "+range2Until+"\n")
	}


	cm := CompareMetrics{
		GraphiteClient{new_token(username, int(time.Now().Unix()), authToken)},
		metric,
		range1From,
		range1Until,
		range2From,
		range2Until,
		thresholdWarning,
		thresholdCritical}

	result, returnCode := cm.analysisOfMetrics()
	println(result)
	os.Exit(returnCode)
}

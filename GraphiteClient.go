package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"crypto/tls"
)

type GraphiteClient struct {
	Username string
	Password string
	URL      string
}

func (c GraphiteClient) download(debug bool) []byte {
	if debug {
		fmt.Println("URL: " + c.URL)
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := http.Client{Timeout: time.Duration(5 * time.Second), Transport: tr}
	req, err := http.NewRequest("POST", c.URL, nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}

	req.SetBasicAuth(c.Username, c.Password)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	if res.StatusCode != 200 {
		fmt.Println("Graphite returned", res.StatusCode)
		os.Exit(3)
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
	defer res.Body.Close()
	return content
}

func json2float(content []byte) float64 {
	type Result struct {
		Target     string      `json:"target"`
		Datapoints [][]float64 `json:"datapoints"`
	}
	var r []Result
	//s := `[{"target": "summarize(qqqq.test.leoleovich.currentProblems, \"99year\", \"avg\")", "datapoints": [[0.5317615073466304, 0]]}]`
	json.Unmarshal(content, &r)
	//fmt.Println(r[0].Datapoints[0][0])
	if len(r) == 0 {
		fmt.Println("Invalid data got from Graphite - Check, that your metric exist")
		os.Exit(3)
	}

	var sum float64
	for _, value := range r {
		if len(value.Datapoints) == 0 {
			fmt.Println("Invalid data got from Graphite - Empty Datapoint set")
			os.Exit(3)
		}
		// Graphite may return null as a value (no data in fact, but you will have 0 in structure)
		sum += value.Datapoints[0][0]
	}
	return sum / float64(len(r))
}
func (c GraphiteClient) getValueOfMetric(metricName string, from string, until string, debug bool) float64 {
	//Big value, so that we get the avg over all data. Then we are sure that we get only one result.
	intervalString := "99year"
	c.URL = fmt.Sprint(c.URL, "target=summarize(", metricName, ",\"", intervalString, "\",\"avg\")&from=", from, "&until=", until, "&format=json")
	return json2float(c.download(debug))
}

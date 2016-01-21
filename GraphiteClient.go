package main
import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"
)


type GraphiteClient struct {
	authToken string
	graphiteBaseURL string
}


func (c GraphiteClient)getFindMetricBaseUrl() string {
	return fmt.Sprint(c.graphiteBaseURL, "metrics/find?", getUrlAuthTokenParam(c), "&")
}

func (c GraphiteClient)getRenderBaseUrl() string {
	return fmt.Sprint(c.graphiteBaseURL, "render?", getUrlAuthTokenParam(c), "&")
}

func getUrlAuthTokenParam(c GraphiteClient) string {
	return fmt.Sprint("__auth_token=", c.authToken)
}

func download(url string, debug bool) []byte {
	if debug {
		fmt.Println("URL: " + url)
	}
	res, err  := http.Get(url)
	if err != nil {
		panic(err)
	}
	if res.StatusCode != 200 {
		fmt.Println("Graphite returned " + strconv.Itoa(res.StatusCode))
		os.Exit(10)
	}
	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
		fmt.Println(string(content))
	}
	defer res.Body.Close()
	return content
}

func json2float (content []byte) float64 {
	type Result struct {
		Target  string `json:"target"`
		Datapoints [][]float64 `json:"datapoints"`
	}
	var r []Result
	//s := `[{"target": "summarize(qqqq.test.leoleovich.currentProblems, \"99year\", \"avg\")", "datapoints": [[0.5317615073466304, 0]]}]`
	json.Unmarshal(content, &r)
	//fmt.Println(r[0].Datapoints[0][0])
	if len(r) == 0 {
		fmt.Println("Invalid data got from Graphite - Check, that your metric exist")
		os.Exit(2)
	}

	var sum float64
	for _, value := range r {
		if len(value.Datapoints) == 0 {
			fmt.Println("Invalid data got from Graphite - Empty Datapoint set")
			os.Exit(2)
		}
		// Graphite may return null as a value (no data in fact, but you will have 0 in structure)
		sum += value.Datapoints[0][0]
	}
	return sum/float64(len(r))
}
func (c GraphiteClient)getValueOfMetric(metricName string, from string, until string, debug bool) float64 {
	//Big value, so that we get the avg over all data. Then we are sure that we get only one result.
	intervalString := "99year"
	url := fmt.Sprint(c.getRenderBaseUrl(), "target=summarize(", metricName,  ",\"",  intervalString,  "\",\"avg\")&from=", from, "&until=", until, "&format=json")

	return json2float(download(url, debug))
}
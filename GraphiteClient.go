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
}

const GRAPHITE_BASE_URL = "https://graphite.innogames.de/"


func (c GraphiteClient)getFindMetricBaseUrl() string {
	return fmt.Sprint(GRAPHITE_BASE_URL, "metrics/find?", getUrlAuthTokenParam(c), "&")
}

func (c GraphiteClient)getRenderBaseUrl() string {
	return fmt.Sprint(GRAPHITE_BASE_URL, "render?", getUrlAuthTokenParam(c), "&")
}

func getUrlAuthTokenParam(c GraphiteClient) string {
	return fmt.Sprint("__auth_token=", c.authToken)
}

func download(url string) []byte {
	fmt.Print(url+"\n")
	res, err  := http.Get(url)
	//fmt.Print(res, err)
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
	//s := `[{"target": "summarize(nagios.wallOfShame.oleg_obleukhov.currentProblems, \"99year\", \"avg\")", "datapoints": [[0.5317615073466304, 0]]}]`
	json.Unmarshal(content, &r)
	//fmt.Println(r[0].Datapoints[0][0])
	if len(r) == 0 {
		fmt.Println("Invalid data got from Graphite - Check, that your metric exist")
		os.Exit(2)
	} else if len(r[0].Datapoints) == 0 {
		fmt.Println("Invalid data got from Graphite - Empty Datapoint set")
		os.Exit(2)
	}
	return r[0].Datapoints[0][0]
}
func (c GraphiteClient)getValueOfMetric(metricName string, from string, until string) float64 {
	intervalString := "99year"; //Big value, so that we get the avg over all data. Then we are sure that we get only one result.
	url := fmt.Sprint(c.getRenderBaseUrl(), "target=summarize(", metricName,  ",\"",  intervalString,  "\",\"avg\")&from=", from, "&until=", until, "&format=json")

	return json2float(download(url))
}
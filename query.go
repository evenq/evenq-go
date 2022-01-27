package evenq

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// QueryRequest event data
type QueryRequest struct {
	EventID       string      `json:"id"`
	PartitionKeys []string    `json:"partitionKeys"`
	From          time.Time   `json:"from"`
	To            time.Time   `json:"to"`
	Items         []Item      `json:"items"`
	Conditions    []Condition `json:"conditions"`
}

// ItemType what to query for
type ItemType string

const (
	TypeList       ItemType = "list"
	TypeNumber     ItemType = "number"
	TypeTimeseries ItemType = "timeseries"
)

// Agg how to aggregate numbers
type Agg string

const (
	AggAvg           Agg = "avg"
	AggSum           Agg = "sum"
	AggSumCumulative Agg = "sum_cumulative"
	AggCount         Agg = "count"
	AggCountUnique   Agg = "count_unique"
	AggMin           Agg = "min"
	AggMax           Agg = "max"
)

// Item within query
type Item struct {
	Key         string   `json:"key,omitempty"`             // the index of the column in the data
	Type        ItemType `json:"type"`                      // list, number, timeseries
	Aggregation Agg      `json:"aggregation"`               // avg, sum, sum_cumulative, count, count_unique, min, max
	Interval    string   `json:"interval,omitempty"`        // (type=timeseries) 5m, 4h
	ListOrder   string   `json:"listOrder,omitempty"`       // (type=list) asc, desc
	ListLimit   int      `json:"listLimit,omitempty"`       // amount of items to be returned in list
	ExcludeNull bool     `json:"listExcludeNull,omitempty"` // whether or not to count null values
}

type Condition struct {
	Key       string      `json:"key"`
	Operation string      `json:"op"`
	Value     interface{} `json:"value"`
}

type Result struct {
	Query   *QueryRequest        `json:"query"`
	Stats   *Stats               `json:"stats"`
	Results *map[int]interface{} `json:"results"`
	Error   *string              `json:"error"`
}

type Stats struct {
	EventsProcessed        int     `json:"eventsProcessed"`
	EventsAnalyzed         int     `json:"eventsAnalyzed"`
	MillionEventsPerSecond float64 `json:"meps"`
	DurationTotal          int64   `json:"durationTotal"`
}

func Query(data []QueryRequest) ([]Result, error) {
	var out []Result

	if clientKey == nil || *clientKey == "" {
		return out, errors.New("[ERROR] missing api key, please run evenq.Init() first")
	}

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(data)
	if err != nil {
		return out, err
	}

	req, err := http.NewRequest("POST", queryEndpoint, &buf)
	if err != nil {
		if isVerbose {
			log.Println("evenq", err.Error())
		}
		return out, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(authHeader, *clientKey)

	resp, err := httpQueryClient.Do(req)
	if err != nil {
		return out, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)

		return out, fmt.Errorf("http statusCode: %v with message %v", resp.StatusCode, string(body))
	}

	err = json.NewDecoder(resp.Body).Decode(&out)

	return out, err
}

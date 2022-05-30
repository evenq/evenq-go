package evenq

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const (
	authHeader     = "x-evenq-key"
	ingestEndpoint = "https://api.evenq.io/v1/events"
	queryEndpoint  = "https://api.evenq.io/v1/queries"
)

var clientKey *string
var isEnabled bool
var isVerbose bool

// Options for the evenq data client
type Options struct {
	ApiKey       string
	MaxBatchSize int
	MaxBatchWait time.Duration
	BatchWorkers int
	Verbose      bool
}

// event holds a single event
// ready to be sent to the evenq server
type event struct {
	Name         string                 `json:"name"`
	PartitionKey *string                `json:"partitionKey,omitempty"`
	TS           *time.Time             `json:"ts,omitempty"`
	Data         map[string]interface{} `json:"data"`
}

type Data map[string]interface{}

var batcher *batch

var httpClient *http.Client
var httpQueryClient *http.Client

func Init(opt Options) {
	clientKey = &opt.ApiKey
	isEnabled = true
	isVerbose = opt.Verbose

	batcher = newBatcher(opt.MaxBatchSize, opt.MaxBatchWait, opt.BatchWorkers)

	transport := &http.Transport{
		MaxIdleConnsPerHost: opt.BatchWorkers,
	}

	httpClient = &http.Client{
		Transport: transport,
		Timeout:   time.Second * 10,
	}

	httpQueryClient = &http.Client{
		Transport: transport,
		Timeout:   time.Minute * 5,
	}
}

// SetEnabled allows you to enable or disable the client to
func SetEnabled(enable bool) {
	isEnabled = enable
}

// Event sends an event with a name and json compatible data
func Event(name string, data Data) {
	if isEnabled {
		batcher.add(event{
			Name: name,
			Data: data,
		})
	}
}

// EventAt sends an event with a name at a custom timestamp
func EventAt(name string, ts time.Time, data Data) {
	if isEnabled {
		batcher.add(event{
			Name: name,
			TS:   &ts,
			Data: data,
		})
	}
}

// PartitionedEvent sends an event with a partition key
func PartitionedEvent(name string, partitionKey string, data Data) {
	if isEnabled {
		batcher.add(event{
			Name:         name,
			PartitionKey: &partitionKey,
			Data:         data,
		})
	}
}

// PartitionedEventAt sends an event with a partition key and a custom timestamp
func PartitionedEventAt(name string, partitionKey string, ts time.Time, data Data) {
	if isEnabled {
		batcher.add(event{
			Name:         name,
			PartitionKey: &partitionKey,
			TS:           &ts,
			Data:         data,
		})
	}
}

// FlushEvents processes any events left in the queue
// and sends them to the server.
func FlushEvents() {
	batcher.flush()
}

// FlushEventsSync processes any events left in the queue
// and sends them to the server. This function
// is blocking until events are sent to ensure that the system
// doesn't shut down before then
func FlushEventsSync() {
	// flush any pending items to process them
	batcher.flush()
	// wait for all processing and network to finish
	batcher.wait()
}

func processBatch(ee []event, wg *sync.WaitGroup) {
	sendBatch(ee)

	// let waitgroup know we're done with these events
	for i := 0; i < len(ee); i++ {
		wg.Done()
	}
}

func sendBatch(ee []event) {
	if clientKey == nil || *clientKey == "" {
		log.Println("[ERROR] missing api key, please run evenq.Init() first")
		return
	}

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	err := enc.Encode(ee)
	if err != nil {
		if isVerbose {
			log.Println("evenq", err.Error())
		}
		return
	}

	req, err := http.NewRequest("POST", ingestEndpoint, &buf)
	if err != nil {
		if isVerbose {
			log.Println("evenq", err.Error())
		}
		return
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(authHeader, *clientKey)

	start := time.Now()

	resp, err := httpClient.Do(req)
	if err != nil {
		if isVerbose {
			log.Println("evenq", err.Error())
		}
		return
	}

	if resp == nil || resp.Body == nil {
		if isVerbose {
			log.Println("evenq missing http response")
		}
		return
	}

	if isVerbose {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf(
			`[INFO] Request Finished in %vms. endpoint: "%v" status: %v body: "%v"`,
			time.Since(start).Milliseconds(),
			req.URL.String(),
			resp.StatusCode,
			string(body),
		)
	} else {
		// Read the body even if we don't intend to use it. Otherwise golang won't pool the connection.
		// See also: http://stackoverflow.com/questions/17948827/reusing-http-connections-in-golang/17953506#17953506
		io.Copy(ioutil.Discard, resp.Body)
	}

	resp.Body.Close()
}

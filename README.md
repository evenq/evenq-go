[![GoDoc](https://godoc.org/github.com/evenq/evenq-go?status.svg)](https://godoc.org/github.com/evenq/evenq-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/evenq/evenq-go)](https://goreportcard.com/report/github.com/evenq/evenq-go)
## The official Evenq client for go
#### Installation
`go get github.com/evenq/evenq-go`

#### Full Example
```
package main

import "github.com/evenq/evenq-go"

func init() {
  // Initialize the Evenq client with your API key.
  evenq.Init(evenq.Options{
    ApiKey: "YOUR_API_KEY"
  })
}

func main() {
  // Send an event with any JSON data
  evenq.Event("hello.world", evenq.Data{
    "string": "hello world",
    "number": 42,
    "boolean": true
  })
  
  // Query your data. Keep in mind that if you just sent your first event,
  // it can take up to a minute until you can query your data
  res, err := evenq.Query([]QueryRequest{
    {
        ID: "test.event",
        From:    time.Now().UTC().Add(-time.Hour * 24),
        To:      time.Now().UTC(),
        Items: []Item{
            {
                Type:        TypeNumber,
                Aggregation: AggCount,
            },
        },
    },
  })
  
  if err != nil {
    panic(err)
  }
    
  js, _ := json.MarshalIndent(res, "", " ")
  log.Println(string(js))
}
```

#### Setup
Initialize with your API key before sending events. You'll only have to do this once in your app lifecycle.
```
evenq.Init(evenq.Options{
    ApiKey:       "YOUR_API_KEY",
    MaxBatchSize: 500,   
    MaxBatchWait: 5,    
    BatchWorkers: 1,     
    Verbose:      true, 
})
```


### Send events
Send a single event with any data in `map[string]interface{}` format.
```
// event timestamped to current time
evenq.Event("your.event", evenq.Data{})

// event with custom timestamp
evenq.EventAt("your.event", time.Now(), evenq.Data{})
```

... Or send an event with a partition key if you want to split up your data.
```
// event with partition key timestamped to current time
evenq.PartitionedEvent("eventName", "somePartition", evenq.Data)

// event with partition key and custom timestamp
evenq.PartitionedEventAt("eventName", "somePartition", time.Now(), evenq.Data)
```

For more info on naming conventions and examples check out our docs at https://app.evenq.io/docs

### Query Events
You can query your data easily with your query helper. 
```
  // Query your data. Keep in mind that if you just sent your first event,
  // it can take up to a minute until you can query your data
  res, err := evenq.Query([]QueryRequest{
    {
        ID: "test.event",
        From:    time.Now().UTC().Add(-time.Hour * 24),
        To:      time.Now().UTC(),
        Items: []Item{
            {
                Type:        TypeNumber,
                Aggregation: AggCount,
            },
        },
    },
  })
  
  if err != nil {
    panic(err)
  }
    
  js, _ := json.MarshalIndent(res, "", " ")
  log.Println(string(js))
}
```

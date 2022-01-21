[![GoDoc](https://godoc.org/github.com/evenq/evenq-go?status.svg)](https://godoc.org/github.com/evenq/evenq-go)
## The official Evenq client for go
#### Installation
`go get github.com/evenq/evenq-go`

#### Full Example
```
package main

import "github.com/evenq/evenq-go"

func init() {
  evenq.Init(evenq.Options{
    ApiKey: "YOUR_API_KEY"
  })
}

func main() {
  evenq.Event("hello.world", evenq.Data{
    "string": "hello world",
    "number": 42,
    "boolean": true
  })
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


#### Send events
Send a single event with any data in `map[string]interface{}` format.
```
evenq.Event("your.event", evenq.Data)
```

... Or send an event with a partition key if you want to split up your data.
```
evenq.EventPartition("eventName", "somePartition", evenq.Data)
```

And that's it!

For more info on naming conventions and examples check out our docs at https://app.evenq.io/docs

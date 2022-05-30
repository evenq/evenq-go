package evenq

import (
	"os"
	"testing"
)

func TestInit(t *testing.T) {
	Init(Options{
		ApiKey:  os.Getenv("EVENQ_TESTKEY"),
		Verbose: true,
	})
}

func TestFlush(t *testing.T) {
	Init(Options{
		ApiKey:       os.Getenv("EVENQ_TESTKEY"),
		Verbose:      true,
		MaxBatchSize: 10000,
		MaxBatchWait: 60,
	})

	Event("test.event", Data{"string": "strValue", "number": 123, "bool": true})
	Event("test.event", Data{"string": "strValue", "number": 123, "bool": true})
	Event("test.event", Data{"string": "strValue", "number": 123, "bool": true})
	Event("test.event", Data{"string": "strValue", "number": 123, "bool": true})
	FlushEventsSync()
	Event("test.event", Data{"string": "strValue", "number": 123, "bool": true})
	Event("test.event", Data{"string": "strValue", "number": 123, "bool": true})
	Event("test.event", Data{"string": "strValue", "number": 123, "bool": true})

	FlushEventsSync()
	FlushEventsSync()
}

func TestEvent(t *testing.T) {
	Event("testEvent", Data{"string": "lalala", "number": 123, "bool": true})
}

func TestPartitionedEvent(t *testing.T) {
	Init(Options{
		ApiKey:  os.Getenv("EVENQ_TESTKEY"),
		Verbose: true,
	})

	PartitionedEvent("testEvent", "pk-b", Data{"string": "lalala", "number": 123, "bool": true})
	PartitionedEvent("testEvent", "pk-b", Data{"string": "lalala", "number": 123, "bool": true})
	PartitionedEvent("testEvent", "pk-b", Data{"string": "lalala", "number": 123, "bool": true})
	PartitionedEvent("testEvent", "pk-b", Data{"string": "lalala", "number": 123, "bool": true})
	PartitionedEvent("testEvent", "pk-b", Data{"string": "lalala", "number": 123, "bool": true})
	FlushEventsSync()
}

func TestEnable(t *testing.T) {
	// we disable event sending
	SetEnabled(false)
	// the api key is not inizalized
	// but that shouldn't matter since
	// the sending is disabled
	Event("testEvent", Data{"string": "lalala", "number": 123, "bool": true})
}

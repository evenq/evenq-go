package evenq

import (
	"encoding/json"
	"log"
	"os"
	"testing"
	"time"
)

func TestQuery(t *testing.T) {
	Init(Options{
		ApiKey: os.Getenv("EVENQ_TESTKEY"),
	})

	res, err := Query([]QueryRequest{
		{
			ID:   "test.event",
			From: time.Now().UTC().Add(-time.Hour * 24),
			To:   time.Now().UTC(),
			Items: []Item{
				{
					Type:        TypeNumber,
					Aggregation: AggCount,
				},
			},
		},
	})

	if err != nil {
		t.Error(err)
	}

	js, _ := json.MarshalIndent(res, "", " ")
	log.Println(string(js))
}

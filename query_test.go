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
			EventID: "test.event",
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

	js, _ := json.MarshalIndent(res, "", " ")

	log.Println(err, string(js))
}

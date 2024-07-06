package main

import (
	"dfss/cmd/es/entities"
	"dfss/cmd/es/store"
	"fmt"
	"github.com/google/uuid"
	"log"
	"strings"
	"time"
)

func main() {
	log.SetFlags(0)

	es, err := store.NewStore("test")
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	// 1. Get cluster info
	es.Info()

	//for _, e := range testData() {
	//	if err = es.Create(&e); err != nil {
	//		log.Fatalln(err)
	//	}
	//}

	data, err := es.Search("info")
	if err != nil {
		log.Fatalf("Error searching index: %s", err)
	}

	for _, d := range data.Hits {
		fmt.Println(d)
	}

	log.Println(strings.Repeat("=", 37))
}

func testData() []entities.Event {
	events := []entities.Event{
		{
			ID:        uuid.UUID{},
			Event:     "Test Alert",
			Severity:  "Info",
			CreatedTs: time.Now().Unix(),
			ClosedTs:  0,
		}, {
			ID:        uuid.UUID{},
			Event:     "RAID 5, SSD degraded",
			Severity:  "High",
			CreatedTs: time.Now().Unix() - 500,
			ClosedTs:  0,
		},
	}
	for _, ev := range events {
		u, _ := uuid.NewV7()

		ev.ID = u
	}
	return events
}

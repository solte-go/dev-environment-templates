```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

func main() {
	mechanism := plain.Mechanism{
		Username: "alice",        // replace <username> with actual
		Password: "alice-secret", // replace <password> with actual
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"localhost:29092"}, // replace with actual brokers
		Topic:     "test-topic",                // replace with actual topic
		Dialer:    dialer,
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})

	defer r.Close()
	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("could not read message: %v\n", err)
		}
		fmt.Printf("message received: %s = %s\n", string(m.Key), string(m.Value))
	}
}
```
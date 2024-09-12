package main

import (
	"context"
	"fmt"
	"time"

	"log"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

const (
	topic    = "test-topic"      // replace with your topic
	broker   = "localhost:29092" // replace with your broker
	username = "admin"           // replace with your username
	password = "admin-secret"    // replace with your password
)

//TODO CHECK IF IT POSSIBLE TO WORK WITHOUT SCRAM

func main() {

	mechanism, err := scram.Mechanism(scram.SHA512, username, password)
	if err != nil {
		log.Panicln(err)
		return
	}

	writer := &kafka.Writer{
		Addr:         kafka.TCP(broker),
		Topic:        topic,
		RequiredAcks: kafka.RequireOne,
		BatchTimeout: 100 * time.Millisecond,
		BatchSize:    1000,
		Balancer:     &kafka.LeastBytes{},
		Transport: &kafka.Transport{
			SASL: mechanism,
		},
	}

	dialer := &kafka.Dialer{
		DualStack:     true,
		SASLMechanism: mechanism,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:         []string{broker}, // replace with your broker
		Topic:           topic,            // replace with your topic
		Dialer:          dialer,
		ReadLagInterval: 0,
		Partition:       0,
		MinBytes:        10e3, // 10KB
		MaxBytes:        10e6, // 10MB
	})

	defer reader.Close()

	log.Printf("Producing message into topic: %v", topic)

	var msgs []kafka.Message

	for i := 0; i < 100; i++ {
		msgs = append(
			msgs,
			kafka.Message{
				Key:   []byte("Key"),
				Value: []byte(fmt.Sprintf("Hello World number %d!", i)),
			},
		)
	}

	err = writer.WriteMessages(context.Background(), msgs...)
	if err != nil {
		log.Fatal("failed to write messages:", err)
	}

	writer.Close()

	log.Printf("Comsuming message from topic: %v", topic)
	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			fmt.Printf("could not read message: %v\n", err)
		}

		log.Printf("message received: %s = %s\n", string(m.Key), string(m.Value))
	}
}

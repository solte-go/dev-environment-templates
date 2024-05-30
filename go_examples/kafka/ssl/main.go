package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log"
	"math/rand"
	"os"
	"time"
)

var (
	kafkaTopic    = "test-topic"
	kafkaBroker   = "localhost:29092"
	certificate   = "../../../configs/kafka/ssl/certs/client.crt"
	privateKey    = "../../../configs/kafka/ssl/certs/client.key"
	caCertificate = "../../../configs/kafka/ssl/certs/sandbox.crt"
)

func createTLSConfig(certificateFile, keyFile, caFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certificateFile, keyFile)
	if err != nil {
		return nil, err
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	tlsConfig := &tls.Config{
		RootCAs:            caCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}

	return tlsConfig, nil
}

func produce(kafkaBroker, kafkaTopic string, tlsConfig *tls.Config) error {
	//dialer := &kafka.Dialer{
	//	TLS: tlsConfig,
	//}

	writer := &kafka.Writer{
		Addr:     kafka.TCP(kafkaBroker),
		Topic:    kafkaTopic,
		Balancer: &kafka.LeastBytes{},
		Transport: &kafka.Transport{
			TLS:         tlsConfig,
			DialTimeout: 10 * time.Second,
		},
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	defer writer.Close()

	for i := 0; i < 100; i++ {
		err := writer.WriteMessages(context.Background(),
			kafka.Message{
				Key:   []byte("Key"),
				Value: []byte(fmt.Sprintf("Hello, kafka message %d with SSL; TS: %v", i, time.Now().Unix())),
			},
		)

		if err != nil {
			return err
		}
		randomSleep()
	}

	log.Printf("Producer: all messages have been send\n")

	return nil
}

func consume(kafkaBroker, kafkaTopic string, tlsConfig *tls.Config) {
	dialer := &kafka.Dialer{
		TLS: tlsConfig,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker},
		Topic:   kafkaTopic,
		GroupID: "test",
		//Partition:   1,
		StartOffset: kafka.LastOffset,
		MinBytes:    10e3,
		MaxBytes:    10e6,
		Dialer:      dialer,
	})

	defer reader.Close()

	for {
		message, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("received: %s\n", string(message.Value))

		err = reader.CommitMessages(context.Background(), message)
		if err != nil {
			log.Fatal(err)
			return
		}
	}
}

func main() {

	tlsConfig, err := createTLSConfig(certificate, privateKey, caCertificate)
	if err != nil {
		log.Fatalf("tls config error: %v\n", err)
	}

	go func() {
		err = produce(kafkaBroker, kafkaTopic, tlsConfig)
		if err != nil {
			log.Fatalf("Producer error: %v\n", err)
		}
	}()

	consume(kafkaBroker, kafkaTopic, tlsConfig)
}

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func randomSleep() {
	minimum := time.Millisecond * 200
	maximum := time.Second
	duration := minimum + time.Duration(r.Int63n(int64(maximum-minimum)))
	time.Sleep(duration)
}

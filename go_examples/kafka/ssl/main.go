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
	"strings"
	"sync"
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

	wg := sync.WaitGroup{}

	for round := 10; round > 0; round-- {
		wg.Add(1)

		go func() {
			var msgs = make([]kafka.Message, 0, 100)
			defer wg.Done()

			for i := 0; i < 100; i++ {
				message := kafka.Message{
					Key:   []byte("Key"),
					Value: []byte(fmt.Sprintf("Hello, kafka message %d with SSL; TS: %v", i, time.Now().Unix())),
				}
				msgs = append(msgs, message)
			}

			writer.WriteMessages(context.Background(), msgs...)
			randomSleep()

			log.Printf("Producer: batch of messages have been send\n")
		}()
	}

	wg.Wait()

	log.Printf("Producer: all messages have been send\n")
	return nil
}

func consume(kafkaBroker, kafkaTopic string, tlsConfig *tls.Config) {
	dialer := &kafka.Dialer{
		TLS: tlsConfig,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{kafkaBroker},
		Topic:       kafkaTopic,
		GroupID:     "test",
		StartOffset: kafka.FirstOffset,
		MinBytes:    1,
		MaxBytes:    10e6,
		Dialer:      dialer,
	})

	defer reader.Close()

	ticker := time.NewTicker(300 * time.Millisecond)

	var batch = make([]kafka.Message, 0, 50)
	var batchNum = 1

	massageCh := make(chan kafka.Message)

	go func() {
		for {
			message, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Fatalln(err)
			}
			massageCh <- message
		}
	}()

	for {
		select {
		case <-ticker.C:
			fmt.Println(strings.Repeat("~", 40))
			fmt.Printf("Counsumer: Requesting batch %d of messages\n", batchNum)

			// create a timer
			idleTimer := time.NewTimer(60 * time.Second)

			for len(batch) < 50 {
				select {
				case <-idleTimer.C:
					fmt.Println(strings.Repeat("=", 40))
					fmt.Println("No activity for 60 seconds. Exiting.")
					fmt.Printf("Counsumer: Processing the remaining messages from batch %d!\n", batchNum)
					return
				case msg := <-massageCh:
					// reset the timer when a new message is read
					if !idleTimer.Stop() {
						<-idleTimer.C
					}
					idleTimer.Reset(60 * time.Second)
					batch = append(batch, msg)
				}
			}
		}

		fmt.Println(strings.Repeat("=", 40))
		fmt.Printf("Counsumer: Processing batch number %d!\n", batchNum)
		batchNum++
		batch = batch[:0]
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
	minimum := time.Millisecond * 100
	maximum := time.Millisecond * 1000
	duration := minimum + time.Duration(r.Int63n(int64(maximum-minimum)))
	time.Sleep(duration)
}

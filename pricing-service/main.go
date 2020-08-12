package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"top-coins/pricing-service/service"
)

const (
	RABBIT_URL    = "amqp://guest:guest@localhost:5672"
	PRICING_QUEUE = "pricing_queue"
)

func main() {
	// Connect to RabbitMQ.
	conn, err := amqp.Dial(RABBIT_URL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	log.Println("Connected to RabbitMQ")

	// Set up the channel.
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(PRICING_QUEUE, false, false,
		false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue", err)
	}

	if err = ch.Qos(1, 0, false); err != nil {
		log.Fatalf("Failed to set prefetch settings: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// Set up the API.
	var c service.HttpClient
	c = http.DefaultClient
	api := service.NewAPI(c)

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			// Fetch data from the upstream API.
			bytes, err := api.FetchCryptocurrencies()
			if err != nil {
				log.Fatalf("Failed to fetch cryptocurrencies: %v", err)
			}

			// Process the bytes into []Cryptocurrency
			cryptos, err := api.ProcessCryptocurrencyBytes(bytes)
			if err != nil {
				log.Fatalf("Failed to get cryptocurrencies: %v", err)
			}

			log.Println(cryptos)
			body, err := json.Marshal(cryptos)
			if err != nil {
				log.Fatalf("Failed to marshal cryptocurrencies: %v", err)
			}

			if err = ch.Publish("", d.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: d.CorrelationId,
				Body:          body,
			}); err != nil {
				log.Fatalf("Failed to publish message: %v", err)
			}
			if err = d.Ack(false); err != nil {
				log.Printf("Failed to acknowledge message: %v", err)
			} else {
				log.Println("Acknowledged message")
			}
		}
	}()
	log.Println("Listening for RPC requests")
	<-forever
}

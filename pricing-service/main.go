package main

import (
	"github.com/streadway/amqp"
	"log"
	"top-coins/pricing-service/api"
)

const (
	RABBIT_URL = "amqp://guest:guest@localhost:5672"
)

func main() {
	// Connect to RabbitMQ.
	conn, err := amqp.Dial(RABBIT_URL)
	if err != nil {
		log.Fatalf("Could not establish AMQP connection: %v", err)
	}
	defer conn.Close()
	log.Println("Established AMQP connection")

	// Set up the channel.
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Could not open ch: %v", err)
	}
	defer ch.Close()
	q, err := ch.QueueDeclare("pricing_queue", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue", err)
	}

	if err = ch.Qos(1, 0, false); err != nil {
		log.Fatalf("Failed to set prefetch settings: %v", err)
	}

	_, err = ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	// Get the data from the upstream API.
	cryptos, err := api.GetCryptocurrencies()
	if err != nil {
		log.Fatalf("Failed to get cryptocurrencies: %v",err)
	}
	log.Println(cryptos)
}

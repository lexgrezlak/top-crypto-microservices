package main

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
)

const (
	RABBIT_URL = "amqp://guest:guest@localhost:5672"
)

// There are also other fields available. Check the docs for more information.
type Currency struct {
	Price float32 `json:"price"`
}

type Quote struct {
	USD Currency `json:"USD"`
	// We could also define fields for other currencies such as EUR
	// but USD is the one we're using - set in the `const` variable above.
}

type Cryptocurrency struct {
	Symbol string `json:"symbol"`
	Quote  Quote  `json:"quote"`
}

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

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	if err = ch.Publish("", "pricing_queue", false, false, amqp.Publishing{
		ContentType: "text/plain",
		ReplyTo:     q.Name,
	}); err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}
	for d := range msgs {
		var cryptos []Cryptocurrency
		_ = json.Unmarshal(d.Body, &cryptos)
		log.Println(cryptos)
	}
}

package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main()  {
	// Connect to RabbitMQ.
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatalf("Could not establish AMQP connection: %v", err)
	}
	defer conn.Close()
	log.Println("Established AMQP connection")

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("Could not open channel: %v", err)
	}
	log.Println(channel)
}
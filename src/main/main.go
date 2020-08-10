package main

import (
	"github.com/streadway/amqp"
	"log"
)

func main()  {
	conn, err := 	amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatalf("Could not establish AMQP connection: %v", err)
	}
	log.Println("Established AMQP connection")
	defer conn.Close()

}
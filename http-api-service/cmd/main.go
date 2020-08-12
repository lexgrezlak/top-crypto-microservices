package main

import (
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"time"
	"top-coins/http-api-service/internal/handler"
)

const (
	RABBIT_URL    = "amqp://guest:guest@localhost:5672"
)

func main() {
	// Connect to RabbitMQ.
	conn, err := amqp.Dial(RABBIT_URL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	log.Println("Established AMQP connection")

	// We're using mux instead of just http.HandleFunc because otherwise
	// there might be a security vulnerability. For more info check
	// https://stackoverflow.com/questions/36921190/difference-between-http-and-default-servemux/36921591
	mux := http.DefaultServeMux

	// Set up handlers.
	mux.Handle("/", handler.GetCryptocurrencies(conn))

	// Set up the server.
	server := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
	// Start the server.
	err = server.ListenAndServe()
	log.Fatalf("Failed to listen: %v", err)
}

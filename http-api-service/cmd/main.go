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

	// We're using mux instead of just http.HandleFunc because otherwise
	// there might be a security vulnerability. For more info check
	// https://stackoverflow.com/questions/36921190/difference-between-http-and-default-servemux/36921591
	mux := http.DefaultServeMux

	// Set up handlers.
	mux.Handle("/", handler.GetCryptocurrencies(conn))

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

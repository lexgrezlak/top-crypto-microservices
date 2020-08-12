package handler

import (
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"top-coins/http-api-service/internal/service"
)

const (
	PRICING_QUEUE = "pricing_queue"
)


// We're returning a handler to enable dependency injection.
func GetCryptocurrencies(conn *amqp.Connection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cryptos, err := service.HandleRPC(conn, PRICING_QUEUE)
		if err != nil {
			log.Printf("Failed to get messages: %v", err)
			return
		}
		w.Write(cryptos)
	}
}
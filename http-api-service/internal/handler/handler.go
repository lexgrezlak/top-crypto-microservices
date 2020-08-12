package handler

import (
	"github.com/streadway/amqp"
	"log"
	"net/http"
	"top-coins/http-api-service/internal/service"
)

const (
	PRICING_QUEUE = "pricing_queue"
	RANKING_QUEUE = "ranking_queue"
	PRICING_LIMIT = 100
	RANKING_LIMIT = PRICING_LIMIT
)

// We're returning a handler to enable dependency injection.
func GetCryptocurrencies(conn *amqp.Connection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pricing, err := service.HandleRPC(conn, PRICING_QUEUE, PRICING_LIMIT)
		if err != nil {
			log.Printf("Failed to handle pricing RPC request: %v", err)
		}
		ranking, err := service.HandleRPC(conn, RANKING_QUEUE, RANKING_LIMIT)
		if err != nil {
			log.Printf("Failed to handle ranking RPC request: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(append(pricing, ranking...))
	}
}
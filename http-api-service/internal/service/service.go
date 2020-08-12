package service

import (
	"errors"
	"github.com/streadway/amqp"
	"log"
	"math/rand"
	"strconv"
)

const (
	TEXT_PLAIN    = "text/plain"
)

// Used to generate correlation ID along with randInt.
// For more info check https://www.rabbitmq.com/tutorials/tutorial-six-go.html
func randString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

// Used to generate correlation ID along with randInt.
// For more info check https://www.rabbitmq.com/tutorials/tutorial-six-go.html
func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func HandleRPC(conn *amqp.Connection, queue string, limit int) ([]byte, error) {
	// Set up the channel.
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue", err)
	}

	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	// https://www.rabbitmq.com/tutorials/tutorial-six-go.html
	// Correlation id is used to differentiate between messages, so that
	// we can discard the ones we don't need.
	corrId := randString(32)

	if err = ch.Publish("", queue, false, false, amqp.Publishing{
		ContentType:   TEXT_PLAIN,
		CorrelationId: corrId,
		Body: []byte(strconv.Itoa(limit)),
		ReplyTo:       q.Name,
	}); err != nil {
		log.Fatalf("Failed to publish message: %v", err)
	}

	for d := range msgs {
		if d.CorrelationId == corrId {
			return d.Body, nil
		}
	}
	return nil, errors.New("delivery not received")
}


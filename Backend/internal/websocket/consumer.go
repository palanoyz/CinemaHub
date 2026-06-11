package websocket

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type BookingMessage struct {
	UserID     string   `json:"user_id"`
	ShowtimeID string   `json:"showtime_id"`
	SeatIDs    []string `json:"seat_ids"`
}

func StartBookingConsumer(conn *amqp.Connection) {
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	q, err := ch.QueueDeclare(
		"bookings", // name
		true,       // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("RabbitMQ Booking Consumer started")

	go func() {
		for d := range msgs {
			var msg BookingMessage
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				log.Printf("Error decoding booking message: %v", err)
				continue
			}

			// SIMULATE BACKGROUND WORK (e.g. Email / PDF Generation)
			log.Printf(" [🐰] ASYNC WORKER: Processing booking for user %s", msg.UserID)
			log.Printf(" [🐰] ASYNC WORKER: Generating tickets for seats: %v", msg.SeatIDs)
			log.Printf(" [🐰] ASYNC WORKER: Sending confirmation email...")
			log.Printf(" [🐰] ASYNC WORKER: Done!")
		}
	}()
}

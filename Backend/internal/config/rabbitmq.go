package config

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RabbitMQConn *amqp.Connection

func ConnectRabbitMQ() *amqp.Connection {
	url := AppConfig.RabbitMQ_URL

	conn, err := amqp.Dial(url)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ: ", err)
	}

	log.Println("Connected to RabbitMQ successfully")
	RabbitMQConn = conn
	return conn
}

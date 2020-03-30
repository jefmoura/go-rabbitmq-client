package settings

import (
	"os"
)

type configuration struct {
	RabbitMQUser      string
	RabbitMQPassword  string
	RabbitMQHost      string
	RabbitMQPort      string
	RabbitMQVHost     string
	RabbitMQQueue     string
}

// Config composes the set of values necessary
// for configuring the application
var Config = configuration{
	RabbitMQUser:      os.Getenv("RABBITMQ_USER"),
	RabbitMQPassword:  os.Getenv("RABBITMQ_PASSWORD"),
	RabbitMQHost:      os.Getenv("RABBITMQ_HOST"),
	RabbitMQPort:      os.Getenv("RABBITMQ_PORT"),
	RabbitMQVHost:     os.Getenv("RABBITMQ_VHOST"),
	RabbitMQQueue:     os.Getenv("RABBITMQ_QUEUE"),
}

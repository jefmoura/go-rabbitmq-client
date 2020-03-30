package main

import (
	"go-rabbitmq-client/rabbitmq"
)

func main() {
	cli := rabbitmq.NewClient()
	cli.Consume()
}

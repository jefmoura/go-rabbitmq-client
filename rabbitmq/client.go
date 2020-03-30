package rabbitmq

import (
	"bytes"
	"fmt"
	"log"
	"time"

	"go-rabbitmq-client/settings"

	"github.com/streadway/amqp"
)

// Client composes the set of values necessary for connecting
// to a Borker, and consuming or producing messages
type Client struct {
	username  string
	password  string
	hostname  string
	port      string
	vhost     string
	queueName string
}

func (c Client) handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func (c Client) connect() *amqp.Connection {
	AMQPConnectionURL := fmt.Sprintf("amqp://%s:%s@%s:%s/%s", c.username, c.password, c.hostname, c.port, c.vhost)
	conn, err := amqp.Dial(AMQPConnectionURL)
	c.handleError(err, "Failed to connect to RabbitMQ")

	return conn
}

func (c Client) openChannel(conn *amqp.Connection) *amqp.Channel {
	channel, err := conn.Channel()
	c.handleError(err, "Failed to open a channel")

	return channel
}

func (c Client) createQueue(channel *amqp.Channel) amqp.Queue {
	queue, err := channel.QueueDeclare(
		c.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	c.handleError(err, "Failed to declare a queue")

	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	c.handleError(err, "Failed to set QoS")

	return queue
}

func (c Client) subscribe(queueName string, channel *amqp.Channel) <-chan amqp.Delivery {
	messages, err := channel.Consume(
		queueName, // queue
		"",        // consumer
		false,     // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	c.handleError(err, "Failed to register a consumer")

	return messages
}

// Consume retrieves messages from the Broker from the specified topic
func (c Client) Consume() {
	conn := c.connect()
	channel := c.openChannel(conn)
	queue := c.createQueue(channel)
	msgs := c.subscribe(queue.Name, channel)

	defer conn.Close()
	defer channel.Close()

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			dotCount := bytes.Count(d.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			log.Printf("Done")
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// NewClient returns Client with values defined in the environment
func NewClient() Client {
	return Client{
		username:  settings.Config.RabbitMQUser,
		password:  settings.Config.RabbitMQPassword,
		hostname:  settings.Config.RabbitMQHost,
		port:      settings.Config.RabbitMQPort,
		vhost:     settings.Config.RabbitMQVHost,
		queueName: settings.Config.RabbitMQQueue,
	}
}

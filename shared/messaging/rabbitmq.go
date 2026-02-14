package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"go-ride/shared/contracts"
	"go-ride/shared/env"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	TripExchange       = "trip"
	DeadLetterExchange = "dlx"
)

var (
	AMQPAddr = env.GetString("RABBITMQ_URI", "amqp://guest:guest@rabbitmq:5672")
)

type RabbitMQ struct {
	conn    *amqp.Connection
	Channel *amqp.Channel
}

func NewRabbitMQ(uri string) (*RabbitMQ, error) {
	conn, err := amqp.Dial(AMQPAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to rabbitmq: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create channel: %v", err)
	}

	rmq := &RabbitMQ{
		conn:    conn,
		Channel: ch,
	}

	if err := rmq.setupExchangesAndQuests(); err != nil {
		rmq.Close()
		return nil, fmt.Errorf("failed to setup exchanges and queues: %v", err)
	}

	return rmq, nil
}

func (r *RabbitMQ) PublishMessage(ctx context.Context, routingKey string, message contracts.AmqpMessage) error {
	log.Printf("publishing message with routing key: %s", routingKey)

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	msg := amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "application/json",
		Body:         jsonMsg,
	}

	return r.Channel.PublishWithContext(ctx,
		TripExchange,
		routingKey,
		false,
		false,
		msg,
	)
}

func (r *RabbitMQ) setupDeadLetterExchange() error {
	// Declare the dead letter exchange
	err := r.Channel.ExchangeDeclare(
		DeadLetterExchange,
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare dead letter exchange: %v", err)
	}

	q, err := r.Channel.QueueDeclare(
		DeadLetterQueue,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare dead letter queue: %v", err)
	}

	// Bind the queue to the exchange with a wildcard routing key
	err = r.Channel.QueueBind(
		q.Name,
		"#", // wildcard routing key to catch all messages
		DeadLetterExchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to bind dead letter queue: %v", err)
	}

	return nil
}

func (r *RabbitMQ) setupExchangesAndQuests() error {
	// First setup the DLQ exchange and queue
	if err := r.setupDeadLetterExchange(); err != nil {
		return err
	}

	err := r.Channel.ExchangeDeclare(
		TripExchange, // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare exchange: %s: %v", TripExchange, err)
	}

	if err := r.declareAndBindQueue(
		FindAvailableDriversQueue,
		[]string{
			contracts.TripEventCreated,
			contracts.TripEventDriverNotInterested,
		},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		NotifyNoDriversFoundQueue,
		[]string{contracts.TripEventNoDriversFound},
		TripExchange,
	); err != nil {
		return err
	}

	if err := r.declareAndBindQueue(
		NotifyDriverAssignQueue,
		[]string{contracts.TripEventDriverAssigned},
		TripExchange,
	); err != nil {
		return err
	}

	return nil
}

func (r *RabbitMQ) declareAndBindQueue(queueName string, messageTypes []string, exchange string) error {
	// Add dead letter configuration
	args := amqp.Table{
		"x-dead-letter-exchange": DeadLetterExchange,
	}

	queue, err := r.Channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		args,      // arguments with DLX config
	)
	if err != nil {
		log.Fatal(err)
	}

	for _, msgType := range messageTypes {
		err := r.Channel.QueueBind(
			queue.Name, // queue name
			msgType,    // routing key
			exchange,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue to %s: %v", queueName, err)
		}
	}

	return nil
}

func (r *RabbitMQ) Close() {
	if r.conn != nil {
		r.conn.Close()
	}
	if r.Channel != nil {
		r.Channel.Close()
	}
}

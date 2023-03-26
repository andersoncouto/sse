package rabbitmq

import amqp "github.com/rabbitmq/amqp091-go"

func OpenChannel() (channel *amqp.Channel, err error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}

func Consume(queue string, ch *amqp.Channel, out chan amqp.Delivery) error {
	msgs, err := ch.Consume(
		queue,         // queue
		"go-consumer", // consumer
		false,         // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)
	if err != nil {
		return err
	}
	for msg := range msgs {
		out <- msg
	}
	return nil
}

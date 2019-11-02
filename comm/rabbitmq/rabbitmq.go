package rabbitmq

import (
	"encoding/json"
	"github.com/streadway/amqp"
)

type RabbitMq struct {
	channel *amqp.Channel
	conn *amqp.Connection
	Name string
	exchange string
}

func New(serverAddr string) *RabbitMq {
	conn, err := amqp.Dial(serverAddr)
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		true,
		false,
		false,
		nil)

	if err != nil {
		panic(err)
	}

	return &RabbitMq{channel : ch, Name: q.Name, conn : conn}
}

func (q *RabbitMq) Bind(exchange string) {
	err := q.channel.QueueBind(
		q.Name,
		"",
		exchange,
		false,
		nil)
	if err != nil {
		panic(err)
	}

	q.exchange = exchange
}

func (q *RabbitMq) Send(queue string, body interface{}) {
	str, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	err = q.channel.Publish(
		"",
		queue,
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body : []byte(str),
		})
	if err != nil {
		panic(err)
	}
}

func (q *RabbitMq) Publish(exchange string, body interface{}) {
	str, err := json.Marshal(body)
	if err != nil {
		panic(err)
	}

	err = q.channel.Publish(exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ReplyTo: q.Name,
			Body : []byte(str),
		})

	if err != nil {
		panic(err)
	}
}

func (q *RabbitMq) Cosume() <-chan amqp.Delivery {
	c, err := q.channel.Consume(q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		panic(err)
	}

	return c
}

func (q *RabbitMq) Close() {
	q.channel.Close()
	q.conn.Close()
}

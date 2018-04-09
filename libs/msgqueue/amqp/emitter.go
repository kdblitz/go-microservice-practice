package amqp

import (
	"github.com/streadway/amqp"
	"github.com/kdblitz/go-microservice-practice/libs/msgqueue"
	"encoding/json"
)

type amqpEventEmitter struct {
	connection amqp.Connection
}

func NewAMQPEventEmitter(conn *amqp.Connection) (msgqueue.EventEmitter, error) {
	emitter := &amqpEventEmitter{
		connection: conn,
	}
	err := emitter.setup()
	if err != nil {
		return nil, err
	}
	return &emitter, nil
}

func (a *amqpEventEmitter) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	return channel.ExchangeDeclare("events", "topic", true, false, false, false, nil)
}

func (a *amqpEventEmitter) emit(event msgqueue.Event) error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	jsonDoc, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := amqp.Publishing{
		Headers: amqp.Table{"x-event-name": event.EventName()},
		ContentType: "application/json",
		Body: jsonDoc,
	}

	return channel.Publish("events", event.EventName(), false, false, msg)
}
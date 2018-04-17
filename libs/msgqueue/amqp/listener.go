package amqp

import (
	"github.com/streadway/amqp"
	"github.com/kdblitz/go-microservice-practice/libs/msgqueue"
	"fmt"
	"github.com/kdblitz/go-microservice-practice/contracts"
	"encoding/json"
)

type amqpEventListener struct {
	connection *amqp.Connection
	queue string
}

func NewAMQPEventListener(conn *amqp.Connection, queue string) (msgqueue.EventListener, error) {
	listener := &amqpEventListener{
		connection: conn,
		queue: queue,
	}

	err := listener.setup()
	if err != nil {
		return nil, err
	}
	return listener, nil
}

func (a *amqpEventListener) setup() error {
	channel, err := a.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	_, err = channel.QueueDeclare(a.queue, true, false, false, false, nil)
	return err
}

func (a *amqpEventListener) Listen(eventNames ...string) (<-chan msgqueue.Event, <-chan error, error) {
	channel, err := a.connection.Channel()
	if err != nil {
		return nil, nil, err
	}
	defer channel.Close()

	for _, eventName := range eventNames {
		if err := channel.QueueBind(a.queue, eventName, "events", false, nil); err != nil {
			return nil, nil, err
		}
	}
	fmt.Println("listening")
	msgs, err := channel.Consume(a.queue, "", false, false, false, false, nil)
	if err != nil {
		return nil, nil, err
	}

	events := make(chan msgqueue.Event)
	errors := make(chan error)

	go func() {
		for msg := range msgs {
			rawEventName, ok := msg.Headers["x-event-name"]
			if !ok {
				errors <- fmt.Errorf("no x-event-name header on msg")
				msg.Nack(false, false)
				continue
			}
			eventName, ok := rawEventName.(string)
			if !ok {
				errors <- fmt.Errorf("x-event-name if not a string, but %t", rawEventName)
				msg.Nack(false, false)
				continue
			}
			var event msgqueue.Event
			switch eventName {
			case "eventCreated":
				event = new(contracts.EventCreatedEvent)
			}
			err := json.Unmarshal(msg.Body, event)
			if err != nil {
				errors <- fmt.Errorf("unmarshall error: %s %s", eventName, err)
				msg.Nack(false, false)
				continue
			}
			events <- event
			msg.Ack(false)
		}
	}()

	return events, errors, nil
}

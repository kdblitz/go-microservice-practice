package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/kdblitz/go-microservice-practice/libs/msgqueue"
	"log"
	"encoding/json"
	"fmt"
	"github.com/kdblitz/go-microservice-practice/contracts"
	"github.com/mitchellh/mapstructure"
)

type kafkaEventListener struct {
	consumer sarama.Consumer
	partitions []int32
}

func NewKafkaEventListener(client sarama.Client, partitions []int32) (msgqueue.EventListener, error) {
	consumer, err := sarama.NewConsumerFromClient(client);
	if err != nil {
		return nil, err
	}

	listener := &kafkaEventListener{
		consumer: consumer,
		partitions: partitions,
	}
	return listener, nil
}

func (ke *kafkaEventListener) Listen(events ...string) (<-chan msgqueue.Event, <-chan error, error){
	var err error

	topic := "events"
	results := make(chan msgqueue.Event)
	errors := make(chan error)

	partitions := ke.partitions
	if len(partitions) == 0 {
		partitions, err = ke.consumer.Partitions(topic)
		if err != nil {
			return nil, nil, err
		}
	}
	log.Printf("topic %s has partitions: %v", topic, partitions)

	for _, partitions := range partitions {
		con, err := ke.consumer.ConsumePartition(topic, partitions, 0)
		if err != nil {
			return nil, nil, err
		}

		go func() {
			for msg := range con.Messages() {
				body := messageEnvelope{}
				err := json.Unmarshal(msg.Value, &body)
				if err != nil {
					errors <- fmt.Errorf("could not decode message: %s", err)
					continue
				}

				var event msgqueue.Event
				switch body.EventName {
				case "eventCreated":
					event = &contracts.EventCreatedEvent{}
				case "locationCreated":
					event = &contracts.LocationCreatedEvent{}
				default:
					errors <- fmt.Errorf("unknown event type: %s", body.EventName)
					continue
				}

				cfg := mapstructure.DecoderConfig{
					Result: event,
					TagName: "json",
				}

				dec, err := mapstructure.NewDecoder(&cfg)
				if err != nil {
					errors <- fmt.Errorf("couldn't initialize decoder %s %s", body.EventName, err)
				}
				err = dec.Decode(body.Payload)
				if err != nil {
					errors <- fmt.Errorf("unmarshall error: %s %s", body.EventName, err)
				}
				results <- event
			}
		}()
	}
	return results, errors, nil
}

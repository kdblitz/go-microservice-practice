package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/kdblitz/go-microservice-practice/libs/msgqueue"
	"encoding/json"
)

type kafkaEventEmitter struct {
	producer sarama.SyncProducer
}

func NewKafkaEventEmitter(client sarama.Client) (msgqueue.EventEmitter, error) {
	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return nil, err
	}

	emitter := &kafkaEventEmitter{
		producer: producer,
	}

	return emitter, nil
}

func (ke *kafkaEventEmitter) Emit(event msgqueue.Event) error {
	envelope := messageEnvelope{event.EventName(), event}
	jsonBody, err := json.Marshal(&envelope)
	if err != nil {
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic: event.PartitionKey(),
		Value: sarama.ByteEncoder(jsonBody),
	}

	_, _, err = ke.producer.SendMessage(msg)
	return err
}

package main

import (
	"github.com/Shopify/sarama"
	"fmt"
)

func main() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	brokers := []string{"192.168.56.101:9092"}
	client, err := sarama.NewClient(brokers, config)

	if err != nil {
		panic(err)
	}
	producer, err := sarama.NewSyncProducerFromClient(client)

	if err != nil {
		panic(err)
	}
	fmt.Println(producer)

}
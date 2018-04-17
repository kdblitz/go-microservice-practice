package configuration

import (
	"github.com/kdblitz/go-microservice-practice/libs/persistence/dblayer"
	"os"
	"fmt"
	"encoding/json"
)

var (
	DBTypeDefault = dblayer.DBTYPE("mongodb")
	DBConnectionDefault = "mongodb://127.0.0.1"
	RestfulEPDefault = "localhost:8181"
	RestfulTLSEPDefault = "localhost:9191"
	MessageBrokerTypeDefault = "kafka"
	AMQPMessageBrokerDefault = "amqp://guest:guest@192.168.56.101:5672"
	KafkaMessageBrokersDefault = []string{"192.168.56.101:9092"}
)

type ServiceConfig struct {
	Databasetype dblayer.DBTYPE `json:"databasetype"`
	DBConnection string `json:"dbconnection"`
	RestfulEndpoint string `json:"restfulapi_endpoint"`
	RestfulTLSEndPoint string `json:"restfulapi_tlsendpoint"`
	MessageBrokerType string `json:"message_broker_type"`
	AMQPMessageBroker string `json:"amqp_message_broker"`
	KafkaMessageBrokers []string `json:"kafka_message_brokers"`
}

func ExtractConfig(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
		RestfulTLSEPDefault,
		MessageBrokerTypeDefault,
		AMQPMessageBrokerDefault,
		KafkaMessageBrokersDefault,
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Config file not found, using defaults")
		return conf, err
	}
	err = json.NewDecoder(file).Decode(&conf)
	if broker := os.Getenv("AMQP_URL"); broker != "" {
		conf.AMQPMessageBroker = broker
	}
	return conf, err
}
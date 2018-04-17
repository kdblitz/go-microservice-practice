package main

import (
	"flag"
	"github.com/kdblitz/go-microservice-practice/libs/configuration"
	"github.com/kdblitz/go-microservice-practice/libs/persistence/dblayer"
	"github.com/streadway/amqp"
	msgqueue_amqp "github.com/kdblitz/go-microservice-practice/libs/msgqueue/amqp"
	"github.com/kdblitz/go-microservice-practice/bookingservice/listener"
	"github.com/kdblitz/go-microservice-practice/libs/msgqueue"
	"github.com/Shopify/sarama"
	"github.com/kdblitz/go-microservice-practice/libs/msgqueue/kafka"
	"github.com/kdblitz/go-microservice-practice/eventsservice/rest"
)

func main()  {
	confPath := flag.String("config", `.\configuration\config.json`, "config file path")
	flag.Parse()
	config, err := configuration.ExtractConfig(*confPath)
	//if err != nil {
	//	panic("Extract config error: " + err.Error())
	//}

	dblayer, err := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
	if err != nil {
		panic(err)
	}

	var eventListener msgqueue.EventListener
	switch config.MessageBrokerType {
	case "amqp":
		mq, err := amqp.Dial(config.AMQPMessageBroker)
		if err != nil {
			panic(err)
		}

		eventListener, err = msgqueue_amqp.NewAMQPEventListener(mq, "events")
		if err != nil {
			panic(err)
		}
	case "kafka":
		conf := sarama.NewConfig()
		conf.Producer.Return.Successes = true
		conn, err := sarama.NewClient(config.KafkaMessageBrokers, conf)
		if err != nil {
			panic(err)
		}

		eventListener, err = kafka.NewKafkaEventListener(conn, []int32{})
	default:
		panic("invalid message broker type:" + config.MessageBrokerType)
	}

	processor := &listener.EventProcessor{
		EventListener: eventListener,
		Database: dblayer,
	}

	processor.ProcessEvents()

	rest.ServeAPI(config.RestfulEndpoint, )
}

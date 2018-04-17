package main

import (
  "flag"
  "github.com/kdblitz/go-microservice-practice/libs/configuration"
  "fmt"
  "github.com/kdblitz/go-microservice-practice/libs/persistence/dblayer"
  "github.com/kdblitz/go-microservice-practice/eventsservice/rest"
  "log"
  "github.com/streadway/amqp"
  msg_queueamqp "github.com/kdblitz/go-microservice-practice/libs/msgqueue/amqp"
  "github.com/kdblitz/go-microservice-practice/libs/msgqueue"
  "github.com/Shopify/sarama"
  "github.com/kdblitz/go-microservice-practice/libs/msgqueue/kafka"
)

func main() {
  confPath := flag.String("conf", `.\configuration\config.json`, "configuration json file")
  flag.Parse()

  config, _ := configuration.ExtractConfig(*confPath)

  var emitter msgqueue.EventEmitter
  switch config.MessageBrokerType {
  case "amqp":
    conn, err := amqp.Dial(config.AMQPMessageBroker)
    if err != nil {
      panic(err)
    }

    emitter, err = msg_queueamqp.NewAMQPEventEmitter(conn)
    if err != nil {
      panic(err)
    }
  case "kafka":
    conf := sarama.NewConfig()
    conf.Producer.Return.Successes = true
    fmt.Println("kafka")
    fmt.Println(config.KafkaMessageBrokers)
    conn, err := sarama.NewClient(config.KafkaMessageBrokers, conf)
    if err != nil {
      panic(err)
    }
    emitter, err = kafka.NewKafkaEventEmitter(conn)
    if err != nil {
      panic(err)
    }
    fmt.Println("success")
  default:
    panic("invalid message broker type:" + config.MessageBrokerType)
  }

  fmt.Println("connecting to db")
  dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
  httpErrChan, httpsErrChan := rest.ServeAPI(config.RestfulEndpoint, config.RestfulTLSEndPoint, dbhandler, emitter)
  select {
  case err := <-httpErrChan:
    log.Fatal("Http error:", err)
  case err := <-httpsErrChan:
    log.Fatal("Https error:", err)
  }
}
package main

import (
  "flag"
  "github.com/kdblitz/go-microservice-practice/libs/configuration"
  "fmt"
  "github.com/kdblitz/go-microservice-practice/libs/persistence/dblayer"
  "github.com/kdblitz/go-microservice-practice/eventsservice/rest"
)

func main() {
  confPath := flag.String("conf", `.\configuration\config.json`, "configuration json file")
  flag.Parse()

  config, _ := configuration.ExtractConfig(*confPath)

  fmt.Println("connecting to db")
  dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
  rest.ServeAPI(config.RestfulEndpoint, dbhandler)
}
package main

import (
  "flag"
  "github.com/kdblitz/go-microservice-practice/libs/configuration"
  "fmt"
  "github.com/kdblitz/go-microservice-practice/libs/persistence/dblayer"
)

func main() {
  confPath := flag.String("conf", `.\configuration\config.json`, "configuration json file")
  flag.Parse()

  config, _ := configuration.ExtractConfig(*confPath)

  fmt.Println("connecting to db")
  dbhandler, _ := dblayer.NewPersistenceLayer(config.Databasetype, config.DBConnection)
  ServeAPI(config.RestfulEndpoint, dbhandler)
}
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
)

type ServiceConfig struct {
	Databasetype dblayer.DBTYPE `json:"databasetype"`
	DBConnection string `json:"dbconnection"`
	RestfulEndpoint string `json:"restfulapi_endpoint"`
}

func ExtractConfig(filename string) (ServiceConfig, error) {
	conf := ServiceConfig{
		DBTypeDefault,
		DBConnectionDefault,
		RestfulEPDefault,
	}
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Config file not found, using defaults")
		return conf, err
	}
	err = json.NewDecoder(file).Decode(&conf)
	return conf, err
}
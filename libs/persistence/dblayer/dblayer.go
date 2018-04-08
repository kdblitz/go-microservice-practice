package dblayer

import (
	"github.com/kdblitz/go-microservice-practice/libs/persistence"
	"github.com/kdblitz/go-microservice-practice/libs/persistence/mongo"
)

type DBTYPE string

const (
	MONGODB DBTYPE = "mongodb"
)

func NewPersistenceLayer(options DBTYPE, connection string) (persistence.DatabaseHandler, error) {
	switch options {
	case MONGODB:
		return mongo.NewMongoDBLayer(connection)
	}
	return nil, nil
}
package persistence

import "gopkg.in/mgo.v2/bson"

type Event struct {
	ID bson.ObjectId `bson:"_id"`
	Name string
	Location Location
}

type Location struct {
	ID bson.ObjectId `bson:"_id"`
	Name string
}

type Hall struct {
	Name string
	Capacity int
}

type Booking struct {
	Date int64
	EventID []byte
	Seats int
}
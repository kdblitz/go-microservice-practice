package listener

import (
	"github.com/kdblitz/go-microservice-practice/libs/msgqueue"
	"github.com/kdblitz/go-microservice-practice/libs/persistence"
	"fmt"
	"github.com/kdblitz/go-microservice-practice/contracts"
	"gopkg.in/mgo.v2/bson"
)

type EventProcessor struct {
	EventListener msgqueue.EventListener
	Database persistence.DatabaseHandler
}

func (ep *EventProcessor) ProcessEvents() error {
	received, errors, err := ep.EventListener.Listen("eventCreated")
	if err != nil {
		return err
	}
	fmt.Println("start process")
	for {
		select {
		case evt := <-received:
			ep.handleEvent(evt)
		case err := <-errors:
			fmt.Printf("received error while processing: %s", err.Error())
		}
	}
}

func (ep *EventProcessor) handleEvent(event msgqueue.Event) {
	fmt.Println("handling evt %T", event)
	switch e := event.(type) {
	case *contracts.EventCreatedEvent:
		fmt.Println("adding event")
		ep.Database.AddEvent(persistence.Event{
			ID: bson.ObjectIdHex(e.ID),
			Name: e.Name,
		})
	case *contracts.LocationCreatedEvent:
		fmt.Printf("todo: handle location")
		//ep.Database.AddLocation(persistence.Location{ID: bson.ObjectId(e.)})
	default:
		fmt.Printf("unknown event: %T", e)
	}
}
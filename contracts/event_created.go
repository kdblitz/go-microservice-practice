package contracts

//import "time"

type EventCreatedEvent struct {
	ID string `json:"id"`
	Name string `json:"name"`
	LocationID string `json:"location_id"`
	//Start time.Time `json:"start_time"`
	//End time.Time `json:"end_time"`
}

func (ec *EventCreatedEvent) EventName() string {
	return "eventCreated"
}

func (ec *EventCreatedEvent) PartitionKey() string {
	return ec.ID
}
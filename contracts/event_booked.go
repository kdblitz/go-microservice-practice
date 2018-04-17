package contracts

type EventBookedEvent struct {
	EventID string `json:"eventId"`
	UserID string `json:"userId"`
}

func (e *EventBookedEvent) EventName() string {
	return "eventBooked"
}

func (e *EventBookedEvent) PartitionKey() string {
	return e.EventID + "_" + e.UserID
}
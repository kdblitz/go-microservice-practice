package contracts

type LocationCreatedEvent struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

func (lc *LocationCreatedEvent) EventName() string {
	return "locationCreated"
}

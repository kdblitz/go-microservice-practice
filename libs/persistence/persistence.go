package persistence

type DatabaseHandler interface {
	AddEvent(Event) ([]byte, error)
	FindEvent([]byte) (Event, error)
	FindEventByName(string) (Event, error)
	FindAllEvents() ([]Event, error)
}
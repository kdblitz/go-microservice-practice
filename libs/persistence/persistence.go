package persistence

type DatabaseHandler interface {
	AddEvent(Event) ([]byte, error)
	AddBookingForUser([]byte, Booking) error
	FindEvent([]byte) (Event, error)
	FindEventByName(string) (Event, error)
	FindAllEvents() ([]Event, error)
}
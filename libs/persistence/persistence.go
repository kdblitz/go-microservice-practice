package persistence

type DatabaseHandler interface {
	AddEvent() ([]byte, error)
	FindEvent() (Event, error)
	FindAllEvent() ([]Event, error)
}
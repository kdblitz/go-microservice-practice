package rest

import (
	"github.com/kdblitz/go-microservice-practice/libs/persistence"
	"github.com/gorilla/mux"
	"net/http"
)

func ServeAPI(endpoint string, databasehandler persistence.DatabaseHandler) error {
	handler := NewEventHandler(databasehandler)

	mainRouter := mux.NewRouter()
	eventsRouter := mainRouter.PathPrefix("/events").Subrouter()

	eventsRouter.Methods("GET").Path("{SearchCriteria}/{query}").HandlerFunc(handler.FindEvent)
	eventsRouter.Methods("GET").Path("").HandlerFunc(handler.AllEvent)
	eventsRouter.Methods("POST").Path("").HandlerFunc(handler.NewEvent)
	return http.ListenAndServe(endpoint, mainRouter)
}

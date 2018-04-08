package main

import (
	"net/http"
	"github.com/kdblitz/go-microservice-practice/libs/persistence"
	"github.com/gorilla/mux"
)

func ServeAPI(endpoint string, databasehandler persistence.DatabaseHandler) {
	handler := NewEventHandler(databasehandler)

	mainRouter := mux.NewRouter()
	eventsRouter := mainRouter.PathPrefix("/events").Subrouter()

	eventsRouter.Methods("GET").Path("{SearchCriteria}/{query}").HandlerFunc(handler.FindEvent)
	eventsRouter.Methods("GET").Path("").HandlerFunc(handler.AllEvent)
	eventsRouter.Methods("POST").Path("").HandlerFunc(handler.NewEvent)
	return http.ListenAndServe(endpoint, mainRouter)
}
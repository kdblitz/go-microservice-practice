package rest

import (
	"github.com/kdblitz/go-microservice-practice/libs/persistence"
	"github.com/gorilla/mux"
	"net/http"
	"github.com/kdblitz/go-microservice-practice/libs/msgqueue"
	"github.com/gorilla/handlers"
)

func ServeAPI(endpoint, tlsEndpoint string, databaseHandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) (chan error, chan error) {
	handler := NewEventHandler(databaseHandler, eventEmitter)

	mainRouter := mux.NewRouter()
	eventsRouter := mainRouter.PathPrefix("/events").Subrouter()

	server := handlers.CORS()(mainRouter)

	eventsRouter.Methods("GET").Path("/{SearchCriteria}/{query}").HandlerFunc(handler.FindEvent)
	eventsRouter.Methods("GET").Path("").HandlerFunc(handler.AllEvent)
	eventsRouter.Methods("GET").Path("/{eventID}").HandlerFunc(handler.OneEventHandler)
	eventsRouter.Methods("POST").Path("").HandlerFunc(handler.NewEvent)

	httpErrChan := make(chan error)
	httpsErrChan := make(chan error)
	go func() {httpErrChan <- http.ListenAndServe(endpoint, server)}()
	go func() {httpsErrChan <- http.ListenAndServeTLS(tlsEndpoint, "cert.pem", "key.pem", server)}()
	return httpErrChan, httpsErrChan
}

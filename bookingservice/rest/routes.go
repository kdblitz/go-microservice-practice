package rest

import (
	"github.com/gorilla/mux"
	"github.com/kdblitz/go-microservice-practice/libs/persistence"
	"github.com/kdblitz/go-microservice-practice/libs/msgqueue"
	"net/http"
	"github.com/gorilla/handlers"
	"time"
)

func ServeAPI(addr string, database persistence.DatabaseHandler, emitter msgqueue.EventEmitter) {
	r := mux.NewRouter()
	r.Methods("post").Path("/events/{eventID}/bookings").Handler(&CreateBookingHandler{emitter, database})

	srv := http.Server{
		Handler: handlers.CORS()(r),
		Addr: addr,
		WriteTimeout: 2 * time.Second,
		ReadTimeout: 1 * time.Second,
	}

	srv.ListenAndServe()
}
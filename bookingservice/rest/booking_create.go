package rest

import (
	"github.com/kdblitz/go-microservice-practice/libs/msgqueue"
	"github.com/kdblitz/go-microservice-practice/libs/persistence"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"encoding/hex"
	"encoding/json"
	"time"
	"github.com/kdblitz/go-microservice-practice/contracts"
)

type createBookingRequest struct {
	Seats int `json:"seats"`
}

type CreateBookingHandler struct {
	emitter msgqueue.EventEmitter
	database persistence.DatabaseHandler
}

func (h *CreateBookingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	eventID, ok := vars["eventID"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `missing "eventID" parameter`)
		return
	}

	eventIDMongo, _ := hex.DecodeString(eventID)
	event, err := h.database.FindEvent(eventIDMongo)
	if err != nil {
		w.WriteHeader(404)
		fmt.Fprintf(w, "event %s could not be loaded: %s", eventID, err)
		return
	}

	bookingRequest := createBookingRequest{}
	err = json.NewDecoder(r.Body).Decode(&bookingRequest)
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprintf(w, "could not decode json: %s", err)
		return
	}

	if bookingRequest.Seats <= 0 {
		w.WriteHeader(400);
		fmt.Fprintf(w, "seats must be more positive: (got %d)", bookingRequest.Seats)
		return
	}

	eventIDAsBytes, _ := event.ID.MarshalText()
	booking := persistence.Booking{
		Date: time.Now().Unix(),
		EventID: eventIDAsBytes,
		Seats: bookingRequest.Seats,
	}

	userID := "myUserID"
	msg := contracts.EventBookedEvent{
		EventID: event.ID.Hex(),
		UserID: userID,
	}

	h.emitter.Emit(&msg)

	h.database.AddBookingForUser([]byte(userID), booking)

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(201)
	json.NewEncoder(w).Encode(&booking)
}

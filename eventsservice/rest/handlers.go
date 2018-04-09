package rest

import (
	"github.com/kdblitz/go-microservice-practice/libs/persistence"
	"net/http"
	"github.com/gorilla/mux"
	"fmt"
	"strings"
	"encoding/hex"
	"encoding/json"
	"github.com/kdblitz/go-microservice-practice/libs/msgqueue"
	"github.com/kdblitz/go-microservice-practice/contracts"
)

type eventServiceHandler struct {
	dbhandler persistence.DatabaseHandler
	eventEmitter msgqueue.EventEmitter
}

func NewEventHandler(databasehandler persistence.DatabaseHandler, eventEmitter msgqueue.EventEmitter) *eventServiceHandler {
	return &eventServiceHandler{
		dbhandler: databasehandler,
		eventEmitter: eventEmitter,
	}
}

func (eh *eventServiceHandler) FindEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r);
	criteria, ok := vars["SearchCriteria"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{error: No search criteria found you can search by id using /id/[id], or name using /name/[name]}`)
		return
	}
	query, ok := vars["query"]
	if !ok {
		w.WriteHeader(400)
		fmt.Fprint(w, `{error: No search criteria found you can search by id using /id/[id], or name using /name/[name]}`)
		return
	}
	var event persistence.Event
	var err error
	switch strings.ToLower(criteria) {
		case "name":
		event, err = eh.dbhandler.FindEventByName(query)
		case "id":
		id, err := hex.DecodeString(query)
		if err == nil {
			event, err = eh.dbhandler.FindEvent(id)
		}
	}
	if err != nil {
		fmt.Fprintf(w, "{error %s}", err)
		return
	}
	w.Header().Set("Content-Type", "application/json;charset=utf8")
	json.NewEncoder(w).Encode(&event)
}

func (eh *eventServiceHandler) AllEvent(w http.ResponseWriter, r *http.Request) {
	events, err := eh.dbhandler.FindAllEvents()
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{error: Error occured while trying to find all events %s}`, err)
	}
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	err = json.NewEncoder(w).Encode(&events)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{error: Error while encoding events to JSON %s`, err)
	}
}

func (eh *eventServiceHandler) NewEvent(w http.ResponseWriter, r *http.Request)  {
	event := persistence.Event{}
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{error: error while decoding event data %s}`, err)
		return
	}
	id, err := eh.dbhandler.AddEvent(event)
	if err != nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, `{error: error occured while saving event %d %s}`, id, err)
		return
	}

	msg := contracts.EventCreatedEvent{
		ID: hex.EncodeToString(id),
		Name: event.Name,
		LocationID: event.Location.ID,
	}
	eh.eventEmitter.Emit(&msg)
}
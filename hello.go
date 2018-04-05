package main

import (
  "github.com/gorilla/mux"
  "net/http"
)

func main() {
  mainRouter := mux.NewRouter()
  eventsRouter := mainRouter.PathPrefix("/events").Subrouter()

  eventsRouter.Methods("GET").Path("{SearchCriteria}/{query}").HandlerFunc(findEvent)
  eventsRouter.Methods("GET").Path("").HandlerFunc(allEvent)
  eventsRouter.Methods("POST").Path("").HandlerFunc(newEvent)
  //fmt.Println("hello world", mainRouter);

  http.ListenAndServe(":8181", mainRouter)
}

func findEvent(w http.ResponseWriter, r *http.Request) {

}

func allEvent(w http.ResponseWriter, r *http.Request) {

}

func newEvent(w http.ResponseWriter, r *http.Request) {

}
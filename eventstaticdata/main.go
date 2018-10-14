package eventstaticdata

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Init registers the endpoints exposed by this package
// with the given Router.
// Also initialises the static data database connection
func Init(r *mux.Router) {

	initConn()

	r.HandleFunc("/events", getEventsHandler).Methods("GET")
	r.HandleFunc("/events", postEventsHandler).Methods("POST")
	r.HandleFunc("/events/{eventId}", getEventHandler).Methods("GET")
	r.HandleFunc("/events/{eventId}/map", postEventMapHandler).Methods("POST")
	r.HandleFunc("/events/{eventId}/map", getEventMapHandler).Methods("GET")
}

func getEventsHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: handle

}

func postEventsHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var postedEvent Event

	err := decoder.Decode(&postedEvent)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprint("Failed to decode event: %s", err),
			http.StatusBadRequest)
		return
	}

	err = validateEvent(&postedEvent)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Invalid event: %s", err)
			http.StatusBadRequest)
		return
	}

	createdEvent, err := addEvent(&postedEvent)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to write event to database: %s", err)
			http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(createdEvent)

}

func getEventHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: handle

}

func postEventMapHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: handle

}

func getEventMapHandler(w http.ResponseWriter, r *http.Request) {

	// TODO: handle

}

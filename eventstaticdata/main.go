package eventstaticdata

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Init registers the endpoints exposed by this package
// with the given Router.
// Also initialises the static data database connection
func Init(r *mux.Router) {

	fetchEnvVars()

	r.HandleFunc("/events", getEventsHandler).Methods("GET")
	r.HandleFunc("/events", postEventsHandler).Methods("POST")
	r.HandleFunc("/events/{eventId}", getEventHandler).Methods("GET")
	r.HandleFunc("/events/{eventId}/map", postEventMapHandler).Methods("POST")
}

func getEventsHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var allEventsReq AllEventsRequest

	err := decoder.Decode(&allEventsReq)

	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to unmarshall get all events request: %s", err),
			http.StatusBadRequest)
		return
	}

	err = validateAllEventsRequest(&allEventsReq)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Invalid request for all Events: %s", err),
			http.StatusBadRequest)
		return
	}

	events, err := getAllEventsByOrganiserID(allEventsReq.OrganiserID)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to get all events by organiser ID: %s", err),
			http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(events)

}

func postEventsHandler(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)

	var event Event

	err := decoder.Decode(&event)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to unmarshall Event: %s", err),
			http.StatusBadRequest)
		return
	}

	err = validateEvent(&event)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Invalid Event: %s", err),
			http.StatusBadRequest)
		return
	}

	err = addEvent(&event)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to write Event to database: %s", err),
			http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(event)

}

func getEventHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	idStr := vars["eventId"]

	if idStr == "" {
		http.Error(
			w,
			fmt.Sprint("Missing Event ID"),
			http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to parse Get Event by ID request: %s", err),
			http.StatusBadRequest)
		return
	}

	event, err := getEventByID(id)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to Get Event by ID: %s", err),
			http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(event)

}

func postEventMapHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	eventIDStr := vars["eventId"]

	if eventIDStr == "" {
		http.Error(
			w,
			fmt.Sprint("Missing Event ID"),
			http.StatusBadRequest)
		return
	}

	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to parse Post Map request: %s", err),
			http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(r.Body)

	var eventMap Map
	err = decoder.Decode(&eventMap)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to unmarshall map: %s", err),
			http.StatusBadRequest)
		return
	}

	err = validateMap(&eventMap)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Invalid Event: %s", err),
			http.StatusBadRequest)
		return
	}

	err = addMap(&eventMap)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to write map to database: %s", err),
			http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(eventMap)

}

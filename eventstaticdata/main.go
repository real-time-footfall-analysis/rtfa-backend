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
	r.HandleFunc("/events/{eventId}/regions", postRegionsHandler).Methods("POST")
	r.HandleFunc("/events/{eventId}/regions", getAllRegionsHandler).Methods("GET")
	r.HandleFunc("/events/{eventId}/regions/{regionId}", getRegionHandler).Methods("GET")
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
	eventMap.EventID = int32(eventID)

	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to unmarshall map: %s", err),
			http.StatusBadRequest)
		return
	}

	err = validateMap(&eventMap, eventID)
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

func postRegionsHandler(w http.ResponseWriter, r *http.Request) {

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
			fmt.Sprintf("Failed to parse event id: %s", err),
			http.StatusBadRequest)
		return
	}

	var regions []Region

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&regions)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to unmarshall regions: %s", err),
			http.StatusBadRequest)
		return
	}

	for _, region := range regions {
		err = validateRegion(&region, eventID)
		if err != nil {
			log.Println(err)
			http.Error(
				w,
				fmt.Sprintf("Invalid Region: %s", err),
				http.StatusBadRequest)
			return
		}
	}

	err = addRegions(&regions)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to write regions to database: %s", err),
			http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(regions)

}

func getAllRegionsHandler(w http.ResponseWriter, r *http.Request) {

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
			fmt.Sprintf("Failed to parse event id: %s", err),
			http.StatusBadRequest)
		return
	}

	regions, err := getRegionsByEventID(eventID)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to Get Regions by Event ID: %s", err),
			http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(regions)

}

func getRegionHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	eventIDStr := vars["eventId"]
	regionIDStr := vars["regionId"]

	if eventIDStr == "" {
		http.Error(
			w,
			fmt.Sprint("Missing Event ID"),
			http.StatusBadRequest)
		return
	}

	if regionIDStr == "" {
		http.Error(
			w,
			fmt.Sprint("Missing Region ID"),
			http.StatusBadRequest)
		return
	}

	eventID, err := strconv.Atoi(eventIDStr)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to parse event id: %s", err),
			http.StatusBadRequest)
		return
	}

	regionID, err := strconv.Atoi(regionIDStr)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to parse region id: %s", err),
			http.StatusBadRequest)
		return
	}

	region, err := getRegionByID(eventID, regionID)
	if err != nil {
		log.Println(err)
		http.Error(
			w,
			fmt.Sprintf("Failed to Get Region by Event ID: %s", err),
			http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(region)

}

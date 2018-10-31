package locationupdate

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// Init registers the endpoints exposed by this package
// with the given Router.
// Also initialises the static data database connection
var queue queue_adapter = &kenisis_queue{}

func Init(r *mux.Router) {

	queue.initConn()

	r.HandleFunc("/update", updateHandler).Methods("POST")
}

const (
	UUID_MIN_LENGTH = 5
)

type update struct {
	UUID       *string `json:"uuid"`
	EventID    *int    `json:"eventId"`
	RegionID   *int    `json:"regionId"`
	Entering   *bool   `json:"entering"`
	OccurredAt *int    `json:"occurredAt"`
}

func notPresentError(writer http.ResponseWriter, name string) {
	log.Println(name + " not present in update")
	http.Error(
		writer,
		fmt.Sprintf(name+" not present in update"),
		http.StatusBadRequest)
}

func notPresentCheck(writer http.ResponseWriter, update update) bool {
	if update.UUID == nil {
		notPresentError(writer, "UUID")
		return true
	}
	if update.EventID == nil {
		notPresentError(writer, "EventID")
		return true
	}
	if update.RegionID == nil {
		notPresentError(writer, "RegionID")
		return true
	}
	if update.Entering == nil {
		notPresentError(writer, "Entering")
		return true
	}
	if update.OccurredAt == nil {
		notPresentError(writer, "OccurredAt")
		return true
	}
	return false
}

func updateHandler(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	var update update

	err := decoder.Decode(&update)

	if err != nil {
		log.Println("Cannot decode update:", err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode update: %s", err),
			http.StatusBadRequest)
		return
	}

	if notPresentCheck(writer, update) {
		return
	}

	if len(*update.UUID) < UUID_MIN_LENGTH {
		log.Println("UUID less than 5 characters in update")
		http.Error(
			writer,
			fmt.Sprintf("UUID less than 5 characters"),
			http.StatusBadRequest)
		return
	}

	if *update.EventID < 0 {
		log.Println("Invalid EventId in update")
		http.Error(
			writer,
			fmt.Sprintf("Invalid EventId"),
			http.StatusBadRequest)
		return
	}

	// TODO: replace with actual timestamp from frontend
	now := int(time.Now().Unix())
	update.OccurredAt = &now

	queue.addLocationUpdate(&update)

	return
}

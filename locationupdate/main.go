package locationupdate

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

// Init registers the endpoints exposed by this package
// with the given Router.
// Also initialises the static data database connection
var queue queue_adapter = &kenisis_queue{}

func Init(r *mux.Router) {

	_ = queue.initConn()

	r.HandleFunc("/update", updateHandler).Methods("POST")
}

const (
	UUID_MIN_LENGTH = 5
)

type Movement_update struct {
	UUID       *string `json:"uuid"`
	EventID    *int    `json:"eventId"`
	RegionID   *int    `json:"regionId"`
	Entering   *bool   `json:"entering"`
	OccurredAt *int    `json:"occurredAt"`
}

func notPresentError(writer http.ResponseWriter, name string) {
	log.Println(name + " not present in movement update")
	http.Error(
		writer,
		fmt.Sprintf(name+" not present in movement update"),
		http.StatusBadRequest)
}

func notPresentCheck(writer http.ResponseWriter, update Movement_update) bool {
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

	var update Movement_update

	err := decoder.Decode(&update)

	if err != nil {
		log.Println("Cannot decode movement update:", err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode movement update: %s", err),
			http.StatusBadRequest)
		return
	}

	if notPresentCheck(writer, update) {
		return
	}

	if len(*update.UUID) < UUID_MIN_LENGTH {
		log.Println("UUID less than 5 characters in movement update")
		http.Error(
			writer,
			fmt.Sprintf("UUID less than 5 characters"),
			http.StatusBadRequest)
		return
	}

	if *update.EventID < 0 {
		log.Println("Invalid EventId in movement update")
		http.Error(
			writer,
			fmt.Sprintf("Invalid EventId"),
			http.StatusBadRequest)
		return
	}

	if *update.OccurredAt < 0 {
		log.Println("Invalid OccurredAt in movement update")
		http.Error(
			writer,
			fmt.Sprintf("Invalid occurredAt (timestamp)"),
			http.StatusBadRequest)
		return
	}

	// Send the data to the kinesis stream
	err = queue.addLocationUpdate(&update)
	if err != nil {
		log.Println("Error sending data to Kinesis")
		log.Println(err.Error())
	}
}

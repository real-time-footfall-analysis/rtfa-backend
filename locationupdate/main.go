package locationupdate

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/real-time-footfall-analysis/rtfa-backend/kinesisqueue"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

const kinesisStreamName = "movement_event_stream"

// Init registers the endpoints exposed by this package
// with the given Router.
// Also initialises the static data database connection
var queue kinesisqueue.KinesisQueueInterface = &kinesisqueue.KenisisQueueClient{}

func Init(r *mux.Router) {

	err := queue.InitConn(kinesisStreamName)
	if err != nil {
		log.Println("Failed to connect to Kinesis: " + kinesisStreamName)
		os.Exit(1)
	}

	r.HandleFunc("/update", updateHandler).Methods("POST")
	r.HandleFunc("/bulkUpdate", bulkKensisUpdateHandler).Methods("POST")
}

const (
	UUID_LENGTH = 36
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

	err = validateUpdate(&update, writer)
	if err != nil {
		return
	}

	// Send the data to the kinesis stream
	err = queue.SendToQueue(update, strconv.Itoa(*update.RegionID))
	if err != nil {
		log.Println("Error sending data to Kinesis")
		log.Println(err.Error())
	}
}

func bulkKensisUpdateHandler(writer http.ResponseWriter, request *http.Request) {

	decoder := json.NewDecoder(request.Body)

	var updates []Movement_update

	err := decoder.Decode(&updates)

	if err != nil {
		log.Println("Cannot decode movement updates: ", err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode movement updates: %s", err),
			http.StatusBadRequest)
		return
	}

	// Validate and insert into the database
	for _, update := range updates {

		err := validateUpdate(&update, writer)
		if err != nil {
			return
		}

		err = queue.SendToQueue(update, strconv.Itoa(*update.RegionID))
		if err != nil {
			log.Printf("Error posting movement update to Kenisis: %+v", update)
			log.Println(err.Error())
		}

	}

}

func validateUpdate(update *Movement_update, writer http.ResponseWriter) error {

	if notPresentCheck(writer, *update) {
		return errors.New("Movement update missing fields")
	}

	if len(*update.UUID) != UUID_LENGTH {
		msg := fmt.Sprintf("UUID not 36 characters in movement update %+v", update)
		log.Println(msg)
		http.Error(
			writer,
			msg,
			http.StatusBadRequest)
		return errors.New(msg)
	}

	if *update.EventID < 0 {
		msg := fmt.Sprintf("Invalid EventId in movement update %+v", update)
		log.Println(msg)
		http.Error(
			writer,
			msg,
			http.StatusBadRequest)
		return errors.New(msg)
	}

	if *update.OccurredAt < 0 {
		msg := fmt.Sprintf("Invalid OccurredAt in movement update %+v", update)
		log.Println(msg)
		http.Error(
			writer,
			msg,
			http.StatusBadRequest)
		return errors.New(msg)
	}

	return nil

}

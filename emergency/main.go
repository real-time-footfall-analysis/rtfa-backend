package emergency

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mitchellh/mapstructure"
	"github.com/real-time-footfall-analysis/rtfa-backend/dynamoDB"
	"github.com/real-time-footfall-analysis/rtfa-backend/pusher"
	"github.com/real-time-footfall-analysis/rtfa-backend/utils"
	"log"
	"net/http"
	"os"
	"strconv"
)

// Init registers the endpoints exposed by this package
// with the given Router.
// Also initialises the static data database connection

const (
	UUID_MIN_LENGTH = 5
)

type emergency_request struct {
	UUID        string `json:"uuid"`
	EventId     int    `json:"eventId"`
	RegionIds   []int  `json:"regionIds"`
	OccurredAt  int    `json:"occurredAt"`
	DealtWith   bool   `json:"dealtWith"`
	Description string `json:"description"`
	Position    *struct {
		Lat float64 `json:"lat"`
		Lng float64 `json:"lng"`
	} `json:"position"`
}

var db dynamoDB.DynamoDBInterface = &dynamoDB.DynamoDBClient{}
var pc pusher.PusherChannelInterface = &pusher.PusherChannelClient{}

func Init(r *mux.Router) {
	err := db.InitConn("emergency_events")
	if err != nil {
		os.Exit(1)
	}

	pc.InitConn()
	r.HandleFunc("/emergency-update", updateHandler).Methods("POST")
	r.HandleFunc("/live/emergency/{eventId}/{lastPoll}", requestHandler).Methods("GET")
}

func updateHandler(writer http.ResponseWriter, request *http.Request) {

	utils.SetAccessControlHeaders(writer)

	decoder := json.NewDecoder(request.Body)

	var emergencyUpdate emergency_request

	// Try and decode the data
	err := decoder.Decode(&emergencyUpdate)
	if err != nil {
		log.Println("Cannot decode emergency_request:", err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode emergency_request: %s", err),
			http.StatusBadRequest)
		return
	}

	// Check the other fields in the data
	if emergencyUpdate.UUID == "" || len(emergencyUpdate.UUID) < UUID_MIN_LENGTH {
		log.Println("uuid field is empty:", err)
		http.Error(
			writer,
			fmt.Sprintf("uuid field is empty %s", err),
			http.StatusBadRequest)
		return
	}
	if emergencyUpdate.RegionIds == nil {
		log.Println("RegionIds missing:", err)
		http.Error(
			writer,
			fmt.Sprintf("RegionIds missing: %s", err),
			http.StatusBadRequest)
		return
	}
	if emergencyUpdate.OccurredAt == 0 {
		log.Println("occurredAt timestamp missing", err)
		http.Error(
			writer,
			fmt.Sprintf("occurredAt timestamp missing: %s", err),
			http.StatusBadRequest)
		return
	}
	if emergencyUpdate.EventId <= 0 {
		log.Println("Invalid EventId in emergency_request")
		http.Error(
			writer,
			fmt.Sprintf("Invalid EventId"),
			http.StatusBadRequest)
		return
	}

	// Send the item to the database
	db.SendItem(emergencyUpdate)

	// Push the item to Pusher
	channelName := strconv.Itoa(emergencyUpdate.EventId)
	data, _ := json.Marshal(emergencyUpdate)
	pc.SendItem(channelName, "emergency-update", data)

	// Return the update to the user
	_ = json.NewEncoder(writer).Encode(emergencyUpdate)
}

func requestHandler(writer http.ResponseWriter, request *http.Request) {

	// Allow cross origin
	utils.SetAccessControlHeaders(writer)

	// Get the payload variables out
	vars := mux.Vars(request)

	// Convert the eventId to an int
	eventId, err := parseRequestArgs(vars, "eventId", writer)
	if err != nil {
		log.Println("Error parsing eventId")
		http.Error(
			writer,
			fmt.Sprintf("Invalid EventId"),
			http.StatusBadRequest)
		return
	}

	// Convert the last poll timestamp to an int
	lastPoll, err := parseRequestArgs(vars, "lastPoll", writer)
	if err != nil {
		log.Println("Error parsing last poll time")
		http.Error(
			writer,
			fmt.Sprintf("Invalid last poll time"),
			http.StatusBadRequest)
		return
	}

	// Scan the table, parse the result and filter out old events
	unparsedRows := db.GetTableScan()
	var parsedRows []emergency_request = make([]emergency_request, len(unparsedRows))
	for index, row := range unparsedRows {
		_ = mapstructure.Decode(row, &parsedRows[index])
	}

	// Extract relevant values
	recents := extractRecentUpdates(parsedRows, eventId, lastPoll)

	// Transmit the result back
	_ = json.NewEncoder(writer).Encode(recents)
}

func extractRecentUpdates(parsed []emergency_request, event int, prevPoll int) (res []emergency_request) {
	// Remove the values that don't satisfy a criteria
	deleted := 0
	for i := range parsed {
		j := i - deleted
		if parsed[j].EventId != event || parsed[j].OccurredAt < prevPoll {
			parsed = parsed[:j+copy(parsed[j:], parsed[j+1:])]
			deleted++
		}
	}
	return parsed
}

func parseRequestArgs(vars map[string]string, varName string, writer http.ResponseWriter) (int, error) {
	id, err := strconv.Atoi(vars[varName])
	if err != nil {
		log.Println("Cannot decode request "+varName, err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode request: %s", err),
			http.StatusBadRequest)
	}
	return id, err
}

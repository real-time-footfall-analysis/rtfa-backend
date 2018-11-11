package emergency

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
	"github.com/real-time-footfall-analysis/rtfa-backend/utils"
	"log"
	"net/http"
	"strconv"
)

// Init registers the endpoints exposed by this package
// with the given Router.
// Also initialises the static data database connection

const (
	UUID_MIN_LENGTH = 5
)

type emergency_request struct {
	UUID       string `json:"uuid"`
	EventId    int    `json:"eventId"`
	RegionIds  []int  `json:"regionIds"`
	OccurredAt int    `json:"occurredAt"`
	Sorted     bool   `json:"sorted"`
}

var db emergencyDbAdapter = &dynamoDbAdaptor{}

func Init(r *mux.Router) {
	db.initConn()
	r.HandleFunc("/emergency_update", updateHandler).Methods("POST")
	r.HandleFunc("/live/emergency/{eventId}/{lastPoll}", requestHandler).Methods("GET")
}

func notPresentError(writer http.ResponseWriter, name string) {
	log.Println(name + " not present in emergency_request")
	http.Error(
		writer,
		fmt.Sprintf(name+" not present in emergency_request"),
		http.StatusBadRequest)
}

func notPresentCheck(writer http.ResponseWriter, update emergency_request) bool {
	if update.UUID == "" {
		notPresentError(writer, "UUID")
		return true
	}
	if update.EventId == 0 {
		notPresentError(writer, "EventID")
		return true
	}
	if update.RegionIds == nil {
		notPresentError(writer, "RegionID")
		return true
	}
	if update.OccurredAt == 0 {
		notPresentError(writer, "OccurredAt")
		return true
	}
	return false
}

func updateHandler(writer http.ResponseWriter, request *http.Request) {
	decoder := json.NewDecoder(request.Body)

	var emergency_update emergency_request

	err := decoder.Decode(&emergency_update)

	if err != nil {
		log.Println("Cannot decode emergency_request:", err)
		http.Error(
			writer,
			fmt.Sprintf("Failed to decode emergency_request: %s", err),
			http.StatusBadRequest)
		return
	}

	if notPresentCheck(writer, emergency_update) {
		return
	}

	if len(emergency_update.UUID) < UUID_MIN_LENGTH {
		log.Println("UUID less than 5 characters in emergency_request")
		http.Error(
			writer,
			fmt.Sprintf("UUID less than 5 characters"),
			http.StatusBadRequest)
		return
	}

	if emergency_update.EventId < 0 {
		log.Println("Invalid EventId in emergency_request")
		http.Error(
			writer,
			fmt.Sprintf("Invalid EventId"),
			http.StatusBadRequest)
		return
	}

	// Send the item
	db.sendItem(emergency_update)
}

func requestHandler(writer http.ResponseWriter, request *http.Request) {

	// Allow cross origin
	utils.SetAccessControlHeaders(writer)

	// Get the payload variables out
	vars := mux.Vars(request)

	// Convert the eventId to an int
	eventId, err := parseRequestArgs(vars, "eventId", writer)
	if err != nil {
		return
	}

	// Convert the last poll timestamp to an int
	lastPoll, err := parseRequestArgs(vars, "lastPoll", writer)
	if err != nil {
		return
	}

	// Scan the table, parse the result and filter out old events
	scan, err := db.getTableScan()
	parsed := parseScan(scan)
	recents := extractRecentUpdates(parsed, eventId, lastPoll)

	// Transmit the result back
	json.NewEncoder(writer).Encode(recents)
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

func parseScan(tableScan *dynamodb.ScanOutput) []emergency_request {
	var parsed []emergency_request

	// Parse the table scanned items
	for _, row := range tableScan.Items {
		// Extract the data
		eventId, _ := strconv.Atoi(*(*row["eventId"]).N)
		uuid := *(*row["uuid"]).S
		sorted := *(*row["sorted"]).BOOL
		occurredAt, _ := strconv.Atoi(*(*row["occurredAt"]).N)

		// Parse the regionIds
		unparsed_regions := (*(row["regionIds"])).L
		var regions []int
		for _, reg := range unparsed_regions {
			regionId, _ := strconv.Atoi(*reg.N)
			regions = append(regions, regionId)
		}

		// Insert into a new emergency request
		parsed_row := emergency_request{EventId: eventId,
			UUID:       uuid,
			OccurredAt: occurredAt,
			RegionIds:  regions,
			Sorted:     sorted}

		// Add to final list
		parsed = append(parsed, parsed_row)
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

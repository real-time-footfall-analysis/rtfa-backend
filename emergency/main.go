package emergency

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gorilla/mux"
	"github.com/real-time-footfall-analysis/rtfa-backend/utils"
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

var db emergencyDbAdapter = &dynamoDbAdaptor{}

func Init(r *mux.Router) {
	db.initConn()
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
	err = db.sendItem(emergencyUpdate)
	if err != nil {
		log.Println(err.Error())
	}

	// Return the update to the user
	json.NewEncoder(writer).Encode(emergencyUpdate)
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
	scan, err := db.getTableScan()
	if err != nil {
		log.Println("Error getting the scan of the emergencies DynamoDB table")
	}

	// Parse the scan and remove old values
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
		// Parse the eventId
		eventId, err := strconv.Atoi(*(*row["eventId"]).N)
		if err != nil {
			continue
		}

		// Parse the timestamp
		occurredAt, err := strconv.Atoi(*(*row["occurredAt"]).N)
		if err != nil {
			continue
		}

		// Extract the uuid and dealtWith with boolean
		uuid := *(*row["uuid"]).S
		dealtWith := *(*row["dealtWith"]).BOOL

		// Extract the description if present
		var description string
		if (*row["description"]).S != nil {
			description = *(*row["description"]).S
		} else {
			description = ""
		}

		// Parse the regionIds
		unparsedRegions := (*(row["regionIds"])).L
		var regions = []int{}
		for _, reg := range unparsedRegions {
			regionId, _ := strconv.Atoi(*reg.N)
			regions = append(regions, regionId)
		}

		// Insert into a new emergency request
		parsedRow := emergency_request{EventId: eventId,
			UUID:        uuid,
			OccurredAt:  occurredAt,
			RegionIds:   regions,
			DealtWith:   dealtWith,
			Description: description,
		}

		// Parse the GPS Position
		if (*row["position"]).M != nil {
			positionMap := (*row["position"]).M
			lat, err := strconv.ParseFloat(*positionMap["lat"].N, 32)
			if err != nil {
				continue
			}

			lng, err := strconv.ParseFloat(*positionMap["lng"].N, 32)
			if err != nil {
				continue
			}

			// Add it to the struct
			parsedRow.Position = &struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			}{Lat: lat, Lng: lng}
		}

		// Add to final list
		parsed = append(parsed, parsedRow)
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
